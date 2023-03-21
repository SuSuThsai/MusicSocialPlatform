package main

import (
	"GraduationDesign/BaseMent/Cloud/CosCloud"
	"GraduationDesign/BaseMent/Config"
	"GraduationDesign/BaseMent/middleware"
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	Config.InitsConfig()
	Config.InitTencentCos()
	gin.SetMode(Config.Conf.SetModel.AppNode)
	cod := gin.New()
	cod.Use(middleware.Cors())
	cod.Use(gin.Recovery())
	cod.Group("")
	{
		cod.POST("upload/:user_id", UploadFace)
	}
	cod.Run(":" + "6929")
}

func UploadFace(c *gin.Context) {
	userId := c.Param("user_id")
	file, _ := c.FormFile("file")
	uel, _ := CosCloud.UpLoadFace(file, userId)
	c.JSON(http.StatusOK, gin.H{
		"data": uel,
	})
}
