package main

import (
	"GraduationDesign/BaseMent/Config"
	"GraduationDesign/BaseMent/Model"
	"GraduationDesign/BaseMent/Model/Cache/PersistentSql"
	"GraduationDesign/BaseMent/Model/Cache/TopRankCache"
	"GraduationDesign/BaseMent/routes"
)

func main() {
	Config.InitsSQL()
	//Config.DB.AutoMigrate(&Model.Article{}, &Model.Forward{}, &Model.Topic{}, &Model.ArticleLike{}, &Model.UserListeningHabits{}, &Model.UserListenTypeCount{}, &Model.UserListenMusicCount{})
	//Config.DB.AutoMigrate(&Model.Comment{}, &Model.CommentLike{})
	//Config.DB.AutoMigrate(&Model.Music{}, &Model.MusicTopic{}, &Model.MusicListen{}, &Model.Album{})
	//Config.DB.AutoMigrate(&Model.MusicRankYear{}, &Model.MusicRankMonth{}, &Model.MusicRankWeek{}, &Model.MusicRankDay{})
	//Config.DB.AutoMigrate(&Model.MusicListRankYear{}, &Model.MusicListRankMonth{}, &Model.MusicListRankWeek{}, &Model.MusicListRankDay{})
	//Config.DB.AutoMigrate(&Model.User{}, &Model.UserInfo{}, Model.UserListeningHabits{}, Model.MusicList{}, Model.Tips{}, Model.SongList{})
	//Config.DB.AutoMigrate(&Model.MusicListLike{}, &Model.MusicLike{})
	//Config.DB.AutoMigrate(&Model.FollowList{}, &Model.ConcernList{},)
	Config.DB.AutoMigrate(&Model.Article{}, &Model.Comment{}, &Model.MusicTopic{}, &Model.CommandMusicCount{}, &Model.CommandMusicListenCount{}, &Model.MusicCommandIsListen{})
	PersistentSql.InitCache()
	TopRankCache.InitTopRankCacheBasement()
	routes.InitRoutes()
}
