package TopRankCache

import (
	"GraduationDesign/BaseMent/Config"
	"GraduationDesign/BaseMent/Model"
	"GraduationDesign/BaseMent/utils"
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"log"
	"strconv"
	"time"
)

var C = context.Background()

func MusicListRankAdd(id uint) int {
	m, d, n := time.Now().Month().String(), time.Now().Weekday().String(), strconv.Itoa(time.Now().Day())
	y1, w1 := time.Now().ISOWeek()
	w := strconv.Itoa(w1)
	y := strconv.Itoa(y1)
	id1 := strconv.Itoa(int(id))
	key := y + " " + m + " " + w + " " + d + " " + n + " " + id1
	err := Config.DBR2.ZIncrBy(c, Config.Conf.Rank.CacheRankMusicListDay, 1, key).Err()
	if err != nil {
		num, _ := Config.DBR2.ZScore(c, Config.Conf.Rank.CacheRankMusicListDay, key).Result()
		log.Println("歌单次数缓存出错:", err, "具体键值信息:musicListId:", id1, "count:", num)
		err = Config.DBR2.ZAdd(c, Config.Conf.Rank.CacheRankMusicListDay, redis.Z{Member: key, Score: 1 + num}).Err()
		if err != nil {
			log.Println("日听排行歌单缓存预载失败", err, "musicListId:", id1)
			return utils.ERROR
		}
	}
	return utils.SUCCESS
}

func GetMusicListRankList() ([]Model.MusicList, int, int64) {
	nums, err := Config.DBR2.ZRevRangeWithScores(c, Config.Conf.Rank.CacheRankMusicListDay, 0, 30).Result()
	fmt.Println(err)
	if err != nil {
		log.Println("获取排行榜失败:", err)
		return nil, utils.ERROR, 0
	}
	musicList, code, total := Model.GetAMusicListList(nums)
	return musicList, code, total
}
