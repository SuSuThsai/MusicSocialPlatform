package v1

import (
	"GraduationDesign/BaseMent/Model"
	"GraduationDesign/BaseMent/middleware"
	"GraduationDesign/BaseMent/utils"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

// Login UP
func Login(c *gin.Context) {
	var data Model.User
	var code int
	var token string
	_ = c.ShouldBind(&data)
	data, code = Model.ValidateLogin(data.UserId, data.Password)
	if code == utils.SUCCESS {
		SetToken(c, data)
	} else {
		c.JSON(http.StatusOK, gin.H{
			"status":  code,
			"user_id": data.UserId,
			"id":      data.ID,
			"token":   token,
			"message": utils.GetErrMsg(code),
		})
	}
}

func SetToken(c *gin.Context, data Model.User) {
	j := middleware.NewJwt()
	claims := middleware.Claims{
		UserId: data.UserId,
		StandardClaims: jwt.StandardClaims{
			NotBefore: time.Now().Unix() - 100,
			//ExpiresAt: time.Now().Unix() + 7200,
			ExpiresAt: time.Now().AddDate(0, 0, 7).Unix(),
			Subject:   "GraduationDesign",
			Issuer:    "Yamada",
		},
	}
	//临时需要加的 现在获取数据有点重复	 后续再优化
	_, data2, _, _, _ := Model.GetUser(data.UserId)
	var maps = make(map[string]interface{})
	maps["id"] = data.ID
	maps["user_id"] = data.UserId
	maps["username"] = data.Username
	maps["role"] = data.Role
	maps["sex"] = data2.Sex
	maps["desc"] = data2.Desc
	maps["pfp"] = data2.Pfp
	maps["background"] = data2.Background
	tokenString, err := j.CreatToken(claims)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  utils.ERROR,
			"message": err,
			"token":   tokenString,
		})
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  utils.SUCCESS,
		"data":    maps,
		"token":   tokenString,
		"message": utils.GetErrMsg(utils.SUCCESS),
	})
	return
}
