package Cache

import (
	"GraduationDesign/BaseMent/Config"
	"GraduationDesign/BaseMent/Model"
	"GraduationDesign/BaseMent/Model/Cache/PersistentSql"
	"GraduationDesign/BaseMent/utils"
	"context"
	"encoding/json"
	"fmt"
	"github.com/redis/go-redis/v9"
	"log"
	"strconv"
	"time"
)

var c = context.Background()
var ArticlePass = 1 * time.Hour
var BasementMinute = 1 * time.Minute
var MusicRankPassBase = 1 * time.Hour

func GetArticle(id uint) (*Model.Article, int) {
	var article *Model.Article
	key := utils.GetCacheNameArticle(id)
	article, err := GetCacheArticle(id)
	if err == redis.Nil || err != nil || article.ID == 0 {
		article, code := Model.CheckAArticle(id)
		if code == utils.SUCCESS {
			err = CreatACacheArticle(id, &article)
			if err == nil {
				return &article, utils.SUCCESS
			}
			return &article, utils.ERROR
		}
	}
	Config.DBR.Expire(c, key, ArticlePass)
	return article, utils.SUCCESS
}

func GetCacheArticle(id uint) (*Model.Article, error) {
	var article Model.Article
	key := utils.GetCacheNameArticle(id)
	data, err := Config.DBR.Get(c, key).Result()
	if err == redis.Nil || err != nil {
		log.Println(err)
		return &article, err
	}
	if err = json.Unmarshal([]byte(data), &article); err != nil {
		log.Println(err)
		return &article, err
	}
	return &article, err
}

func CreatACacheArticle(id uint, article *Model.Article) error {
	key := utils.GetCacheNameArticle(id)
	data, err := json.Marshal(article)
	if err != nil {
		return err
	}
	err = Config.DBR.Set(c, key, data, ArticlePass).Err()
	if err != nil {
		return err
	}
	return nil
}

func DeleteACacheArticle(id uint) {
	key := utils.GetCacheNameArticle(id)
	_, _ = Config.DBR.Del(c, key).Result()
}

func UpdateCacheArticle(id uint, article *Model.Article) {
	key := utils.GetCacheNameArticle(id)
	_, err := Config.DBR.Get(c, key).Result()
	if err == redis.Nil || err != nil {
		log.Println(err)
	}
	_, _ = Config.DBR.Del(c, key).Result()
	err = CreatACacheArticle(id, article)
	if err != nil {
		log.Println(err)
	}
}

func UpdateCacheMusicListened(id uint) int {
	err := PersistentSql.CaCheMusicListen.SendScheduleMsg(strconv.Itoa(int(id))+"1", time.Now().Add(0), PersistentSql.WithRetryCount(2))
	if err != nil {
		return utils.ERROR
	}
	return utils.SUCCESS
}

func UpdateMusicLikeCount(id uint) int64 {
	key := utils.GetCacheNameMusicLikeCount(id)
	data, err := Config.DBR.Get(c, key).Result()
	if err != nil || len(data) <= 2 {
		code, total := Model.MusicsLikeCount(id)
		if code == utils.ERROR {
			return 0
		}
		data1, _ := json.Marshal(total)
		err = Config.DBR.Set(c, key, string(data1), ArticlePass).Err()
		if err != nil {
			log.Println("保存喜欢音乐失败!", err, "musicId:", id)
		}
		return total
	}
	Config.DBR.Expire(c, key, ArticlePass)
	var count int64
	json.Unmarshal([]byte(data), &count)
	return count
}

func CheckCacheMusicIsLike(MusicId uint, userId uint) bool {
	key := utils.GetCacheNameMusicLike(MusicId, userId)
	data, err := Config.DBR.Get(c, key).Result()
	if err != nil || len(data) <= 2 {
		status := Model.CheckSongLike(MusicId, userId)
		data1, _ := json.Marshal(status)
		err = Config.DBR.Set(c, key, string(data1), 5*BasementMinute).Err()
		if err != nil {
			log.Println("歌曲喜欢缓存储存失败。", err)
		}
		return status
	}
	Config.DBR.Expire(c, key, 5*BasementMinute)
	var status bool
	json.Unmarshal([]byte(data), &status)
	return status
}

func UpdateCaCheMusicLike(musicId uint, userId uint) int {
	a, b := strconv.Itoa(int(musicId)), strconv.Itoa(int(userId))
	id := a + " " + b + "2"
	err := PersistentSql.CaCheMusicLike.SendScheduleMsg(id, time.Now().Add(0), PersistentSql.WithRetryCount(2))
	if err != nil {
		return utils.ERROR
	}
	return utils.SUCCESS
}

func UpdateCaCheMusicDisLike(musicId uint, userId uint) int {
	a, b := strconv.Itoa(int(musicId)), strconv.Itoa(int(userId))
	id := a + " " + b + "3"
	err := PersistentSql.CaCheMusicLike.SendScheduleMsg(id, time.Now().Add(0), PersistentSql.WithRetryCount(2))
	if err != nil {
		return utils.ERROR
	}
	return utils.SUCCESS
}

// Have no idea weather to use it or not
func UpdateCaCheArticleComment(data *Model.Comment) int {
	data1, _ := json.Marshal(data)
	err := PersistentSql.CaCheArticleComment.SendScheduleMsg(string(data1)+"4", time.Now().Add(0), PersistentSql.WithRetryCount(2))
	if err != nil {
		return utils.ERROR
	}
	return utils.SUCCESS
}

func GetCacheArticleCommentOne(id uint, rid uint, pageSize int, pageNum int) ([]Model.Comment, int, int64) {
	key := utils.GetCacheNameCommentONE(id, rid)
	type Comment struct {
		comments []Model.Comment
		total    int64
	}
	data, err := Config.DBR.Get(c, key).Result()
	if err != nil || len(data) <= 2 {
		var comment Comment
		data1, total := Model.GetCommentListInAComment(id, rid, pageSize, pageNum)
		comment.comments = data1
		comment.total = total
		data2, _ := json.Marshal(&comment)
		err = Config.DBR.Set(c, key, string(data2), 5*BasementMinute).Err()
		if err != nil {
			log.Println("评论缓存储存失败。", err)
		}
		return comment.comments, utils.SUCCESS, comment.total
	}
	Config.DBR.Expire(c, key, 5*BasementMinute)
	var result Comment
	json.Unmarshal([]byte(data), &result)
	return result.comments, utils.SUCCESS, result.total
}

func GetCacheArticleComments(id uint, pageSize int, pageNum int) ([]Model.Comment, int, int64) {
	key := utils.GetCacheNameComments(id)
	type Comment struct {
		comments []Model.Comment
		total    int64
	}
	data, err := Config.DBR.Get(c, key).Result()
	if err != nil || len(data) <= 2 {
		var comment Comment
		data1, total := Model.GetCommentListByArticleId(id, pageSize, pageNum)
		comment.comments = data1
		comment.total = total
		data2, _ := json.Marshal(&comment)
		err = Config.DBR.Set(c, key, string(data2), 10*BasementMinute).Err()
		if err != nil {
			log.Println("文章评论缓存储存失败。", err)
		}
		return comment.comments, utils.SUCCESS, comment.total
	}
	Config.DBR.Expire(c, key, 10*BasementMinute)
	var result Comment
	json.Unmarshal([]byte(data), &result)
	return result.comments, utils.SUCCESS, result.total
}

func GetCacheComments(id uint, id2 uint, pageSize int, pageNum int) ([]Model.Comment, int, int64) {
	//key := utils.GetCacheNameComments(id)
	data1, total := Model.GetCommentListByArticleIdT(id, id2, pageSize, pageNum)
	return data1, utils.SUCCESS, total
}

func CheckACaCheArticleLike(articleId uint, userId string) bool {
	key := utils.GetCacheNameArticleLike(articleId, userId)
	data, err := Config.DBR.Get(c, key).Result()
	if err != nil || len(data) <= 2 {
		status := Model.CheckArticleLike(articleId, userId)
		data1, _ := json.Marshal(status)
		err = Config.DBR.Set(c, key, string(data1), ArticlePass).Err()
		if err != nil {
			log.Println("喜欢标记失败", err, "articleId:", articleId, "userid:", userId)
		}
		return status
	}
	Config.DBR.Expire(c, key, ArticlePass)
	var status bool
	json.Unmarshal([]byte(data), &status)
	return status
}

func UpdateCaCheArticleLike(articleId uint, userId string) int {
	key := utils.GetCacheNameArticleLike(articleId, userId)
	data1, _ := json.Marshal(true)
	_ = Config.DBR.Set(c, key, string(data1), ArticlePass).Err()
	articleId1 := strconv.Itoa(int(articleId))
	id := articleId1 + " " + userId + "5"
	err := PersistentSql.CaCheArticleLike.SendScheduleMsg(id, time.Now().Add(0), PersistentSql.WithRetryCount(2))
	if err != nil {
		return utils.ERROR
	}
	return utils.SUCCESS
}

func UpdateCaCheArticleDisLike(articleId uint, userId string) int {
	key := utils.GetCacheNameArticleLike(articleId, userId)
	data1, _ := json.Marshal(false)
	_ = Config.DBR.Set(c, key, string(data1), ArticlePass).Err()
	articleId1 := strconv.Itoa(int(articleId))
	id := articleId1 + " " + userId + "6"
	err := PersistentSql.CaCheArticleLike.SendScheduleMsg(id, time.Now().Add(0), PersistentSql.WithRetryCount(2))
	if err != nil {
		return utils.ERROR
	}
	return utils.SUCCESS
}

func UpdateCaCheArticleRead(articleId uint) int {
	id := strconv.Itoa(int(articleId)) + "7"
	err := PersistentSql.CaCheArticleRead.SendScheduleMsg(id, time.Now().Add(0), PersistentSql.WithRetryCount(2))
	if err != nil {
		return utils.ERROR
	}
	return utils.SUCCESS
}

func UpdateCaCheArticleForward(articleId uint, userId string) int {
	id := strconv.Itoa(int(articleId)) + " " + userId + "e"
	err := PersistentSql.CacheArticleForward.SendScheduleMsg(id, time.Now().Add(0), PersistentSql.WithRetryCount(1))
	if err != nil {
		return utils.ERROR
	}
	return utils.SUCCESS
}

func UpdateCacheCommentLike(ArticleId uint, commentId uint, userId string) int {
	key := utils.GetCacheNameCommentLike(commentId, userId)
	_ = Config.DBR.Set(c, key, true, 10*BasementMinute).Err()
	id := strconv.Itoa(int(ArticleId))
	id2 := strconv.Itoa(int(commentId))
	id3 := id + " " + id2 + " " + userId + "8"
	err := PersistentSql.CacheCommentLike.SendScheduleMsg(id3, time.Now().Add(0), PersistentSql.WithRetryCount(2))
	if err != nil {
		return utils.ERROR
	}
	return utils.SUCCESS
}

func UpdateCacheCommentIsLike(commentId uint, userId string) bool {
	key := utils.GetCacheNameCommentLike(commentId, userId)
	data, err := Config.DBR.Get(c, key).Result()
	if err != nil || len(data) <= 2 {
		status := Model.CheckCommentLike(commentId, userId)
		data1, _ := json.Marshal(status)
		err = Config.DBR.Set(c, key, string(data1), 10*BasementMinute).Err()
		if err != nil {
			log.Println("评论喜欢换成失败:", err, "commentId:", commentId, "userId:", userId)
		}
		return status
	}
	Config.DBR.Expire(c, key, 10*BasementMinute)
	var result bool
	json.Unmarshal([]byte(data), &result)
	return result
}

func UpdateCacheCommentDisLike(ArticleId uint, commentId uint, userId string) int {
	key := utils.GetCacheNameCommentLike(commentId, userId)
	_ = Config.DBR.Set(c, key, false, 10*BasementMinute).Err()
	id := strconv.Itoa(int(ArticleId))
	id2 := strconv.Itoa(int(commentId))
	id3 := id + " " + id2 + " " + userId + "9"
	err := PersistentSql.CacheCommentLike.SendScheduleMsg(id3, time.Now().Add(0), PersistentSql.WithRetryCount(2))
	if err != nil {
		return utils.ERROR
	}
	return utils.SUCCESS
}

func UpdateCacheMusicListIsLike(MusicListId uint, userId uint) bool {
	key := utils.GetCacheNameMusicListLike(MusicListId, userId)
	data, err := Config.DBR.Get(c, key).Result()
	if err != nil || len(data) <= 2 {
		status := Model.CheckMusicListLike(MusicListId, userId)
		data1, _ := json.Marshal(status)
		err = Config.DBR.Set(c, key, string(data1), 10*BasementMinute).Err()
		if err != nil {
			log.Println("歌单喜欢缓存失败:", err, "MusicListId:", MusicListId, "userId:", userId)
		}
		return status
	}
	Config.DBR.Expire(c, key, 10*BasementMinute)
	var result bool
	json.Unmarshal([]byte(data), &result)
	return result
}

func UpdateCacheMusicListLike(MusicListId uint, userId uint) int {
	key := strconv.Itoa(int(MusicListId)) + " " + strconv.Itoa(int(userId)) + "f"
	err := PersistentSql.CaCheMusicListLike.SendScheduleMsg(key, time.Now().Add(0), PersistentSql.WithRetryCount(2))
	if err != nil {
		return utils.ERROR
	}
	return utils.SUCCESS
}

func UpdateCacheMusicListDisLike(MusicListId uint, userId uint) int {
	key := strconv.Itoa(int(MusicListId)) + " " + strconv.Itoa(int(userId)) + "g"
	err := PersistentSql.CaCheMusicListLike.SendScheduleMsg(key, time.Now().Add(0), PersistentSql.WithRetryCount(2))
	if err != nil {
		return utils.ERROR
	}
	return utils.SUCCESS
}

//read for the sort list but consider it not update frequently so to use postgres to save
//func UpdateCacheCacheMusicListSort(ArticleId uint, commentId uint, userId string) int {
//	id := strconv.Itoa(int(ArticleId))
//	id2 := strconv.Itoa(int(commentId))
//	id3 := id + " " + id2 + " " + userId + "a"
//	err := PersistentSql.CacheMusicListSort.SendScheduleMsg(id3, time.Now().Add(0), PersistentSql.WithRetryCount(2))
//	if err != nil {
//		return utils.ERROR
//	}
//	return utils.SUCCESS
//}
//
//func UpdateCacheCacheMusicSort(ArticleId uint, commentId uint, userId string) int {
//	id := strconv.Itoa(int(ArticleId))
//	id2 := strconv.Itoa(int(commentId))
//	id3 := id + " " + id2 + " " + userId + "b"
//	err := PersistentSql.CacheMusicSort.SendScheduleMsg(id3, time.Now().Add(0), PersistentSql.WithRetryCount(2))
//	if err != nil {
//		return utils.ERROR
//	}
//	return utils.SUCCESS
//}

// PersistentMusicDayRank Persistent Musics Day Rank
func PersistentMusicDayRank() {
	nums, err := Config.DBR2.ZRevRangeWithScores(context.Background(), Config.Conf.Rank.CacheRankMusicDay, 0, -1).Result()
	fmt.Println(nums)
	if err != nil {
		log.Println("缓存获取排行榜失败:", err)
	}
	for i := 0; i < len(nums); i++ {
		a := nums[i].Member.(string)
		b := strconv.Itoa(int(nums[i].Score))
		key := a + " " + b + "c"
		err = PersistentSql.CacheMusicDayListen.SendScheduleMsg(key, time.Now().Add(0), PersistentSql.WithRetryCount(2))
		if err != nil {
			log.Println("排名歌曲持久化失败:", err, "musicId:", a[len(a)-1], "count:", b, "Rank:", i)
		}
	}
}

// PersistentMusicListDayRank Persistent MusicLists Day Rank
func PersistentMusicListDayRank() {
	nums, err := Config.DBR2.ZRevRangeWithScores(context.Background(), Config.Conf.Rank.CacheRankMusicListDay, 0, -1).Result()
	if err != nil {
		log.Println("获取排行榜失败:", err)
	}
	for i := 0; i < len(nums); i++ {
		a := nums[i].Member.(string)
		b := strconv.Itoa(int(nums[i].Score))
		key := a + " " + b + "d"
		err = PersistentSql.CacheMusicListDayListen.SendScheduleMsg(key, time.Now().Add(0), PersistentSql.WithRetryCount(2))
		if err != nil {
			log.Println("排名歌单持久化失败:", err, "musicListId:", a[len(a)-1], "count:", b, "Rank:", i)
		}
	}
}

func GetACacheMusicListRankWeek() ([]Model.MusicList, int, int64) {
	key := utils.GetCacheNameMusicListRank("week")
	data, err := Config.DBR2.Get(context.Background(), key).Result()
	if err != redis.Nil && err != nil {
		log.Println(err)
		return nil, utils.ERROR, 0
	} else if err == redis.Nil || len(data) <= 2 {
		return nil, utils.ERROR, 0
	}
	var musicsList struct {
		musicsList []Model.MusicList
		total      int64
	}
	json.Unmarshal([]byte(data), &musicsList)
	return musicsList.musicsList, utils.SUCCESS, musicsList.total
}

func UpdateMusicListRankWeek() {
	var data []Model.MusicList
	key := utils.GetCacheNameMusicListRank("week")
	y, w1 := time.Now().ISOWeek()
	m := utils.GetCNTimeMonth(time.Now().Month().String())
	data, _, total := Model.GetMusicListRankWeek(y, m, w1)
	var musicsList struct {
		musicsList []Model.MusicList
		total      int64
	}
	musicsList.musicsList = data
	musicsList.total = total
	data1, err := json.Marshal(musicsList)
	if err != nil {
		log.Println(err)
	}
	_, err = Config.DBR2.Get(context.Background(), key).Result()
	if err != redis.Nil && err != nil {
		log.Println(err)
	}
	if total == 0 {
		Config.DBR2.Expire(context.Background(), key, 49*ArticlePass)
		log.Println("周歌单排行使用了旧数据！")
	}
	_, _ = Config.DBR2.Del(context.Background(), key).Result()
	err = Config.DBR2.Set(context.Background(), key, string(data1), 49*ArticlePass).Err()
	if err != nil {
		log.Println("周听歌单榜缓存失败:", err, "y", "m", m, "w", w1)
	}
}

func GetACacheMusicListRankMonth() ([]Model.MusicList, int, int64) {
	key := utils.GetCacheNameMusicListRank("month")
	data, err := Config.DBR2.Get(context.Background(), key).Result()
	if err != redis.Nil && err != nil {
		log.Println(err)
		return nil, utils.ERROR, 0
	} else if err == redis.Nil || len(data) <= 2 {
		return nil, utils.ERROR, 0
	}
	var musicsList struct {
		musicsList []Model.MusicList
		total      int64
	}
	json.Unmarshal([]byte(data), &musicsList)
	return musicsList.musicsList, utils.SUCCESS, musicsList.total
}

func UpdateMusicListRankMonth() {
	var data []Model.MusicList
	key := utils.GetCacheNameMusicListRank("month")
	y, _ := time.Now().ISOWeek()
	m := utils.GetCNTimeMonth(time.Now().Month().String())
	data, _, total := Model.GetMusicListRankMonth(y, m)
	var musicsList struct {
		musicsList []Model.MusicList
		total      int64
	}
	musicsList.musicsList = data
	musicsList.total = total
	data1, err := json.Marshal(musicsList)
	if err != nil {
		log.Println(err)
	}
	_, err = Config.DBR2.Get(context.Background(), key).Result()
	if err != redis.Nil && err != nil {
		log.Println(err)
	}
	if total == 0 {
		Config.DBR2.Expire(context.Background(), key, 193*ArticlePass)
		log.Println("月歌单排行榜使用了数据缓存！")
	}
	_, _ = Config.DBR2.Del(context.Background(), key).Result()
	err = Config.DBR2.Set(context.Background(), key, string(data1), 193*ArticlePass).Err()
	if err != nil {
		log.Println("月听歌单榜缓存失败:", err, "y", "m", m)
	}
}

func GetACacheMusicListRankYear() ([]Model.MusicList, int, int64) {
	key := utils.GetCacheNameMusicListRank("year")
	data, err := Config.DBR2.Get(context.Background(), key).Result()
	if err != redis.Nil && err != nil {
		log.Println(err)
		return nil, utils.ERROR, 0
	} else if err == redis.Nil || len(data) <= 2 {
		return nil, utils.ERROR, 0
	}
	var musicsList struct {
		musicsList []Model.MusicList
		total      int64
	}
	json.Unmarshal([]byte(data), &musicsList)
	return musicsList.musicsList, utils.SUCCESS, musicsList.total
}

func UpdateMusicListRankYear() {
	var data []Model.MusicList
	key := utils.GetCacheNameMusicListRank("year")
	y, _ := time.Now().ISOWeek()
	data, _, total := Model.GetMusicListRankYear(y)
	var musicsList struct {
		musicsList []Model.MusicList
		total      int64
	}
	musicsList.musicsList = data
	musicsList.total = total
	data1, err := json.Marshal(musicsList)
	if err != nil {
		log.Println(err)
	}
	_, err = Config.DBR2.Get(context.Background(), key).Result()
	if err != redis.Nil && err != nil {
		log.Println(err)
	}
	if total == 0 {
		Config.DBR2.Expire(context.Background(), key, 1440*ArticlePass)
		log.Println("年歌单排行使用了旧数据！")
	}
	_, _ = Config.DBR2.Del(context.Background(), key).Result()
	err = Config.DBR2.Set(context.Background(), key, string(data1), 1440*ArticlePass).Err()
	if err != nil {
		log.Println("年听歌单榜缓存失败:", err, "y")
	}
}

func GetACacheMusicRankWeek() ([]Model.Music, int, int64) {
	key := utils.GetCacheNameMusicRank("week")
	data, err := Config.DBR2.Get(context.Background(), key).Result()
	if err != redis.Nil && err != nil {
		log.Println(err)
		return nil, utils.ERROR, 0
	} else if err == redis.Nil || len(data) <= 2 {
		return nil, utils.ERROR, 0
	}
	var musics struct {
		music []Model.Music
		total int64
	}
	json.Unmarshal([]byte(data), &musics)
	return musics.music, utils.SUCCESS, musics.total
}

func UpdateMusicRankWeek() {
	var data []Model.Music
	key := utils.GetCacheNameMusicRank("week")
	y, w1 := time.Now().ISOWeek()
	m := utils.GetCNTimeMonth(time.Now().Month().String())
	data, _, total := Model.GetMusicRankWeek(y, m, w1)
	var musics struct {
		music []Model.Music
		total int64
	}
	musics.music = data
	musics.total = total
	data1, err := json.Marshal(musics)
	if err != nil {
		log.Println(err)
	}
	_, err = Config.DBR2.Get(context.Background(), key).Result()
	if err != redis.Nil && err != nil {
		log.Println(err)
	}
	if total == 0 {
		Config.DBR2.Expire(context.Background(), key, 49*ArticlePass).Err()
		log.Println("周歌曲排行仍使用旧数据!")
	}
	_, _ = Config.DBR2.Del(context.Background(), key).Result()
	err = Config.DBR2.Set(context.Background(), key, string(data1), 49*ArticlePass).Err()
	if err != nil {
		log.Println("周听歌榜缓存失败:", err, "y", "m", m, "w", w1)
	}
}

func GetACacheMusicRankMonth() ([]Model.Music, int, int64) {
	key := utils.GetCacheNameMusicRank("month")
	data, err := Config.DBR2.Get(context.Background(), key).Result()
	if err != redis.Nil && err != nil {
		log.Println(err)
		return nil, utils.ERROR, 0
	} else if err == redis.Nil || len(data) <= 2 {
		return nil, utils.ERROR, 0
	}
	var musics struct {
		music []Model.Music
		total int64
	}
	json.Unmarshal([]byte(data), &musics)
	return musics.music, utils.SUCCESS, musics.total
}

func UpdateMusicRankMonth() {
	var data []Model.Music
	key := utils.GetCacheNameMusicRank("month")
	y, _ := time.Now().ISOWeek()
	m := utils.GetCNTimeMonth(time.Now().Month().String())
	data, _, total := Model.GetMusicRankMonth(y, m)
	var musics struct {
		music []Model.Music
		total int64
	}
	musics.music = data
	musics.total = total
	data1, err := json.Marshal(musics)
	if err != nil {
		log.Println(err)
	}
	_, err = Config.DBR2.Get(context.Background(), key).Result()
	if err != redis.Nil && err != nil {
		log.Println(err)
	}
	if total == 0 {
		Config.DBR2.Expire(context.Background(), key, 193*ArticlePass)
		log.Println("月歌曲排行仍使用旧数据!")
	}
	_, _ = Config.DBR2.Del(context.Background(), key).Result()
	err = Config.DBR2.Set(context.Background(), key, string(data1), 193*ArticlePass).Err()
	if err != nil {
		log.Println("月听歌榜缓存失败:", err, "y", "m", m)
	}
}

func GetACacheMusicRankYear() ([]Model.Music, int, int64) {
	key := utils.GetCacheNameMusicRank("year")
	data, err := Config.DBR2.Get(context.Background(), key).Result()
	if err != redis.Nil && err != nil {
		log.Println(err)
		return nil, utils.ERROR, 0
	} else if err == redis.Nil || len(data) <= 2 {
		return nil, utils.ERROR, 0
	}
	var musics struct {
		music []Model.Music
		total int64
	}
	json.Unmarshal([]byte(data), &musics)
	return musics.music, utils.SUCCESS, musics.total
}

func UpdateMusicRankYear() {
	var data []Model.Music
	key := utils.GetCacheNameMusicRank("year")
	y, _ := time.Now().ISOWeek()
	data, _, total := Model.GetMusicRankYear(y)
	var musics struct {
		music []Model.Music
		total int64
	}
	musics.music = data
	musics.total = total
	data1, err := json.Marshal(musics)
	if err != nil {
		log.Println(err)
	}
	_, err = Config.DBR2.Get(context.Background(), key).Result()
	if err != redis.Nil && err != nil {
		log.Println(err)
	}
	if total == 0 {
		Config.DBR2.Expire(context.Background(), key, 1440*ArticlePass)
		log.Println("年歌曲排行仍使用旧数据!")
	}
	_, _ = Config.DBR2.Del(context.Background(), key).Result()
	err = Config.DBR2.Set(context.Background(), key, string(data1), 1440*ArticlePass).Err()
	if err != nil {
		log.Println("年听歌榜缓存失败:", err, "y")
	}
}

func UpdateCacheGetTheTopicList() ([]Model.Topic, int64) {
	key := utils.GetCacheNameTopic("topic")
	data, err := Config.DBR2.Get(c, key).Result()
	var Topic struct {
		topics []Model.Topic
		total  int64
	}
	if err != nil || len(data) <= 2 {
		data1, total := Model.GetTheTopicList()
		if total == 0 {
			return data1, total
		}
		Topic.topics = data1
		Topic.total = total
		data2, _ := json.Marshal(Topic.topics)
		err = Config.DBR2.Set(c, key, string(data2), 24*ArticlePass).Err()
		if err != nil {
			log.Println("topic列表缓存失败", err)
		}
		return data1, total
	}
	_ = json.Unmarshal([]byte(data), &Topic.topics)
	return Topic.topics, int64(len(Topic.topics))
}
