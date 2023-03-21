package main

import (
	"GraduationDesign/BaseMent/routes"
)

func main() {
	//Config.InitsSQL()
	//Config.DB.AutoMigrate(&Model.Article{}, &Model.Forward{}, &Model.Topic{}, &Model.ArticleLike{})
	//Config.DB.AutoMigrate(&Model.Comment{}, &Model.CommentLike{})
	//Config.DB.AutoMigrate(&Model.Music{}, &Model.MusicTopic{}, &Model.MusicListen{}, &Model.Album{})
	//Config.DB.AutoMigrate(&Model.MusicRankYear{}, &Model.MusicRankMonth{}, &Model.MusicRankWeek{}, &Model.MusicRankDay{})
	//Config.DB.AutoMigrate(&Model.MusicListRankYear{}, &Model.MusicListRankMonth{}, &Model.MusicListRankWeek{}, &Model.MusicListRankDay{})
	//Config.DB.AutoMigrate(&Model.User{}, &Model.UserInfo{}, Model.UserListeningHabits{}, Model.MusicList{}, Model.Tips{}, Model.SongList{})
	//Config.DB.AutoMigrate(&Model.MusicListLike{}, &Model.MusicLike{})
	//Config.DB.AutoMigrate(&Model.FollowList{}, &Model.ConcernList{})
	//PersistentSql.InitCache()
	//TopRankCache.InitTopRankCacheBasement()
	routes.InitRoutes()
}
