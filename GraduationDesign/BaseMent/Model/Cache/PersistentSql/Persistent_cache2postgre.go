package PersistentSql

import (
	"GraduationDesign/BaseMent/Model"
	"GraduationDesign/BaseMent/utils"
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"log"
	"strconv"
	"strings"
	"sync"
	"time"
)

// QueueForTime is a message queue supporting delayed/scheduled delivery based on redis
type QueueForTime struct {
	// name for this Queue. Make sure the name is unique in redis database
	name          string
	redisCli      *redis.Client
	cb            func(string) bool
	pendingKey    string // sorted set: message id -> delivery time
	readyKey      string // list
	unAckKey      string // sorted set: message id -> retry time
	retryKey      string // list
	retryCountKey string // hash: message id -> remain retry count
	garbageKey    string // set: message id
	ticker        *time.Ticker
	logger        *log.Logger
	close         chan struct{}

	maxConsumeDuration time.Duration
	msgTTL             time.Duration
	defaultRetryCount  uint
	fetchInterval      time.Duration
	fetchLimit         uint

	concurrent uint
}

type hashTagKeyOpt int

// UseHashTagKey add hashtags to redis keys to ensure all keys of this queue are allocated in the same hash slot.
// If you are using Codis/AliyunRedisCluster/TencentCloudRedisCluster, add this option to NewQueue
// WARNING! Changing (add or remove) this option will cause DelayQueue failing to read existed data in redis
// see more:  https://redis.io/docs/reference/cluster-spec/#hash-tags
func UseHashTagKey() interface{} {
	return hashTagKeyOpt(1)
}

// NewQueue creates a new queue, use DelayQueue.StartConsume to consume or DelayQueue.SendScheduleMsg to publish message
// callback returns true to confirm successful consumption. If callback returns false or not return within maxConsumeDuration, DelayQueue will re-deliver this message
func NewQueue(name string, cli *redis.Client, callback func(string) bool, opts ...interface{}) *QueueForTime {
	if name == "" {
		panic("name is required")
	}
	if cli == nil {
		panic("cli is required")
	}
	if callback == nil {
		panic("callback is required")
	}
	useHashTag := false
	for _, opt := range opts {
		switch opt.(type) {
		case hashTagKeyOpt:
			useHashTag = true
		}
	}
	var keyPrefix string
	if useHashTag {
		keyPrefix = "{dp:" + name + "}"
	} else {
		keyPrefix = "dp:" + name
	}
	return &QueueForTime{
		name:               name,
		redisCli:           cli,
		cb:                 callback,
		pendingKey:         keyPrefix + ":pending",
		readyKey:           keyPrefix + ":ready",
		unAckKey:           keyPrefix + ":unack",
		retryKey:           keyPrefix + ":retry",
		retryCountKey:      keyPrefix + ":retry:cnt",
		garbageKey:         keyPrefix + ":garbage",
		close:              make(chan struct{}, 1),
		maxConsumeDuration: 5 * time.Second,
		msgTTL:             time.Hour,
		logger:             log.Default(),
		defaultRetryCount:  3,
		fetchInterval:      time.Second,
		concurrent:         1,
	}
}

// WithLogger customizes logger for queue
func (q *QueueForTime) WithLogger(logger *log.Logger) *QueueForTime {
	q.logger = logger
	return q
}

// WithFetchInterval customizes the interval at which consumer fetch message from redis
func (q *QueueForTime) WithFetchInterval(d time.Duration) *QueueForTime {
	q.fetchInterval = d
	return q
}

// WithMaxConsumeDuration customizes max consume duration
// If no acknowledge received within WithMaxConsumeDuration after message delivery, DelayQueue will try to deliver this message again
func (q *QueueForTime) WithMaxConsumeDuration(d time.Duration) *QueueForTime {
	q.maxConsumeDuration = d
	return q
}

// WithFetchLimit limits the max number of unack (processing) messages
func (q *QueueForTime) WithFetchLimit(limit uint) *QueueForTime {
	q.fetchLimit = limit
	return q
}

// WithConcurrent sets the number of concurrent consumers
func (q *QueueForTime) WithConcurrent(c uint) *QueueForTime {
	if c == 0 {
		return q
	}
	q.concurrent = c
	return q
}

// WithDefaultRetryCount customizes the max number of retry, it effects of messages in this queue
// use WithRetryCount during DelayQueue.SendScheduleMsg or DelayQueue.SendDelayMsg to specific retry count of particular message
func (q *QueueForTime) WithDefaultRetryCount(count uint) *QueueForTime {
	q.defaultRetryCount = count
	return q
}

func (q *QueueForTime) genMsgKey(idStr string) string {
	return "dp:" + q.name + ":msg:" + idStr
}

type retryCountOpt int

// WithRetryCount set retry count for a msg
// example: queue.SendDelayMsg(payload, duration, delayqueue.WithRetryCount(3))
func WithRetryCount(count int) interface{} {
	return retryCountOpt(count)
}

type msgTTLOpt time.Duration

// WithMsgTTL set ttl for a msg
// example: queue.SendDelayMsg(payload, duration, delayqueue.WithMsgTTL(Hour))
func WithMsgTTL(d time.Duration) interface{} {
	return msgTTLOpt(d)
}

// SendScheduleMsg submits a message delivered at given time
func (q *QueueForTime) SendScheduleMsg(payload string, t time.Time, opts ...interface{}) error {
	// parse options
	retryCount := q.defaultRetryCount
	for _, opt := range opts {
		switch o := opt.(type) {
		case retryCountOpt:
			retryCount = uint(o)
		case msgTTLOpt:
			q.msgTTL = time.Duration(o)
		}
	}
	// generate id
	idStr := uuid.Must(uuid.NewRandom()).String()
	ctx := context.Background()
	now := time.Now()
	// store msg
	msgTTL := t.Sub(now) + q.msgTTL // delivery + q.msgTTL
	err := q.redisCli.Set(ctx, q.genMsgKey(idStr), payload, msgTTL).Err()
	if err != nil {
		return fmt.Errorf("store msg failed: %v", err)
	}
	// store retry count
	err = q.redisCli.HSet(ctx, q.retryCountKey, idStr, retryCount).Err()
	if err != nil {
		return fmt.Errorf("store retry count failed: %v", err)
	}
	// put to pending
	err = q.redisCli.ZAdd(ctx, q.pendingKey, redis.Z{Score: float64(t.Unix()), Member: idStr}).Err()
	if err != nil {
		return fmt.Errorf("push to pending failed: %v", err)
	}
	return nil
}

// SendDelayMsg submits a message delivered after given duration
func (q *QueueForTime) SendDelayMsg(payload string, duration time.Duration, opts ...interface{}) error {
	t := time.Now().Add(duration)
	return q.SendScheduleMsg(payload, t, opts...)
}

// pending2ReadyScript atomically moves messages from pending to ready
// keys: pendingKey, readyKey
// argv: currentTime
const pending2ReadyScript = `
local msgs = redis.call('ZRangeByScore', KEYS[1], '0', ARGV[1])  -- get ready msg
if (#msgs == 0) then return end
local args2 = {'LPush', KEYS[2]} -- push into ready
for _,v in ipairs(msgs) do
	table.insert(args2, v) 
    if (#args2 == 4000) then
		redis.call(unpack(args2))
		args2 = {'LPush', KEYS[2]}
	end
end
if (#args2 > 2) then 
	redis.call(unpack(args2))
end
redis.call('ZRemRangeByScore', KEYS[1], '0', ARGV[1])  -- remove msgs from pending
`

func (q *QueueForTime) pending2Ready() error {
	now := time.Now().Unix()
	ctx := context.Background()
	keys := []string{q.pendingKey, q.readyKey}
	err := q.redisCli.Eval(ctx, pending2ReadyScript, keys, now).Err()
	if err != nil && err != redis.Nil {
		return fmt.Errorf("pending2ReadyScript failed: %v", err)
	}
	return nil
}

// ready2UnackScript atomically moves messages from ready to unack
// keys: readyKey/retryKey, unackKey
// argv: retryTime
const ready2UnackScript = `
local msg = redis.call('RPop', KEYS[1])
if (not msg) then return end
redis.call('ZAdd', KEYS[2], ARGV[1], msg)
return msg
`

func (q *QueueForTime) ready2Unack() (string, error) {
	retryTime := time.Now().Add(q.maxConsumeDuration).Unix()
	ctx := context.Background()
	keys := []string{q.readyKey, q.unAckKey}
	ret, err := q.redisCli.Eval(ctx, ready2UnackScript, keys, retryTime).Result()
	if err == redis.Nil {
		return "", err
	}
	if err != nil {
		return "", fmt.Errorf("ready2UnackScript failed: %v", err)
	}
	str, ok := ret.(string)
	if !ok {
		return "", fmt.Errorf("illegal result: %#v", ret)
	}
	return str, nil
}

func (q *QueueForTime) retry2Unack() (string, error) {
	retryTime := time.Now().Add(q.maxConsumeDuration).Unix()
	ctx := context.Background()
	keys := []string{q.retryKey, q.unAckKey}
	ret, err := q.redisCli.Eval(ctx, ready2UnackScript, keys, retryTime, q.retryKey, q.unAckKey).Result()
	if err == redis.Nil {
		return "", redis.Nil
	}
	if err != nil {
		return "", fmt.Errorf("ready2UnackScript failed: %v", err)
	}
	str, ok := ret.(string)
	if !ok {
		return "", fmt.Errorf("illegal result: %#v", ret)
	}
	return str, nil
}

func (q *QueueForTime) callback(idStr string) error {
	ctx := context.Background()
	payload, err := q.redisCli.Get(ctx, q.genMsgKey(idStr)).Result()
	if err == redis.Nil {
		return nil
	}
	if err != nil {
		// Is an IO error?
		return fmt.Errorf("get message payload failed: %v", err)
	}
	ack := q.cb(payload)
	if ack {
		//join in a persistent way
		switch payload[len(payload)-1] {
		case '1':
			a, _ := strconv.Atoi(payload[:len(payload)-1])
			code := Model.MusicsListen(uint(a))
			if code == utils.ERROR {
				log.Println("MusicListen persistent fail!", code, "musicId:", a)
			}
		case '2':
			id1 := strings.Split(payload, " ")
			id2, _ := strconv.Atoi(id1[0])
			id3, _ := strconv.Atoi(id1[1][:len(id1[1])-1])
			code := Model.UserSongLike(uint(id2), uint(id3))
			if code == utils.ERROR {
				log.Println("MusicLike persistent fail", code, "musicId:", id2, "userId:", id3)
			}
		case '3':
			id1 := strings.Split(payload, " ")
			id2, _ := strconv.Atoi(id1[0])
			id3, _ := strconv.Atoi(id1[1][:len(id1[1])-1])
			code := Model.UserSongDisLike(uint(id2), uint(id3))
			if code == utils.ERROR {
				log.Println("MusicDisLike persistent fail", code, "musicId:", id2, "userId:", id3)
			}
		case '4':
			var data2 Model.Comment
			_ = json.Unmarshal([]byte(payload[:len(payload)-1]), &data2)
			code := Model.CreatAComment(&data2)
			if code == utils.ERROR {
				log.Println("NewComment persistent fail", code, "data:", payload[:len(payload)-1])
			}
		case '5':
			id1 := strings.Split(payload, " ")
			id2, _ := strconv.Atoi(id1[0])
			code := Model.LikeArticle(uint(id2), id1[1][:len(id1[1])-1])
			if code == utils.ERROR {
				log.Println("ArticleLike persistent fail", code, "articleId:", id2, "user_id", id1[1][:len(id1[1])-1])
			}
		case '6':
			id1 := strings.Split(payload, " ")
			id2, _ := strconv.Atoi(id1[0])
			code := Model.DisLikeArticle(uint(id2), id1[1][:len(id1[1])-1])
			if code == utils.ERROR {
				log.Println("ArticleDisLike persistent fail", code, "articleId:", id2, "user_id", id1[1][:len(id1[1])-1])
			}
		case '7':
			id1, _ := strconv.Atoi(payload[:len(payload)-1])
			code := Model.RecordReadCount(uint(id1))
			if code == utils.ERROR {
				log.Println("ArticleRead persistent fail", code, "articleId:", id1)
			}
		case '8':
			id4 := strings.Split(payload, " ")
			id5, _ := strconv.Atoi(id4[0])
			id6, _ := strconv.Atoi(id4[1])
			code := Model.LikeComment(uint(id5), uint(id6), id4[2][:len(id4[2])-1])
			if code == utils.ERROR {
				log.Println("CommentLike persistent fail", code, "articleId:", id5, "CommentId:", id6, "user_id:", id4[2][:len(id4[2])-1])
			}
		case '9':
			id4 := strings.Split(payload, " ")
			id5, _ := strconv.Atoi(id4[0])
			id6, _ := strconv.Atoi(id4[1])
			code := Model.DisLikeComment(uint(id5), uint(id6), id4[2][:len(id4[2])-1])
			if code == utils.ERROR {
				log.Println("CommentDisLike persistent fail", code, "articleId:", id5, "CommentId:", id6, "user_id:", id4[2][:len(id4[2])-1])
			}
		//read for the sort list but consider it not update frequently so to use postgres to save
		//case 'a':
		//
		//case 'b':
		case 'c':
			data := strings.Split(payload, " ")
			y := data[0]
			m := utils.GetCNTimeMonth(data[1])
			w := data[2]
			d := utils.GetCNTimeWeek(data[3])
			number := data[4]
			id := data[5]
			count := data[6][:len(data[6])-1]
			id1, _ := strconv.Atoi(id)
			count1, _ := strconv.Atoi(count)
			code := Model.MusicRankCount(uint(id1), y, m, w, d, number, count1)
			if code == utils.ERROR {
				log.Println("MusicRankCount persistent fail", code, "y", y, "m", m, "w", w, "d", d)
			}
		case 'd':
			data := strings.Split(payload, " ")
			y := data[0]
			m := utils.GetCNTimeMonth(data[1])
			w := data[2]
			d := utils.GetCNTimeWeek(data[3])
			number := data[4]
			id := data[5]
			count := data[6][:len(data[6])-1]
			id1, _ := strconv.Atoi(id)
			count1, _ := strconv.Atoi(count)
			code := Model.MusicListRankCount(uint(id1), y, m, w, d, number, count1)
			if code == utils.ERROR {
				log.Println("MusicListRankCount persistent fail", code, "y", y, "m", m, "w", w, "d", d)
			}
		case 'e':
			data := strings.Split(payload, " ")
			id, _ := strconv.Atoi(data[0])
			code := Model.ForwardArticle(uint(id), data[1][:len(data[1])-1])
			if code == utils.ERROR {
				log.Println("ForwardArticleCount persistent fail", code, "articleId:", id, "userId:", data[1][:len(data[1])-1])
			}
		case 'f':
			data := strings.Split(payload, " ")
			id, _ := strconv.Atoi(data[0])
			id1, _ := strconv.Atoi(data[1][:len(data)-1])
			code := Model.LikeMusicList(uint(id), uint(id1))
			if code == utils.ERROR {
				log.Println("LikeMusicListCount persistent fail", code, "articleId:", id, "userId:", data[1][:len(data[1])-1])
			}
		case 'g':
			data := strings.Split(payload, " ")
			id, _ := strconv.Atoi(data[0])
			id1, _ := strconv.Atoi(data[1][:len(data)-1])
			code := Model.DisLikeMusicList(uint(id), uint(id1))
			if code == utils.ERROR {
				log.Println("DisLikeMusicListCount persistent fail", code, "articleId:", id, "userId:", data[1][:len(data[1])-1])
			}
		}
		err = q.ack(idStr)
	} else {
		err = q.nack(idStr)
	}
	return err
}

// batchCallback calls DelayQueue.callback in batch. callback is executed concurrently according to property DelayQueue.concurrent
// batchCallback must wait all callback finished, otherwise the actual number of processing messages may beyond DelayQueue.FetchLimit
func (q *QueueForTime) batchCallback(ids []string) {
	if len(ids) == 1 || q.concurrent == 1 {
		for _, id := range ids {
			err := q.callback(id)
			if err != nil {
				q.logger.Printf("consume msg %s failed: %v", id, err)
			}
		}
		return
	}
	ch := make(chan string, len(ids))
	for _, id := range ids {
		ch <- id
	}
	close(ch)
	wg := sync.WaitGroup{}
	concurrent := int(q.concurrent)
	if concurrent > len(ids) { // too many goroutines is no use
		concurrent = len(ids)
	}
	wg.Add(concurrent)
	for i := 0; i < concurrent; i++ {
		go func() {
			defer wg.Done()
			for id := range ch {
				err := q.callback(id)
				if err != nil {
					q.logger.Printf("consume msg %s failed: %v", id, err)
				}
			}
		}()
	}
	wg.Wait()
}

func (q *QueueForTime) ack(idStr string) error {
	ctx := context.Background()
	err := q.redisCli.ZRem(ctx, q.unAckKey, idStr).Err()
	if err != nil {
		return fmt.Errorf("remove from unack failed: %v", err)
	}
	// msg key has ttl, ignore result of delete
	_ = q.redisCli.Del(ctx, q.genMsgKey(idStr)).Err()
	q.redisCli.HDel(ctx, q.retryCountKey, idStr)
	return nil
}

func (q *QueueForTime) nack(idStr string) error {
	ctx := context.Background()
	// update retry time as now, unack2Retry will move it to retry immediately
	err := q.redisCli.ZAdd(ctx, q.unAckKey, redis.Z{
		Member: idStr,
		Score:  float64(time.Now().Unix()),
	}).Err()
	if err != nil {
		return fmt.Errorf("negative ack failed: %v", err)
	}
	return nil
}

// unack2RetryScript atomically moves messages from unack to retry which remaining retry count greater than 0,
// and moves messages from unack to garbage which  retry count is 0
// Because DelayQueue cannot determine garbage message before eval unack2RetryScript, so it cannot pass keys parameter to redisCli.Eval
// Therefore unack2RetryScript moves garbage message to garbageKey instead of deleting directly
// keys: unackKey, retryCountKey, retryKey, garbageKey
// argv: currentTime
const unack2RetryScript = `
local unack2retry = function(msgs)
	local retryCounts = redis.call('HMGet', KEYS[2], unpack(msgs)) -- get retry count
	for i,v in ipairs(retryCounts) do
		local k = msgs[i]
		if v ~= false and v ~= nil and v ~= '' and tonumber(v) > 0 then
			redis.call("HIncrBy", KEYS[2], k, -1) -- reduce retry count
			redis.call("LPush", KEYS[3], k) -- add to retry
		else
			redis.call("HDel", KEYS[2], k) -- del retry count
			redis.call("SAdd", KEYS[4], k) -- add to garbage
		end
	end
end
local msgs = redis.call('ZRangeByScore', KEYS[1], '0', ARGV[1])  -- get retry msg
if (#msgs == 0) then return end
if #msgs < 4000 then
	unack2retry(msgs)
else
	local buf = {}
	for _,v in ipairs(msgs) do
		table.insert(buf, v)
		if #buf == 4000 then
			unack2retry(buf)
			buf = {}
		end
	end
	if (#buf > 0) then
		unack2retry(buf)
	end
end
redis.call('ZRemRangeByScore', KEYS[1], '0', ARGV[1])  -- remove msgs from unack
`

func (q *QueueForTime) unack2Retry() error {
	ctx := context.Background()
	keys := []string{q.unAckKey, q.retryCountKey, q.retryKey, q.garbageKey}
	now := time.Now()
	err := q.redisCli.Eval(ctx, unack2RetryScript, keys, now.Unix()).Err()
	if err != nil && err != redis.Nil {
		return fmt.Errorf("unack to retry script failed: %v", err)
	}
	return nil
}

func (q *QueueForTime) garbageCollect() error {
	ctx := context.Background()
	msgIds, err := q.redisCli.SMembers(ctx, q.garbageKey).Result()
	if err != nil {
		return fmt.Errorf("smembers failed: %v", err)
	}
	if len(msgIds) == 0 {
		return nil
	}
	// allow concurrent clean
	msgKeys := make([]string, 0, len(msgIds))
	for _, idStr := range msgIds {
		msgKeys = append(msgKeys, q.genMsgKey(idStr))
	}
	err = q.redisCli.Del(ctx, msgKeys...).Err()
	if err != nil && err != redis.Nil {
		return fmt.Errorf("del msgs failed: %v", err)
	}
	err = q.redisCli.SRem(ctx, q.garbageKey, msgIds).Err()
	if err != nil && err != redis.Nil {
		return fmt.Errorf("remove from garbage key failed: %v", err)
	}
	return nil
}

func (q *QueueForTime) consume() error {
	// pending to ready
	err := q.pending2Ready()
	if err != nil {
		return err
	}
	// consume
	ids := make([]string, 0, q.fetchLimit)
	for {
		idStr, err := q.ready2Unack()
		if err == redis.Nil { // consumed all
			break
		}
		if err != nil {
			return err
		}
		ids = append(ids, idStr)
		if q.fetchLimit > 0 && len(ids) >= int(q.fetchLimit) {
			break
		}
	}
	if len(ids) > 0 {
		q.batchCallback(ids)
	}
	// unack to retry
	err = q.unack2Retry()
	if err != nil {
		return err
	}
	err = q.garbageCollect()
	if err != nil {
		return err
	}
	// retry
	ids = make([]string, 0, q.fetchLimit)
	for {
		idStr, err := q.retry2Unack()
		if err == redis.Nil { // consumed all
			break
		}
		if err != nil {
			return err
		}
		ids = append(ids, idStr)
		if q.fetchLimit > 0 && len(ids) >= int(q.fetchLimit) {
			break
		}
	}
	if len(ids) > 0 {
		q.batchCallback(ids)
	}
	return nil
}

// StartConsume creates a goroutine to consume message from DelayQueue
// use `<-done` to wait consumer stopping
func (q *QueueForTime) StartConsume() (done <-chan struct{}) {
	q.ticker = time.NewTicker(q.fetchInterval)
	done0 := make(chan struct{})
	go func() {
	tickerLoop:
		for {
			select {
			case <-q.ticker.C:
				err := q.consume()
				if err != nil {
					log.Printf("consume error: %v", err)
				}
			case <-q.close:
				break tickerLoop
			}
		}
		close(done0)
	}()
	return done0
}

// StopConsume stops consumer goroutine
func (q *QueueForTime) StopConsume() {
	close(q.close)
	if q.ticker != nil {
		q.ticker.Stop()
	}
}

//import (
//	"context"
//	"fmt"
//	"github.com/google/uuid"
//	"github.com/redis/go-redis/v9"
//	"log"
//	"sync"
//	"time"
//)
//
//type QueueForTime struct {
//	Name         string
//	DBR          *redis.Client
//	Cb           func(string) bool
//	PendingKey   string
//	ReadyKey     string
//	UnAckKey     string
//	RetryKey     string
//	RetryTimeKey string
//	MsgId        string
//	Ticker       *time.Ticker
//	Log          *log.Logger
//	Close        chan struct{}
//
//	MaxConsumeDuration time.Duration
//	MsgTtl             time.Duration
//	DefaultRetryTime   uint
//	FetchInterval      time.Duration
//	FetchLimit         uint
//
//	Concurrent uint
//}
//
//type HashTagKeyOpt int
//
//func UseHashTagKey() interface{} {
//	return HashTagKeyOpt(1)
//}
//
//func NewQueue(name string, RDB *redis.Client, callback func(string) bool, opts ...interface{}) *QueueForTime {
//	if name == "" || RDB == nil || callback == nil {
//		return nil
//	}
//	UseHashTag := false
//	for _, opt := range opts {
//		switch opt.(type) {
//		case HashTagKeyOpt:
//			UseHashTag = true
//		}
//	}
//	var keyPreFix string
//	if UseHashTag {
//		keyPreFix = "{dp:" + name + "}"
//	}
//	keyPreFix = "dp:" + name
//	return &QueueForTime{
//		Name:               name,
//		DBR:                RDB,
//		Cb:                 callback,
//		PendingKey:         keyPreFix + ":pending",
//		UnAckKey:           keyPreFix + ":ready",
//		RetryKey:           keyPreFix + ":retry",
//		RetryTimeKey:       keyPreFix + ":retry:time",
//		MsgId:              keyPreFix + ":msg",
//		Close:              make(chan struct{}, 1),
//		MaxConsumeDuration: 5 * time.Second,
//		Log:                log.Default(),
//		DefaultRetryTime:   2,
//		FetchInterval:      time.Second,
//		Concurrent:         1,
//	}
//}
//
//func (q *QueueForTime) WithLogger(logger *log.Logger) *QueueForTime {
//	q.Log = logger
//	return q
//}
//
//func (q *QueueForTime) WithFetchInterval(d time.Duration) *QueueForTime {
//	q.FetchInterval = d
//	return q
//}
//
//func (q *QueueForTime) WithMaxConsumeDuration(d time.Duration) *QueueForTime {
//	q.MaxConsumeDuration = d
//	return q
//}
//
//func (q *QueueForTime) WithFetchLimit(limit uint) *QueueForTime {
//	q.FetchLimit = limit
//	return q
//}
//
//func (q *QueueForTime) WithConcurrent(c uint) *QueueForTime {
//	if c == 0 {
//		return q
//	}
//	q.Concurrent = c
//	return q
//}
//
//func (q *QueueForTime) WithDefaultRetryCount(count uint) *QueueForTime {
//	q.DefaultRetryTime = count
//	return q
//}
//
//func (q *QueueForTime) genMsgKey(idStr string) string {
//	return "dp:" + q.Name + ":msg:" + idStr
//}
//
//type RetryCountOpt int
//
//func WithRetryCount(count int) interface{} {
//	return RetryCountOpt(count)
//}
//
//type MsgTtlOpt time.Duration
//
//func WithMsgTTL(d time.Duration) interface{} {
//	return MsgTtlOpt(d)
//}
//
//func (q *QueueForTime) SendScheduleMsg(payload string, t time.Time, opts ...interface{}) error {
//	// parse options
//	retryCount := q.DefaultRetryTime
//	for _, opt := range opts {
//		switch o := opt.(type) {
//		case RetryCountOpt:
//			retryCount = uint(o)
//		case MsgTtlOpt:
//			q.MsgTtl = time.Duration(o)
//		}
//	}
//	// generate id
//	idStr := uuid.Must(uuid.NewRandom()).String()
//	ctx := context.Background()
//	now := time.Now()
//	// store msg
//	msgTTL := t.Sub(now) + q.MsgTtl // delivery + q.msgTTL
//	err := q.DBR.Set(ctx, q.genMsgKey(idStr), payload, msgTTL).Err()
//	if err != nil {
//		return fmt.Errorf("store msg failed: %v", err)
//	}
//	// store retry count
//	err = q.DBR.HSet(ctx, q.RetryTimeKey, idStr, retryCount).Err()
//	if err != nil {
//		return fmt.Errorf("store retry count failed: %v", err)
//	}
//	// put to pending
//	err = q.DBR.ZAdd(ctx, q.PendingKey, redis.Z{Score: float64(t.Unix()), Member: idStr}).Err()
//	if err != nil {
//		return fmt.Errorf("push to pending failed: %v", err)
//	}
//	return nil
//}
//
//// pending2ReadyScript atomically moves messages from pending to ready
//// keys: PendingKey, readyKey
//// argv: currentTime
//const pending2ReadyScript = `
//local msgs = redis.call('ZRangeByScore', KEYS[1], '0', ARGV[1])  -- get ready msg
//if (#msgs == 0) then return end
//local args2 = {'LPush', KEYS[2]} -- push into ready
//for _,v in ipairs(msgs) do
//	table.insert(args2, v)
//    if (#args2 == 4000) then
//		redis.call(unpack(args2))
//		args2 = {'LPush', KEYS[2]}
//	end
//end
//if (#args2 > 2) then
//	redis.call(unpack(args2))
//end
//redis.call('ZRemRangeByScore', KEYS[1], '0', ARGV[1])  -- remove msgs from pending
//`
//
//func (q *QueueForTime) pending2Ready() error {
//	now := time.Now().Unix()
//	ctx := context.Background()
//	keys := []string{q.PendingKey, q.ReadyKey}
//	err := q.DBR.Eval(ctx, pending2ReadyScript, keys, now).Err()
//	if err != nil && err != redis.Nil {
//		return fmt.Errorf("pending2ReadyScript failed: %v", err)
//	}
//	return nil
//}
//
//// ready2UnackScript atomically moves messages from ready to unack
//// keys: readyKey/RetryKey, UnAckKey
//// argv: retryTime
//const ready2UnackScript = `
//local msg = redis.call('RPop', KEYS[1])
//if (not msg) then return end
//redis.call('ZAdd', KEYS[2], ARGV[1], msg)
//return msg
//`
//
//func (q *QueueForTime) ready2Unack() (string, error) {
//	retryTime := time.Now().Add(q.MaxConsumeDuration).Unix()
//	ctx := context.Background()
//	keys := []string{q.ReadyKey, q.UnAckKey}
//	ret, err := q.DBR.Eval(ctx, ready2UnackScript, keys, retryTime).Result()
//	if err == redis.Nil {
//		return "", err
//	}
//	if err != nil {
//		return "", fmt.Errorf("ready2UnackScript failed: %v", err)
//	}
//	str, ok := ret.(string)
//	if !ok {
//		return "", fmt.Errorf("illegal result: %#v", ret)
//	}
//	return str, nil
//}
//
//func (q *QueueForTime) retry2Unack() (string, error) {
//	retryTime := time.Now().Add(q.MaxConsumeDuration).Unix()
//	ctx := context.Background()
//	keys := []string{q.RetryKey, q.UnAckKey}
//	ret, err := q.DBR.Eval(ctx, ready2UnackScript, keys, retryTime, q.RetryKey, q.UnAckKey).Result()
//	if err == redis.Nil {
//		return "", redis.Nil
//	}
//	if err != nil {
//		return "", fmt.Errorf("ready2UnackScript failed: %v", err)
//	}
//	str, ok := ret.(string)
//	if !ok {
//		return "", fmt.Errorf("illegal result: %#v", ret)
//	}
//	return str, nil
//}
//
//func (q *QueueForTime) callback(idStr string) error {
//	ctx := context.Background()
//	payload, err := q.DBR.Get(ctx, q.genMsgKey(idStr)).Result()
//	if err == redis.Nil {
//		return nil
//	}
//	if err != nil {
//		// Is an IO error?
//		return fmt.Errorf("get message payload failed: %v", err)
//	}
//	ack := q.Cb(payload)
//	if ack {
//		fmt.Println(payload, "ack")
//		err = q.ack(idStr)
//	} else {
//		fmt.Println(payload, "ackrepeate")
//		err = q.nack(idStr)
//	}
//	return err
//}
//
//// batchCallback calls QueueForTime.callback in batch. callback is executed concurrently according to property QueueForTime.concurrent
//// batchCallback must wait all callback finished, otherwise the actual number of processing messages may beyond QueueForTime.FetchLimit
//func (q *QueueForTime) batchCallback(ids []string) {
//	if len(ids) == 1 || q.Concurrent == 1 {
//		for _, id := range ids {
//			err := q.callback(id)
//			if err != nil {
//				q.Log.Printf("consume msg %s failed: %v", id, err)
//			}
//		}
//		return
//	}
//	ch := make(chan string, len(ids))
//	for _, id := range ids {
//		ch <- id
//	}
//	close(ch)
//	wg := sync.WaitGroup{}
//	concurrent := int(q.Concurrent)
//	if concurrent > len(ids) { // too many goroutines is no use
//		concurrent = len(ids)
//	}
//	wg.Add(concurrent)
//	for i := 0; i < concurrent; i++ {
//		go func() {
//			defer wg.Done()
//			for id := range ch {
//				err := q.callback(id)
//				if err != nil {
//					q.Log.Printf("consume msg %s failed: %v", id, err)
//				}
//			}
//		}()
//	}
//	wg.Wait()
//}
//
//func (q *QueueForTime) ack(idStr string) error {
//	ctx := context.Background()
//	err := q.DBR.ZRem(ctx, q.UnAckKey, idStr).Err()
//	if err != nil {
//		return fmt.Errorf("remove from unack failed: %v", err)
//	}
//	// msg key has ttl, ignore result of delete
//	_ = q.DBR.Del(ctx, q.genMsgKey(idStr)).Err()
//	q.DBR.HDel(ctx, q.RetryTimeKey, idStr)
//	return nil
//}
//
//func (q *QueueForTime) nack(idStr string) error {
//	ctx := context.Background()
//	// update retry time as now, unack2Retry will move it to retry immediately
//	err := q.DBR.ZAdd(ctx, q.UnAckKey, redis.Z{
//		Member: idStr,
//		Score:  float64(time.Now().Unix()),
//	}).Err()
//	if err != nil {
//		return fmt.Errorf("negative ack failed: %v", err)
//	}
//	return nil
//}
//
//// unack2RetryScript atomically moves messages from unack to retry which remaining retry count greater than 0,
//// and moves messages from unack to garbage which  retry count is 0
//// Because QueueForTime cannot determine garbage message before eval unack2RetryScript, so it cannot pass keys parameter to DBR.Eval
//// Therefore unack2RetryScript moves garbage message to garbageKey instead of deleting directly
//// keys: UnAckKey, RetryTimeKey, RetryKey, garbageKey
//// argv: currentTime
//const unack2RetryScript = `
//local unack2retry = function(msgs)
//	local retryCounts = redis.call('HMGet', KEYS[2], unpack(msgs)) -- get retry count
//	for i,v in ipairs(retryCounts) do
//		local k = msgs[i]
//		if v ~= false and v ~= nil and v ~= '' and tonumber(v) > 0 then
//			redis.call("HIncrBy", KEYS[2], k, -1) -- reduce retry count
//			redis.call("LPush", KEYS[3], k) -- add to retry
//		else
//			redis.call("HDel", KEYS[2], k) -- del retry count
//			redis.call("SAdd", KEYS[4], k) -- add to garbage
//		end
//	end
//end
//local msgs = redis.call('ZRangeByScore', KEYS[1], '0', ARGV[1])  -- get retry msg
//if (#msgs == 0) then return end
//if #msgs < 4000 then
//	unack2retry(msgs)
//else
//	local buf = {}
//	for _,v in ipairs(msgs) do
//		table.insert(buf, v)
//		if #buf == 4000 then
//			unack2retry(buf)
//			buf = {}
//		end
//	end
//	if (#buf > 0) then
//		unack2retry(buf)
//	end
//end
//redis.call('ZRemRangeByScore', KEYS[1], '0', ARGV[1])  -- remove msgs from unack
//`
//
//func (q *QueueForTime) unack2Retry() error {
//	ctx := context.Background()
//	keys := []string{q.UnAckKey, q.RetryTimeKey, q.RetryKey, q.MsgId}
//	now := time.Now()
//	err := q.DBR.Eval(ctx, unack2RetryScript, keys, now.Unix()).Err()
//	if err != nil && err != redis.Nil {
//		return fmt.Errorf("unack to retry script failed: %v", err)
//	}
//	return nil
//}
//
//func (q *QueueForTime) garbageCollect() error {
//	ctx := context.Background()
//	msgIds, err := q.DBR.SMembers(ctx, q.MsgId).Result()
//	if err != nil {
//		return fmt.Errorf("smembers failed: %v", err)
//	}
//	if len(msgIds) == 0 {
//		return nil
//	}
//	// allow concurrent clean
//	msgKeys := make([]string, 0, len(msgIds))
//	for _, idStr := range msgIds {
//		msgKeys = append(msgKeys, q.genMsgKey(idStr))
//	}
//	err = q.DBR.Del(ctx, msgKeys...).Err()
//	if err != nil && err != redis.Nil {
//		return fmt.Errorf("del msgs failed: %v", err)
//	}
//	err = q.DBR.SRem(ctx, q.MsgId, msgIds).Err()
//	if err != nil && err != redis.Nil {
//		return fmt.Errorf("remove from garbage key failed: %v", err)
//	}
//	return nil
//}
//
//func (q *QueueForTime) consume() error {
//	// pending to ready
//	err := q.pending2Ready()
//	if err != nil {
//		return err
//	}
//	// consume
//	ids := make([]string, 0, q.FetchLimit)
//	for {
//		idStr, err := q.ready2Unack()
//		if err == redis.Nil { // consumed all
//			break
//		}
//		if err != nil {
//			return err
//		}
//		ids = append(ids, idStr)
//		if q.FetchLimit > 0 && len(ids) >= int(q.FetchLimit) {
//			break
//		}
//	}
//	if len(ids) > 0 {
//		q.batchCallback(ids)
//	}
//	// unack to retry
//	err = q.unack2Retry()
//	if err != nil {
//		return err
//	}
//	err = q.garbageCollect()
//	if err != nil {
//		return err
//	}
//	// retry
//	ids = make([]string, 0, q.FetchLimit)
//	for {
//		idStr, err := q.retry2Unack()
//		if err == redis.Nil { // consumed all
//			break
//		}
//		if err != nil {
//			return err
//		}
//		ids = append(ids, idStr)
//		if q.FetchLimit > 0 && len(ids) >= int(q.FetchLimit) {
//			break
//		}
//	}
//	if len(ids) > 0 {
//		q.batchCallback(ids)
//	}
//	return nil
//}
//
//// StartConsume creates a goroutine to consume message from QueueForTime
//// use `<-done` to wait consumer stopping
//func (q *QueueForTime) StartConsume() (done <-chan struct{}) {
//	q.Ticker = time.NewTicker(q.FetchInterval)
//	done0 := make(chan struct{})
//	go func() {
//	tickerLoop:
//		for {
//			select {
//			case <-q.Ticker.C:
//				err := q.consume()
//				fmt.Println("11111")
//				if err != nil {
//					log.Printf("consume error: %v", err)
//				}
//			case <-q.Close:
//				fmt.Println("close")
//				break tickerLoop
//			}
//		}
//		close(done0)
//	}()
//	return done0
//}
//
//// StopConsume stops consumer goroutine
//func (q *QueueForTime) StopConsume() {
//	close(q.Close)
//	if q.Ticker != nil {
//		q.Ticker.Stop()
//	}
//}
