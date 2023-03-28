package routes

import (
	v1 "GraduationDesign/BaseMent/api/v1"
	"GraduationDesign/BaseMent/middleware"
	"github.com/gin-gonic/gin"
)

func InitRoutes() {
	//debug 开发模式 release生产模式
	gin.SetMode("debug")
	cod := gin.New()
	cod.Use(middleware.Cors())
	cod.Use(gin.Recovery())
	Basement := cod.Group("")
	{
		login := cod.Group("")
		{
			login.POST("Login", v1.Login)
			login.POST("Adduser", v1.AddUser)
			login.GET("useridisexit/:user_id", v1.CheckUserIdIsExit)
		}
		find := cod.Group("Find")
		{
			find.GET("User_id", v1.CheckUserId)
			find.GET("Username/:username", v1.CheckUserName)
		}
		music := cod.Group("")
		{
			music.POST("MusicAdd", v1.SongAdd)
		}
		//professional
		Basement.GET(":user_id", v1.CheckUserId)
		//Basement.GET(":user_id/dynamic", v1.GetUserActivitiesAndComments)
		Basement.GET(":user_id/media", v1.GetUserArticles)
		Basement.GET(":user_id/music_like", v1.SearchUserMusicsLike)
		Basement.GET(":user_id/musicList_like", v1.SearchUserMusicListLike)
		Basement.GET(":user_id/music_list", v1.GetUserMusicList)
		Basement.GET(":user_id/music_listTips/:list_id", v1.GetUserMusicListTips)
		//search
		//Basement.GET("searchAll", v1.SearchActivitiesAndComments)
		Basement.GET("searchArticle", v1.SearchActivities)
		Basement.GET("GetArticles/days", v1.SearchArticleDays)
		Basement.GET("searchMusic", v1.SearchMusic)
		Basement.GET("searchTopic", v1.SearchTopics)
		//Basement.GET("searchComment", v1.SearchComments)
		Basement.GET("searchMusicList", v1.SearchMusicList)
		//article
		Basement.GET("Article/:article_id", v1.CheckAArticle)
		Basement.GET("GetaArticleCommentList/:article_id", v1.GetCommentListByArticleId)
		Basement.GET("GetaArticleCommentList", v1.GetCommentListByTypeId)
		Basement.POST("Article/record/:id", v1.RecordReadCount)
		//topic
		Basement.GET("TopicList", v1.GetTheTopicTopList)
		//comment
		//Basement.GET("Article/:article_id/Comment/:comment_id", v1.CheckAComment)
		//music
		Basement.POST("Music/listen/:id", v1.MusicDayRankAdd, v1.SongListen)
		Basement.GET("Music/search", v1.SearchAllMusics)
		Basement.GET("Music/Habbty/:id", v1.MusicHabbtyGet)
		Basement.POST("Music/Habbtyadd", v1.MusicHabbtyAdd)
		Basement.POST("Music/music_like_count/:id", v1.MusicLikeCount)
		Basement.GET("Music/music_rank_day", v1.GetMusicRankList)
		Basement.GET("Music/music_rank_week", v1.GetMusicRankWeekList)
		Basement.GET("Music/music_rank_month", v1.GetMusicRankMonthList)
		Basement.GET("Music/music_rank_year", v1.GetMusicRankYearList)
		//musiclist
		Basement.POST("MusicList/listRank/:id", v1.MusicListDayRankAdd)
		Basement.GET("MusicList/search", v1.SearchAllMusicLists)
		Basement.GET("MusicList/Habbty/:id", v1.MusicListHabbtyGet)
		Basement.GET("MusicList/musicList_rank_day", v1.MusicListRankDayList)
		Basement.GET("MusicList/musicList_rank_week", v1.GetMusicListRankWeekList)
		Basement.GET("MusicList/musicList_rank_month", v1.GetMusicListRankMonthList)
		Basement.GET("MusicList/musicList_rank_year", v1.GetMusicListRankYearList)

		//upload
		Basement.POST("upload/MusicL", v1.UploadMusicL)
		Basement.POST("upload/MusicW", v1.UploadMusicW)
	}

	User := cod.Group("").Use(middleware.JwtToken())
	{
		User.POST("Articles/concerns", v1.GetUserConcernArticles)
		User.POST("Music/UserListen/:id", v1.UserSongListen, v1.CountUserTypeListened)
		User.GET("Music/GetListen", v1.GetUserSongListen)
		User.GET("Comments/:user_id", v1.GetUserAllComments)

		//userCurd
		User.POST("account/setting", v1.EditUser)
		User.POST("account/name", v1.EditUserName)
		User.POST("account/ChangePwd", v1.ChangePassword)
		User.POST("account/Delete", v1.DeleteUser)
		User.POST("account/face", v1.EditUserFace)
		User.POST("account/background", v1.EditUserBKGD)
		User.GET("userpasswordexit/:password", v1.CheckUserPassword)
		//article
		User.POST("Article/edit", v1.EditArticle)
		User.POST("Article/create", v1.AddArticle)
		User.POST("Article/uploadPictures", v1.EditArticlePictures)
		User.POST("Article/delete", v1.DeleteArticle)
		User.POST("Article/Like/:id", v1.LikeArticle)
		User.POST("Article/disLike/:id", v1.DisLikeArticle)
		User.GET("Article/checkIsLike/:id", v1.CheckArticleLike)
		User.POST("Article/forward", v1.ForwardArticle)

		//fans fans folows
		User.GET("followers/:user_id", v1.GetFollows)
		User.GET("concerns/:user_id", v1.GetConcern)
		//User.POST("follower_like", v1.UserFollowerLike)
		//User.POST("follower_dislike", v1.UserFollowerDisLike)
		//User.GET("follower_isLike", v1.CheckFollowerLike)
		User.POST("concern_like", v1.UserConcernLike)
		User.POST("concern_dislike", v1.UserConcernDisLike)
		User.GET("concern_isLike/:user_id", v1.CheckConcernLike)

		//Comment
		User.POST("Comment/Create", v1.AddAComment)
		User.POST("Comment/Delete", v1.DeleteComment)
		User.POST("Comment/Like", v1.LikeComment)
		User.POST("Comment/DisLike", v1.DisLikeComment)
		User.GET("Comment/checkIsLike/:id", v1.CheckCommentLike)

		//song
		User.POST("Music/create")
		User.POST("Music/delete")
		User.POST("Music/music_like/:id", v1.UserSongLike)
		User.POST("Music/music_dislike/:id", v1.UserSongDisLike)
		User.GET("Music/music_isLike/:id", v1.CheckSongLike)
		User.GET("Music/professional", v1.GetAUserProfessionalMusics)
		User.GET("Music/professional/days", v1.GetAUserProfessionalMusicsDays)

		//musicList
		User.POST("MusicList/create", v1.AddMusicList)
		User.POST("MusicList/delete/:list_id", v1.DeleteMusicList)
		User.GET("MusicList/GetmusicListSongs/:list_id", v1.GetMusicListSong)
		User.POST("MusicList/AddmusicListSongs", v1.AddMusicListSong)
		User.POST("MusicList/DeletemusicListSongs", v1.DeleteMusicListSong)
		User.POST("MusicList/edit/:id", v1.EditMusicList)
		User.POST("MusicList/editimg/:id", v1.EditUserMusicListImg)
		User.POST("MusicList/musicList_dislike/:id", v1.DisLikeMusicList)
		User.GET("MusicList/musicList_isLike/:id", v1.CheckMusicListLike)
	}
	//port, err := utils.GetFreePort()
	//if err != nil {
	//	log.Panicln("获取端口失败,启用端口5053：", err.Error())
	//}
	cod.Run(":" + "6939")
}
