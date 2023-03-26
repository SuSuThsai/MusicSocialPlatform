package main

import (
	"GraduationDesign/BaseMent/Config"
	"GraduationDesign/BaseMent/Grpc"
	"GraduationDesign/BaseMent/middleware"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	Config.InitsConfig()
	fmt.Println(Config.Conf.FileGrpc, "11111")
	Config.InitFileRpc()
	fmt.Println(Config.FileRpc, "33333")
	gin.SetMode("debug")
	cod := gin.New()
	cod.Use(middleware.Cors())
	cod.Use(gin.Recovery())
	Basement := cod.Group("")
	{
		Basement.POST("GetMusicTips", GetMusicTips)
	}
	cod.Run(":" + "6940")
}

func GetMusicTips(c *gin.Context) {
	file, _ := c.FormFile("file")
	data, data2 := Grpc.TransferFile(file, "2")
	c.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
		"data":   data,
		"data2":  data2,
	})
}
