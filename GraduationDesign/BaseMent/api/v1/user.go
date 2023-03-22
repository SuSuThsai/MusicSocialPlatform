package v1

import (
	"GraduationDesign/BaseMent/Cloud/CosCloud"
	"GraduationDesign/BaseMent/Model"
	"GraduationDesign/BaseMent/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"strings"
)

func AddUser(c *gin.Context) {
	var data Model.User
	var msg string
	var code int
	_ = c.ShouldBind(&data)
	msg, code = utils.Validate(&Model.User{UserId: data.UserId, Username: data.Username, Password: data.Password})
	if code != utils.SUCCESS {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  code,
			"message": msg,
		})
		c.Abort()
		return
	}
	code = Model.CheckUser(data.UserId)
	if code == utils.SUCCESS {
		code = Model.CreatUser(&data)
		c.JSON(http.StatusOK, gin.H{
			"status":  code,
			"data":    data.UserId + " " + data.Username,
			"message": utils.GetErrMsg(code),
		})
		return
	}
	c.JSON(http.StatusInternalServerError, gin.H{
		"status":  code,
		"data":    data.UserId + " " + data.Username,
		"message": utils.GetErrMsg(code),
	})
}

func EditUserFace(c *gin.Context) {
	faceFile, _ := c.FormFile("faceFile")
	userId := c.GetString("user_id")
	code := Model.CheckUpUser(userId)
	if code == utils.SUCCESS {
		code = utils.CheckPicturePFPIsValidate(faceFile)
	}
	if code != utils.SUCCESS {
		c.JSON(http.StatusNotAcceptable, gin.H{
			"status":  code,
			"message": utils.GetErrMsg(code),
		})
		return
	}
	faceUrl, _ := CosCloud.UpLoadFace(faceFile, userId)
	code = Model.EditUserFace(userId, faceUrl)
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

func EditUserBKGD(c *gin.Context) {
	userid := c.GetString("user_id")
	backGroundFile, _ := c.FormFile("backGroundFile")
	code := Model.CheckUpUser(userid)
	if code == utils.SUCCESS {
		code = utils.CheckPictureBackgroundIsValidate(backGroundFile)
	}
	if code != utils.SUCCESS {
		c.JSON(http.StatusNotAcceptable, gin.H{
			"status":  code,
			"message": utils.GetErrMsg(code),
		})
		return
	}
	backgroundUrl, _ := CosCloud.UpLoadBackGround(backGroundFile, userid)
	code = Model.EditUserBKGD(userid, backgroundUrl)
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

func EditUserMusicListImg(c *gin.Context) {
	listId := c.Param("id")
	backGroundFile, _ := c.FormFile("File")
	code := utils.CheckPictureBackgroundIsValidate(backGroundFile)
	if code != utils.SUCCESS {
		c.JSON(http.StatusNotAcceptable, gin.H{
			"status":  code,
			"message": utils.GetErrMsg(code),
		})
		return
	}
	backgroundUrl, _ := CosCloud.UpLoadMusicListImg(backGroundFile, listId)
	listId2, _ := strconv.Atoi(listId)
	code = Model.EditMusicListImg(uint(listId2), backgroundUrl)
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

// EditUserName Edit UserName
func EditUserName(c *gin.Context) {
	var data Model.User
	_ = c.ShouldBind(&data)
	code := Model.CheckUpUser(data.UserId)
	if code == utils.SUCCESS {
		code = Model.EditUserName(data.UserId, &data)
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

// EditUser edit
func EditUser(c *gin.Context) {
	var data Model.UserInfo
	_ = c.ShouldBind(&data)
	user, code := Model.CheckUpUserId(data.UserId)
	a := c.GetString("user_id")
	if user.UserId != a {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  1015,
			"message": utils.GetErrMsg(1015),
		})
		return
	}
	if code == utils.SUCCESS {
		code = Model.EditUser(user.UserId, &data)
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

// DeleteUser delete
func DeleteUser(c *gin.Context) {
	id := c.Param("user_id")
	code := Model.DeleteUser(id)

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

// ChangePassword change-password
func ChangePassword(c *gin.Context) {
	var data Model.User
	_ = c.ShouldBind(&data)
	data.UserId = c.GetString("user_id")
	code := Model.ChangePassword(&data)
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

// CheckUserIdIsExit Check A UserId Is Exit
func CheckUserIdIsExit(c *gin.Context) {
	id := c.Param("user_id")
	code := Model.CheckUser(id)
	c.JSON(http.StatusOK, gin.H{
		"status": code,
	})
}

// CheckUserPassword Check A Password Is Exit
func CheckUserPassword(c *gin.Context) {
	password := c.Param("password")
	id := c.GetString("user_id")
	_, code := Model.ValidateLogin(id, password)
	c.JSON(http.StatusOK, gin.H{
		"status": code,
	})
}

// CheckUserId checkId
func CheckUserId(c *gin.Context) {
	id := c.Param("user_id")
	data, data2, data3, data4, code := Model.GetUser(id)
	if data.ID == 0 {
		code = utils.TargetNotExit
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  code,
			"data":    "",
			"message": utils.GetErrMsg(code),
		})
	}
	var maps = make(map[string]interface{})
	maps["id"] = data.ID
	maps["user_id"] = data.UserId
	maps["username"] = data.Username
	maps["role"] = data.Role
	maps["sex"] = data2.Sex
	maps["desc"] = data2.Desc
	maps["pfp"] = data2.Pfp
	maps["background"] = data2.Background
	maps["concerns"] = data3
	maps["follows"] = data4
	maps["works"], _ = Model.CheckUserHabitty(data.Username)
	if code == http.StatusOK {
		c.JSON(http.StatusOK, gin.H{
			"status":  code,
			"data":    maps,
			"message": utils.GetErrMsg(code),
		})
		return
	}
	c.JSON(http.StatusInternalServerError, gin.H{
		"status":  code,
		"data":    maps,
		"message": utils.GetErrMsg(code),
	})
}

// CheckUserName checkName
func CheckUserName(c *gin.Context) {
	pageSize, _ := strconv.Atoi(c.Query("pageSize"))
	pageNum, _ := strconv.Atoi(c.Query("pageNum"))
	x := strings.Split(c.Param("username"), " ")
	username := x[0]
	a := x[1]
	//switch {
	//case pageSize <= 0:
	//	pageSize = 10
	//}
	//if pageNum == 0 {
	//	pageNum = 1
	//}
	data, data2, total, data3, data4, data5 := Model.GetUsers(username, pageSize, pageNum, a)
	var mapss []map[string]interface{}
	for i, _ := range data {
		var maps = make(map[string]interface{})
		maps["id"] = data[i].ID
		maps["user_id"] = data[i].UserId
		maps["username"] = data[i].Username
		maps["sex"] = data2[i].Sex
		maps["desc"] = data2[i].Desc
		maps["pfp"] = data2[i].Pfp
		maps["background"] = data2[i].Background
		maps["concerns"] = len(data3[i])
		maps["follows"] = len(data4[i])
		maps["isconcern"] = data5[i]
		maps["works"], _ = Model.CheckUserHabitty(data[i].Username)
		mapss = append(mapss, maps)
	}
	c.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
		//"data":    data,
		//"data2":   data2,
		"data":    mapss,
		"total":   total,
		"message": utils.GetErrMsg(http.StatusOK),
	})
}

func GetFollows(c *gin.Context) {
	id := c.Param("user_id")
	follows := Model.GetUserFollow(id)
	var data []Model.User
	var data2 []Model.UserInfo
	var data3 []int
	var data4 []int
	for i := 0; i < len(follows); i++ {
		a, b, x, d, _ := Model.GetUser(follows[i].UserId)
		data = append(data, a)
		data2 = append(data2, b)
		data3 = append(data3, len(x))
		data4 = append(data4, len(d))
	}
	var mapss []map[string]interface{}
	for i, _ := range data {
		var maps = make(map[string]interface{})
		maps["id"] = data[i].ID
		maps["user_id"] = data[i].UserId
		maps["username"] = data[i].Username
		maps["sex"] = data2[i].Sex
		maps["desc"] = data2[i].Desc
		maps["pfp"] = data2[i].Pfp
		maps["background"] = data2[i].Background
		maps["concerns"] = data3[i]
		maps["follows"] = data4[i]
		maps["works"], _ = Model.CheckUserHabitty(data[i].Username)
		mapss = append(mapss, maps)
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"data":    mapss,
		"message": utils.GetErrMsg(http.StatusOK),
	})
}

func GetConcern(c *gin.Context) {
	id := c.Param("user_id")
	concerns := Model.GetUserConcern(id)
	var data []Model.User
	var data2 []Model.UserInfo
	var data3 []int
	var data4 []int
	for i := 0; i < len(concerns); i++ {
		a, b, x, d, _ := Model.GetUser(concerns[i].ConcernId)
		data = append(data, a)
		data2 = append(data2, b)
		data3 = append(data3, len(x))
		data4 = append(data4, len(d))
	}
	var mapss []map[string]interface{}
	for i, _ := range data {
		var maps = make(map[string]interface{})
		maps["id"] = data[i].ID
		maps["user_id"] = data[i].UserId
		maps["username"] = data[i].Username
		maps["sex"] = data2[i].Sex
		maps["desc"] = data2[i].Desc
		maps["pfp"] = data2[i].Pfp
		maps["background"] = data2[i].Background
		maps["concerns"] = data3[i]
		maps["follows"] = data4[i]
		maps["works"], _ = Model.CheckUserHabitty(data[i].Username)
		mapss = append(mapss, maps)
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"data":    mapss,
		"message": utils.GetErrMsg(http.StatusOK),
	})
}

// GetAUserProfessionalMusics Get A UserProfessional Musics
func GetAUserProfessionalMusics(c *gin.Context) {
	userId := c.GetString("user_id")
	data, code := Model.GetCommandMusic(userId)
	if len(data) == 0 || code == utils.ERROR {
		a, _ := Model.CheckUserHabitty(userId)
		var b []string
		for i := 0; i < len(a); i++ {
			b = append(b, a[i].Habits)
		}
		musics, _ := Model.SearchMusicsProfessional(b)
		code = Model.CountCommandMusic(musics, userId)
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  code,
		"data":    data,
		"message": utils.GetErrMsg(code),
	})
}

// GetAUserProfessionalMusicsDays Get A UserProfessional Musics
func GetAUserProfessionalMusicsDays(c *gin.Context) {
	title, _ := strconv.Atoi(c.Query("day"))
	userId := c.GetString("user_id")
	musics, code := Model.GetCommandMusicDays(userId, title)
	c.JSON(http.StatusOK, gin.H{
		"status":  code,
		"data":    musics,
		"message": utils.GetErrMsg(code),
	})
}

// CountUserTypeListened Count User Type Listened
func CountUserTypeListened(c *gin.Context) {
	musicId, _ := strconv.Atoi(c.Param("id"))
	userId := c.GetString("user_id")
	data := Model.GetMusicHabit(uint(musicId), "")
	for i := 0; i < len(data); i++ {
		_ = Model.CountUserTypeListened(userId, data[i].Tip)
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  200,
		"message": utils.GetErrMsg(200),
	})
}
