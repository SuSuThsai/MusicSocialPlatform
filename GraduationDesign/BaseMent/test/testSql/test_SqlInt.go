package main

import (
	"GraduationDesign/BaseMent/Config"
	"GraduationDesign/BaseMent/Model"
)

func main() {
	Config.InitsConfig()
	Config.InitsPSQL()
	//InitRedis()
	//Config.DB.AutoMigrate(&Model.User{}, &Model.UserInfo{}, &Model.SongList{}, &Model.MusicLike{}, &Model.MusicList{}, &Model.Music{}, &Model.Album{}, &Model.ConcernList{}, &Model.FollowList{})
	//Config.DB.AutoMigrate(&Model.User{}, &Model.UserInfo{}, &Model.MusicList{}, &Model.SongList{})
	//c := Model.UserInfo{Desc: "娜娜天下第一"}
	//d := Model.ConcernList{ConcernId: "1234"}
	//f := Model.FollowList{FollowerId: "123"}
	//g := Model.MusicLike{MusicId: 1}
	//k := Model.SongList{MusicId: 1}
	//h := Model.MusicList{LName: "开心", SongList: []Model.SongList{k}}
	//j := Model.Music{Name: "娜娜"}
	//i := Model.Album{Name: "摇滚", Music: []Model.Music{j}}
	//u := Model.User{UserId: "123", Username: "犬饲", Password: "123", UserInfo: c, MusicList: []Model.MusicList{h}}
	//u1 := Model.User{UserId: "1234", Username: "犬饲", Password: "123", UserInfo: Model.UserInfo{}}
	//Config.DB.Create(&u)
	//Config.DB.Create(&c)
	//Config.DB.Select(clause.Associations).Delete(&Model.User{Model: gorm.Model{ID: 1}})
	//BaseMent.DB.Create(&u1)
	//BaseMent.DB.Create(&i)

	//Comment,Article
	//Config.DB.AutoMigrate(&Model.Article{}, &Model.Comment{})
	//data, total := Model.GetCommentListByArticleId(1, 1, 0)
	//fmt.Println(data)
	//fmt.Println(total)
	Model.DeleteArticle(1)
}
