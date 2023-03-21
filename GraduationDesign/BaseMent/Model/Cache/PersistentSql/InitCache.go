package PersistentSql

import "GraduationDesign/BaseMent/Config"

var CaCheArticleComment *QueueForTime
var CaCheArticleLike *QueueForTime
var CaCheArticleRead *QueueForTime
var CaCheMusicListen *QueueForTime
var CaCheMusicLike *QueueForTime
var CacheMusicSort *QueueForTime
var CacheMusicListSort *QueueForTime
var CacheCommentLike *QueueForTime
var CacheMusicDayListen *QueueForTime
var CacheMusicListDayListen *QueueForTime
var CacheArticleForward *QueueForTime
var CaCheMusicListLike *QueueForTime

//var TestQueue *QueueForTime
//var TestQueue2 *QueueForTime

func InitCache() {
	CaCheArticleComment = NewQueue(Config.Conf.Cache.CacheArticleCommentCount, Config.DBR, func(payload string) bool {
		return true
	}).WithConcurrent(3)
	CaCheArticleLike = NewQueue(Config.Conf.Cache.CacheArticleLikeCount, Config.DBR, func(payload string) bool {
		return true
	}).WithConcurrent(5)
	CacheArticleForward = NewQueue(Config.Conf.Cache.CacheArticleForward, Config.DBR, func(payload string) bool {
		return true
	}).WithConcurrent(2)
	CaCheArticleRead = NewQueue(Config.Conf.Cache.CacheArticleReadCount, Config.DBR, func(payload string) bool {
		return true
	}).WithConcurrent(6)
	CaCheMusicListen = NewQueue(Config.Conf.Cache.CacheMusicListenCount, Config.DBR, func(payload string) bool {
		return true
	}).WithConcurrent(6)
	CaCheMusicLike = NewQueue(Config.Conf.Cache.CacheMusicLikeCount, Config.DBR, func(payload string) bool {
		return true
	}).WithConcurrent(4)
	CaCheMusicListLike = NewQueue(Config.Conf.Cache.CacheMusicListLikeCount, Config.DBR, func(payload string) bool {
		return true
	}).WithConcurrent(3)
	CacheMusicSort = NewQueue(Config.Conf.Cache.CacheMusicSortCount, Config.DBR, func(payload string) bool {
		return true
	}).WithConcurrent(6)
	CacheMusicListSort = NewQueue(Config.Conf.Cache.CacheMusicListSortCount, Config.DBR, func(payload string) bool {
		return true
	}).WithConcurrent(3)
	CacheCommentLike = NewQueue(Config.Conf.Cache.CacheCommentLikeCount, Config.DBR, func(payload string) bool {
		return true
	}).WithConcurrent(4)
	CacheMusicDayListen = NewQueue(Config.Conf.Cache.CacheMusicDayListen, Config.DBR, func(payload string) bool {
		return true
	}).WithConcurrent(4)
	CacheMusicListDayListen = NewQueue(Config.Conf.Cache.CacheMusicListDayListen, Config.DBR, func(payload string) bool {
		return true
	}).WithConcurrent(4)
	//TestQueue = NewQueue("test-example", Config.DBR, func(payload string) bool {
	//	return true
	//}).WithConcurrent(1)
	//TestQueue2 = NewQueue("test-example", Config.DBR, func(payload string) bool {
	//	return true
	//}).WithConcurrent(1)
	go func() {
		for {
			done := CaCheArticleComment.StartConsume()
			done = CaCheArticleLike.StartConsume()
			done = CacheArticleForward.StartConsume()
			done = CaCheArticleRead.StartConsume()
			done = CaCheMusicListen.StartConsume()
			done = CaCheMusicLike.StartConsume()
			done = CaCheMusicListLike.StartConsume()
			done = CacheMusicSort.StartConsume()
			done = CacheMusicListSort.StartConsume()
			done = CacheCommentLike.StartConsume()
			done = CacheMusicDayListen.StartConsume()
			done = CacheMusicListDayListen.StartConsume()
			//done = TestQueue.StartConsume()
			//done = TestQueue2.StartConsume()
			<-done
		}
	}()
}
