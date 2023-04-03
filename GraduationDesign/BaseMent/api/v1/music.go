package v1

import (
	"GraduationDesign/BaseMent/Cloud/CosCloud"
	"GraduationDesign/BaseMent/Cloud/FtpAndSsh"
	"GraduationDesign/BaseMent/Config"
	"GraduationDesign/BaseMent/Grpc"
	"GraduationDesign/BaseMent/Model"
	"GraduationDesign/BaseMent/Model/Cache"
	"GraduationDesign/BaseMent/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func SongAdd(c *gin.Context) {
	var music Model.Music
	_ = c.ShouldBind(&music)
	code := Model.AddMusic(&music)
	if code == utils.SUCCESS {
		c.JSON(http.StatusOK, gin.H{
			"status":  code,
			"message": utils.GetErrMsg(code),
		})
		return
	}
	c.JSON(http.StatusInternalServerError, gin.H{
		"status":  code,
		"message": utils.GetErrMsg(code),
	})
}

func MusicHabbtyAdd(c *gin.Context) {
	var habbity Model.MusicTopic
	_ = c.ShouldBind(&habbity)
	code := Model.CreatMusicHabit(habbity.Id, habbity.Name, habbity.Tip)
	if code == utils.SUCCESS {
		c.JSON(http.StatusOK, gin.H{
			"status":  code,
			"message": utils.GetErrMsg(code),
		})
		return
	}
	c.JSON(http.StatusInternalServerError, gin.H{
		"status":  code,
		"message": utils.GetErrMsg(code),
	})
}

func MusicHabbtyGet(c *gin.Context) {
	musicId, _ := strconv.Atoi(c.Param("id"))
	data := Model.GetMusicHabit(uint(musicId), "")
	c.JSON(http.StatusOK, gin.H{
		"status":  200,
		"data":    data,
		"message": utils.GetErrMsg(200),
	})
	return
}

func MusicListHabbtyGet(c *gin.Context) {
	listId, _ := strconv.Atoi(c.Param("id"))
	data := Model.GetMusicListHabit(uint(listId))
	c.JSON(http.StatusOK, gin.H{
		"status":  200,
		"data":    data,
		"message": utils.GetErrMsg(200),
	})
	return
}

func SongListen(c *gin.Context) {
	MusicId, _ := strconv.Atoi(c.Param("id"))
	code := Cache.UpdateCacheMusicListened(uint(MusicId))
	if code == utils.ERROR {
		code = Model.MusicsListen(uint(MusicId))
	}
	if code == utils.SUCCESS {
		c.JSON(http.StatusOK, gin.H{
			"status":  code,
			"message": utils.GetErrMsg(code),
		})
		return
	}
	c.JSON(http.StatusInternalServerError, gin.H{
		"status":  code,
		"message": utils.GetErrMsg(code),
	})
}

func GetUserCommandMusicCount(c *gin.Context) {
	MusicId, _ := strconv.Atoi(c.Param("id"))
	userId := c.GetString("user_id")
	code := Model.GetUserCommandMusicCount(userId, uint(MusicId))
	c.JSON(http.StatusOK, gin.H{
		"status":  code,
		"message": utils.GetErrMsg(code),
	})
	return
}

func UserSongListen(c *gin.Context) {
	MusicId, _ := strconv.Atoi(c.Param("id"))
	userId := c.GetString("user_id")
	if Config.GlobalUserCommandListen[userId] != nil && Config.GlobalUserCommandListen[userId][uint(MusicId)] == true {
		Model.GetUserCommandMusicCount(userId, uint(MusicId))
	}
	code := Model.CountUserMusicListened(userId, uint(MusicId))
	if code == utils.SUCCESS {
		c.JSON(http.StatusOK, gin.H{
			"status":  code,
			"message": utils.GetErrMsg(code),
		})
		return
	}
	c.JSON(http.StatusInternalServerError, gin.H{
		"status":  code,
		"message": utils.GetErrMsg(code),
	})
}

func GetUserSongListen(c *gin.Context) {
	userId := c.GetString("user_id")
	data, data2, code := Model.GetUserMusicListened(userId)
	var data3 [][]Model.MusicTopic
	for i := 0; i < len(data2); i++ {
		data3 = append(data3, Model.GetMusicHabit(data2[i].Id, ""))
	}
	if code == utils.SUCCESS {
		c.JSON(http.StatusOK, gin.H{
			"status":  code,
			"data":    data,
			"data2":   data2,
			"data3":   data3,
			"message": utils.GetErrMsg(code),
		})
		return
	}
	c.JSON(http.StatusInternalServerError, gin.H{
		"status":  code,
		"data":    data,
		"data2":   data2,
		"data3":   data3,
		"message": utils.GetErrMsg(code),
	})
}

func MusicLikeCount(c *gin.Context) {
	MusicId, _ := strconv.Atoi(c.Param("id"))
	total := Cache.UpdateMusicLikeCount(uint(MusicId))
	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"id":      MusicId,
		"total":   total,
		"message": utils.GetErrMsg(http.StatusOK),
	})
}

func CheckSongLike(c *gin.Context) {
	MusicId, _ := strconv.Atoi(c.Param("id"))
	userId := c.GetString("user_id")
	a, _, _, _, _ := Model.GetUser(userId)
	status := Cache.CheckCacheMusicIsLike(uint(MusicId), a.ID)
	c.JSON(http.StatusOK, gin.H{
		"status": status,
	})
	return
}

func UserSongLike(c *gin.Context) {
	MusicId, _ := strconv.Atoi(c.Param("id"))
	userId := c.GetString("user_id")
	a, _, _, _, _ := Model.GetUser(userId)
	code := Cache.UpdateCaCheMusicLike(uint(MusicId), a.ID)
	if code == utils.ERROR {
		code = Model.UserSongLike(uint(MusicId), a.ID)
	}
	if code == utils.SUCCESS {
		c.JSON(http.StatusOK, gin.H{
			"status":  code,
			"message": utils.GetErrMsg(code),
		})
		return
	}
	c.JSON(http.StatusInternalServerError, gin.H{
		"status":  code,
		"message": utils.GetErrMsg(code),
	})
}

func UserSongDisLike(c *gin.Context) {
	MusicId, _ := strconv.Atoi(c.Param("id"))
	userId := c.GetString("user_id")
	a, _, _, _, _ := Model.GetUser(userId)
	code := Cache.UpdateCaCheMusicDisLike(uint(MusicId), a.ID)
	if code == utils.ERROR {
		code = Model.UserSongDisLike(uint(MusicId), a.ID)
	}
	if code == utils.SUCCESS {
		c.JSON(http.StatusOK, gin.H{
			"status":  code,
			"message": utils.GetErrMsg(code),
		})
		return
	}
	c.JSON(http.StatusInternalServerError, gin.H{
		"status":  code,
		"message": utils.GetErrMsg(code),
	})
}

func CheckFollowerLike(c *gin.Context) {
	followerId := c.PostForm("user_id")
	a := c.GetString("user_id")
	user, _ := Model.CheckUpUserUserid(a)
	status := Model.CheckFollowerLike(followerId, user.ID)
	c.JSON(http.StatusOK, gin.H{
		"status": status,
	})
	return
}

func UserFollowerLike(c *gin.Context) {
	followerId := c.PostForm("user_id")
	a := c.GetString("user_id")
	user, _ := Model.CheckUpUserUserid(a)
	code := Model.UserFollowerLike(followerId, user.ID)
	if code == utils.SUCCESS {
		c.JSON(http.StatusOK, gin.H{
			"status":  code,
			"message": utils.GetErrMsg(code),
		})
		return
	}
	c.JSON(http.StatusInternalServerError, gin.H{
		"status":  code,
		"message": utils.GetErrMsg(code),
	})
}

func UserFollowerDisLike(c *gin.Context) {
	followerId := c.PostForm("user_id")
	a := c.GetString("user_id")
	user, _ := Model.CheckUpUserUserid(a)
	code := Model.UserFollowerDisLike(followerId, user.ID)
	if code == utils.SUCCESS {
		c.JSON(http.StatusOK, gin.H{
			"status":  code,
			"message": utils.GetErrMsg(code),
		})
		return
	}
	c.JSON(http.StatusInternalServerError, gin.H{
		"status":  code,
		"message": utils.GetErrMsg(code),
	})
}

func CheckConcernLike(c *gin.Context) {
	concernId := c.Param("user_id")
	a := c.GetString("user_id")
	status := Model.CheckConcernLike(concernId, a)
	c.JSON(http.StatusOK, gin.H{
		"status": status,
	})
	return
}

func UserConcernLike(c *gin.Context) {
	concernId := c.PostForm("user_id")
	a := c.GetString("user_id")
	code := Model.UserConcernLike(concernId, a)
	if code == utils.SUCCESS {
		c.JSON(http.StatusOK, gin.H{
			"status":  code,
			"message": utils.GetErrMsg(code),
		})
		return
	}
	c.JSON(http.StatusInternalServerError, gin.H{
		"status":  code,
		"message": utils.GetErrMsg(code),
	})
}

func UserConcernDisLike(c *gin.Context) {
	concernId := c.PostForm("user_id")
	a := c.GetString("user_id")
	code := Model.UserConcernDisLike(concernId, a)
	if code == utils.SUCCESS {
		c.JSON(http.StatusOK, gin.H{
			"status":  code,
			"message": utils.GetErrMsg(code),
		})
		return
	}
	c.JSON(http.StatusInternalServerError, gin.H{
		"status":  code,
		"message": utils.GetErrMsg(code),
	})
}

func CheckMusicListLike(c *gin.Context) {
	MusicListLikeId, _ := strconv.Atoi(c.Param("id"))
	userId := c.GetString("user_id")
	a, _, _, _, _ := Model.GetUser(userId)
	status := Cache.UpdateCacheMusicListIsLike(uint(MusicListLikeId), a.ID)
	//status := Model.CheckMusicListLike(uint(MusicListLikeId), uint(userId))
	c.JSON(http.StatusOK, gin.H{
		"status": status,
	})
	return
}

func LikeMusicList(c *gin.Context) {
	MusicListLikeId, _ := strconv.Atoi(c.Param("id"))
	userId := c.GetString("user_id")
	a, _, _, _, _ := Model.GetUser(userId)
	code := Cache.UpdateCacheMusicListLike(uint(MusicListLikeId), a.ID)
	if code == utils.ERROR {
		code = Model.LikeMusicList(uint(MusicListLikeId), a.ID)
	}
	if code == utils.SUCCESS {
		c.JSON(http.StatusOK, gin.H{
			"status":  code,
			"message": utils.GetErrMsg(code),
		})
		return
	}
	c.JSON(http.StatusInternalServerError, gin.H{
		"status":  code,
		"message": utils.GetErrMsg(code),
	})
}

func DisLikeMusicList(c *gin.Context) {
	MusicListLikeId, _ := strconv.Atoi(c.Param("id"))
	userId := c.GetString("user_id")
	a, _, _, _, _ := Model.GetUser(userId)
	code := Cache.UpdateCacheMusicListDisLike(uint(MusicListLikeId), a.ID)
	if code == utils.ERROR {
		code = Model.DisLikeMusicList(uint(MusicListLikeId), a.ID)
	}
	if code == utils.SUCCESS {
		c.JSON(http.StatusOK, gin.H{
			"status":  code,
			"message": utils.GetErrMsg(code),
		})
		return
	}
	c.JSON(http.StatusInternalServerError, gin.H{
		"status":  code,
		"message": utils.GetErrMsg(code),
	})
}

func AddMusicList(c *gin.Context) {
	var data Model.MusicList
	_ = c.ShouldBind(&data)
	userId := c.GetString("user_id")
	user, _ := Model.CheckUpUserUserid(userId)
	data.UserId = user.ID
	msg, code := utils.Validate(&Model.MusicList{UserId: data.UserId, LName: data.LName})
	if code != utils.SUCCESS {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  code,
			"message": msg,
		})
		c.Abort()
		return
	}
	code, id := Model.CreatMusicList(&data)
	if code == utils.SUCCESS {
		c.JSON(http.StatusOK, gin.H{
			"status":  code,
			"data":    id,
			"message": utils.GetErrMsg(code),
		})
		return
	}
	c.JSON(http.StatusInternalServerError, gin.H{
		"status":  code,
		"data":    id,
		"message": utils.GetErrMsg(code),
	})
}

func EditMusicList(c *gin.Context) {
	var data Model.MusicList
	_ = c.ShouldBind(&data)
	id, _ := strconv.Atoi(c.Param("id"))
	code := Model.CheckAMusicList(uint(id))
	if code == utils.SUCCESS {
		code = Model.EditMusicList(uint(id), &data)
	}
	if code == http.StatusOK {
		c.JSON(http.StatusOK, gin.H{
			"status":  code,
			"message": utils.GetErrMsg(code),
		})
		return
	}
	c.JSON(http.StatusInternalServerError, gin.H{
		"status":  code,
		"message": utils.GetErrMsg(code),
	})
}

func DeleteMusicList(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("list_id"))
	code := Model.DeleteMusicList(uint(id))
	if code == http.StatusOK {
		c.JSON(http.StatusOK, gin.H{
			"status":  code,
			"message": utils.GetErrMsg(code),
		})
		return
	}
	c.JSON(http.StatusInternalServerError, gin.H{
		"status":  code,
		"message": utils.GetErrMsg(code),
	})
}

func SearchAllMusics(c *gin.Context) {
	pageSize, _ := strconv.Atoi(c.Query("pageSize"))
	pageNum, _ := strconv.Atoi(c.Query("pageNum"))
	musics, code, total := Model.SearchAllMusics(pageSize, pageNum)
	var data3 [][]Model.MusicTopic
	for i := 0; i < len(musics); i++ {
		tips := Model.GetMusicHabit(musics[i].Id, "")
		data3 = append(data3, tips)
	}
	if code == utils.ERROR {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": code,
			"data":   musics,
			"data2":  data3,
			"total":  total,
			"msg":    utils.GetErrMsg(code),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status": code,
		"data":   musics,
		"data2":  data3,
		"total":  total,
		"msg":    utils.GetErrMsg(code),
	})
}

func SearchAllMusicLists(c *gin.Context) {
	pageSize, _ := strconv.Atoi(c.Query("pageSize"))
	pageNum, _ := strconv.Atoi(c.Query("pageNum"))
	musics, code, total := Model.SearchAllMusicList(pageSize, pageNum)
	var data3 [][]Model.Tips
	for i := 0; i < len(musics); i++ {
		tips, _, _ := Model.GetUserMusicListTips(musics[i].ID)
		if len(tips) > 5 {
			tips = tips[:5]
		}
		data3 = append(data3, tips)
	}
	if code == utils.ERROR {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": code,
			"data":   musics,
			"data2":  data3,
			"total":  total,
			"msg":    utils.GetErrMsg(code),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status": code,
		"data":   musics,
		"data2":  data3,
		"total":  total,
		"msg":    utils.GetErrMsg(code),
	})
}

func SearchUserMusicsLike(c *gin.Context) {
	userId := c.Param("user_id")
	pageSize, _ := strconv.Atoi(c.Query("pageSize"))
	pageNum, _ := strconv.Atoi(c.Query("pageNum"))
	a, _, _, _, _ := Model.GetUser(userId)
	musics, musicData, code, total := Model.SearchUserMusicsLike(a.ID, pageSize, pageNum)
	var data3 [][]Model.MusicTopic
	for i := 0; i < len(musicData); i++ {
		data3 = append(data3, Model.GetMusicHabit(musicData[i].Id, ""))
	}
	if code == utils.ERROR {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": code,
			"data":   musics,
			"data2":  musicData,
			"data3":  data3,
			"total":  total,
			"msg":    utils.GetErrMsg(code),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status": code,
		"data":   musics,
		"data2":  musicData,
		"data3":  data3,
		"total":  total,
		"msg":    utils.GetErrMsg(code),
	})
}

func SearchUserMusicListLike(c *gin.Context) {
	userId := c.Param("user_id")
	pageSize, _ := strconv.Atoi(c.Query("pageSize"))
	pageNum, _ := strconv.Atoi(c.Query("pageNum"))
	a, _, _, _, _ := Model.GetUser(userId)
	musics, music2, code, total := Model.SearchUserMusicsListLike(a.ID, pageSize, pageNum)
	var data3 [][]Model.Tips
	for i := 0; i < len(musics); i++ {
		tips, _, _ := Model.GetUserMusicListTips(musics[i].ID)
		if len(tips) > 5 {
			tips = tips[:5]
		}
		data3 = append(data3, tips)
	}
	if code == utils.ERROR {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": code,
			"data":   musics,
			"data2":  music2,
			"data3":  data3,
			"total":  total,
			"msg":    utils.GetErrMsg(code),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status": code,
		"data":   musics,
		"data2":  music2,
		"data3":  data3,
		"total":  total,
		"msg":    utils.GetErrMsg(code),
	})
}

func SearchMusicList(c *gin.Context) {
	pageSize, _ := strconv.Atoi(c.Query("pageSize"))
	pageNum, _ := strconv.Atoi(c.Query("pageNum"))
	title := c.Query("title")
	switch {
	case pageSize <= 0:
		pageSize = 10
	}
	if pageNum == 0 {
		pageNum = 1
	}
	musicList, code, total := Model.SearchMusicLists(title, pageSize, pageNum)
	var data3 [][]Model.Tips
	for i := 0; i < len(musicList); i++ {
		tips, _, _ := Model.GetUserMusicListTips(musicList[i].ID)
		if len(tips) > 5 {
			tips = tips[:5]
		}
		data3 = append(data3, tips)
	}
	if code == utils.SUCCESS {
		c.JSON(http.StatusOK, gin.H{
			"status":  code,
			"data":    musicList,
			"data2":   data3,
			"total":   total,
			"message": utils.GetErrMsg(code),
		})
		return
	}
	c.JSON(http.StatusInternalServerError, gin.H{
		"status":  code,
		"data":    musicList,
		"data2":   data3,
		"total":   total,
		"message": utils.GetErrMsg(code),
	})
}

func GetUserMusicListTips(c *gin.Context) {
	listId, _ := strconv.Atoi(c.Param("list_id"))
	tips, code, total := Model.GetUserMusicListTips(uint(listId))
	if code == utils.SUCCESS {
		c.JSON(http.StatusOK, gin.H{
			"status":  code,
			"data":    tips,
			"total":   total,
			"message": utils.GetErrMsg(code),
		})
		return
	}
	c.JSON(http.StatusInternalServerError, gin.H{
		"status":  code,
		"data":    tips,
		"total":   total,
		"message": utils.GetErrMsg(code),
	})
}

func GetUserMusicList(c *gin.Context) {
	pageSize, _ := strconv.Atoi(c.Query("pageSize"))
	pageNum, _ := strconv.Atoi(c.Query("pageNum"))
	UserId := c.Param("user_id")
	switch {
	case pageSize <= 0:
		pageSize = 10
	}
	if pageNum == 0 {
		pageNum = 1
	}
	a, _, _, _, _ := Model.GetUser(UserId)
	musicList, code, total := Model.GetUserMusicList(a.ID, pageSize, pageNum)
	var data3 [][]Model.Tips
	for i := 0; i < len(musicList); i++ {
		tips, _, _ := Model.GetUserMusicListTips(musicList[i].ID)
		if len(tips) > 5 {
			tips = tips[:5]
		}
		data3 = append(data3, tips)
	}
	if code == utils.SUCCESS {
		c.JSON(http.StatusOK, gin.H{
			"status":  code,
			"data":    musicList,
			"data2":   data3,
			"total":   total,
			"message": utils.GetErrMsg(code),
		})
		return
	}
	c.JSON(http.StatusInternalServerError, gin.H{
		"status":  code,
		"data":    musicList,
		"total":   total,
		"message": utils.GetErrMsg(code),
	})
}

func UploadMusicL(c *gin.Context) {
	file, _ := c.FormFile("file")
	code := utils.CheckMusicLIsValidate(file)
	if code != utils.SUCCESS {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  code,
			"message": utils.GetErrMsg(code),
		})
	}
	//Todo 完成grpc的传送结汇对应标签
	result, _ := Grpc.TransferFile(file, "1")
	url, _ := CosCloud.UpLoadMusicL(file)
	code = FtpAndSsh.UploadFileMusicL(file)
	c.JSON(http.StatusOK, gin.H{
		"status":  code,
		"data":    result,
		"url":     url,
		"message": utils.GetErrMsg(code),
	})
}

func UploadMusicW(c *gin.Context) {
	file, _ := c.FormFile("file")
	code := utils.CheckMusicWIsValidate(file)
	if code != utils.SUCCESS {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  code,
			"message": utils.GetErrMsg(code),
		})
	}
	url, _ := CosCloud.UpLoadMusicW(file)
	//Todo 完成grpc的传送结汇对应标签
	result, _ := Grpc.TransferFile(file, "2")
	code = FtpAndSsh.UploadFileMusicW(file)
	c.JSON(http.StatusOK, gin.H{
		"status":  code,
		"data":    result,
		"url":     url,
		"message": utils.GetErrMsg(code),
	})
}

func AddMusicListSong(c *gin.Context) {
	listId, _ := strconv.Atoi(c.PostForm("list_id"))
	musicId, _ := strconv.Atoi(c.PostForm("music_id"))
	userId := c.GetString("user_id")
	user, _ := Model.CheckUpUserUserid(userId)
	code := Model.AddMusicListSong(uint(listId), uint(musicId), user.ID)
	c.JSON(http.StatusOK, gin.H{
		"status":  code,
		"message": utils.GetErrMsg(code),
	})
}

func DeleteMusicListSong(c *gin.Context) {
	listId, _ := strconv.Atoi(c.PostForm("list_id"))
	musicId, _ := strconv.Atoi(c.PostForm("music_id"))
	userId := c.GetString("user_id")
	user, _ := Model.CheckUpUserUserid(userId)
	code := Model.DeleteMusicListSong(uint(listId), uint(musicId), user.ID)
	c.JSON(http.StatusOK, gin.H{
		"status":  code,
		"message": utils.GetErrMsg(code),
	})
}

func GetMusicListSong(c *gin.Context) {
	listId, _ := strconv.Atoi(c.Param("list_id"))
	data, code := Model.GetMusicListSong(uint(listId))
	var data3 [][]Model.MusicTopic
	for i := 0; i < len(data); i++ {
		data3 = append(data3, Model.GetMusicHabit(data[i].Id, ""))
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  code,
		"data":    data,
		"data2":   data3,
		"message": utils.GetErrMsg(code),
	})
}
