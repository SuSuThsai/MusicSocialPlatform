package main

import (
	"GraduationDesign/BaseMent/Cloud/FtpAndSsh"
	"GraduationDesign/BaseMent/Config"
	"GraduationDesign/BaseMent/middleware"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	Config.InitsConfig()
	fmt.Println(Config.Conf.SFtp)
	Config.InitFtp()
	gin.SetMode(Config.Conf.SetModel.AppNode)
	cod := gin.New()
	cod.Use(middleware.Cors())
	cod.Use(gin.Recovery())
	cod.Group("")
	{
		cod.POST("upload", UploadFileMusicLW)
	}
	cod.Run(":" + "6929")
}

func UploadFileMusicLW(c *gin.Context) {
	file, _ := c.FormFile("file")
	uel := FtpAndSsh.UploadFileMusicL(file)
	c.JSON(http.StatusOK, gin.H{
		"data": uel,
	})
}
