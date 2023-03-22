package v1

import (
	"GraduationDesign/BaseMent/Cloud/CosCloud"
	"GraduationDesign/BaseMent/Model"
	"GraduationDesign/BaseMent/Model/Cache"
	"GraduationDesign/BaseMent/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func AddArticle(c *gin.Context) {
	var data Model.Article
	var code int
	var msg string
	_ = c.ShouldBind(&data)
	a := c.GetString("user_id")
	if data.UserId != a {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  1015,
			"message": utils.GetErrMsg(1015),
		})
		return
	}
	msg, code = utils.Validate(&Model.Article{UserId: data.UserId})
	if code != utils.SUCCESS {
		c.JSON(http.StatusNotAcceptable, gin.H{
			"status":  code,
			"message": msg,
		})
		c.Abort()
		return
	}
	code, data.ID = Model.CreatArticle(&data)
	if code == utils.SUCCESS {
		c.JSON(http.StatusOK, gin.H{
			"status":  code,
			"id":      data.ID,
			"message": utils.GetErrMsg(code),
		})
		return
	}
	c.JSON(http.StatusInternalServerError, gin.H{
		"status":  code,
		"message": utils.GetErrMsg(code),
	})
}

func EditArticle(c *gin.Context) {
	var data Model.Article
	_ = c.ShouldBind(&data)
	id, _ := strconv.Atoi(c.Query("id"))
	code := Model.CheckArticle(uint(id))
	if code == utils.SUCCESS {
		data, code = Model.EditArticle(uint(id), &data)
		if code == utils.SUCCESS {
			Cache.UpdateCacheArticle(uint(id), &data)
		}
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

func EditArticlePictures(c *gin.Context) {
	ArticlePictureFile, _ := c.FormFile("ArticlePictureFile")
	b := c.PostForm("id")
	articleId, _ := strconv.Atoi(b)
	userId := c.GetString("user_id")
	code := Model.CheckUpUser(userId)
	if articleId == 0 {
		c.JSON(http.StatusNotAcceptable, gin.H{
			"status":  code,
			"message": utils.GetErrMsg(code),
		})
		return
	}
	if code == utils.SUCCESS {
		code = utils.CheckPictureBackgroundIsValidate(ArticlePictureFile)
	}
	if code != utils.SUCCESS {
		c.JSON(http.StatusNotAcceptable, gin.H{
			"status":  code,
			"message": utils.GetErrMsg(code),
		})
		return
	}
	articleId2 := strconv.Itoa(articleId)
	articleUrl, _ := CosCloud.UpLoadArticlePicture(ArticlePictureFile, userId, articleId2)
	code = Model.EditArticlePicture(uint(articleId), articleUrl)
	if code == http.StatusOK {
		c.JSON(http.StatusOK, gin.H{
			"status":  code,
			"url":     articleUrl,
			"message": utils.GetErrMsg(code),
		})
		return
	}
	c.JSON(http.StatusInternalServerError, gin.H{
		"status":  code,
		"message": utils.GetErrMsg(code),
	})
}

func CheckAArticle(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	//article, code := Model.CheckAArticle(uint(id))
	article, code := Cache.GetArticle(uint(id))
	if code == http.StatusOK {
		c.JSON(http.StatusOK, gin.H{
			"status":  code,
			"data":    article,
			"message": utils.GetErrMsg(code),
		})
		return
	}
	c.JSON(http.StatusInternalServerError, gin.H{
		"status":  code,
		"data":    article,
		"message": utils.GetErrMsg(code),
	})
}

func DeleteArticle(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	Cache.DeleteACacheArticle(uint(id))
	code := Model.DeleteArticle(uint(id))
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

func GetTheTopicTopList(c *gin.Context) {
	data, total := Cache.UpdateCacheGetTheTopicList()
	c.JSON(http.StatusOK, gin.H{
		"status":    http.StatusOK,
		"TopicList": data,
		"total":     total,
		"message":   utils.GetErrMsg(http.StatusOK),
	})
}

func LikeArticle(c *gin.Context) {
	articleId, _ := strconv.Atoi(c.Param("id"))
	userId := c.GetString("user_id")
	code := Cache.UpdateCaCheArticleLike(uint(articleId), userId)
	if code == utils.ERROR {
		code = Model.LikeArticle(uint(articleId), userId)
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

func DisLikeArticle(c *gin.Context) {
	articleId, _ := strconv.Atoi(c.Param("id"))
	userId := c.GetString("user_id")
	code := Cache.UpdateCaCheArticleDisLike(uint(articleId), userId)
	if code == utils.ERROR {
		code = Model.DisLikeArticle(uint(articleId), userId)
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

func RecordReadCount(c *gin.Context) {
	articleId, _ := strconv.Atoi(c.Param("id"))
	code := Cache.UpdateCaCheArticleRead(uint(articleId))
	if code == utils.ERROR {
		Model.RecordReadCount(uint(articleId))
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

func CheckArticleLike(c *gin.Context) {
	articleId, _ := strconv.Atoi(c.Param("id"))
	userId := c.GetString("user_id")
	status := Cache.CheckACaCheArticleLike(uint(articleId), userId)
	//status := Model.CheckArticleLike(uint(articleId), userId)
	c.JSON(http.StatusOK, gin.H{
		"status": status,
	})
	return
}

func ForwardArticle(c *gin.Context) {
	articleId, _ := strconv.Atoi(c.Param("id"))
	userId := c.Param("user_id")
	code := Cache.UpdateCaCheArticleForward(uint(articleId), userId)
	if code == utils.ERROR {
		code = Model.ForwardArticle(uint(articleId), userId)
	}
	if code == utils.SUCCESS {
		c.JSON(http.StatusOK, gin.H{
			"status": true,
		})
	}
	c.JSON(http.StatusOK, gin.H{
		"status": false,
	})
	return
}

func GetUserArticles(c *gin.Context) {
	pageSize, _ := strconv.Atoi(c.Query("pageSize"))
	pageNum, _ := strconv.Atoi(c.Query("pageNum"))
	userid := c.Param("user_id")
	//switch {
	//case pageSize > 30:
	//	pageSize = 30
	//case pageSize <= 0:
	//	pageSize = 10
	//}
	//if pageNum == 0 {
	//	pageNum = 1
	//}
	articles, code, total := Model.SearchUserActivities(userid, pageSize, pageNum)
	if code == utils.SUCCESS {
		c.JSON(http.StatusOK, gin.H{
			"status":   code,
			"articles": articles,
			"total":    total,
			"message":  utils.GetErrMsg(code),
		})
		return
	}
	c.JSON(http.StatusInternalServerError, gin.H{
		"status":   code,
		"articles": articles,
		"total":    total,
		"message":  utils.GetErrMsg(code),
	})
}

func GetUserConcernArticles(c *gin.Context) {
	pageSize, _ := strconv.Atoi(c.Query("pageSize"))
	pageNum, _ := strconv.Atoi(c.Query("pageNum"))
	userId := c.PostForm("user_id")
	concernList, _ := c.GetPostFormArray("concerner_list")
	//switch {
	//case pageSize > 30:
	//	pageSize = 30
	//case pageSize <= 0:
	//	pageSize = 10
	//}
	//if pageNum == 0 {
	//	pageNum = 1
	//}
	articles, code, total := Model.SearchUserActivities7Days(userId, pageSize, pageNum)
	for i := 0; i < len(concernList); i++ {
		if concernList[i] == userId {
			continue
		}
		article, _, total1 := Model.SearchUserActivities7Days(concernList[i], pageSize, pageNum)
		articles = append(articles, article...)
		total += total1
	}
	if code == utils.SUCCESS {
		c.JSON(http.StatusOK, gin.H{
			"status":   code,
			"articles": articles,
			"total":    total,
			"message":  utils.GetErrMsg(code),
		})
		return
	}
	c.JSON(http.StatusInternalServerError, gin.H{
		"status":   code,
		"articles": articles,
		"total":    total,
		"message":  utils.GetErrMsg(code),
	})
}

func GetUserActivitiesAndComments(c *gin.Context) {
	pageSize, _ := strconv.Atoi(c.Query("pageSize"))
	pageNum, _ := strconv.Atoi(c.Query("pageNum"))
	userId := c.Param("user_id")
	switch {
	case pageSize <= 0:
		pageSize = 10
	}
	if pageNum == 0 {
		pageNum = 1
	}
	articles, comments, code, total := Model.SearchActivitiesAndComments(userId, pageSize, pageNum)
	if code == utils.SUCCESS {
		c.JSON(http.StatusOK, gin.H{
			"status":   code,
			"articles": articles,
			"comments": comments,
			"total":    total,
			"message":  utils.GetErrMsg(code),
		})
		return
	}
	c.JSON(http.StatusInternalServerError, gin.H{
		"status":   code,
		"articles": articles,
		"comments": comments,
		"total":    total,
		"message":  utils.GetErrMsg(code),
	})
}

func SearchActivitiesAndComments(c *gin.Context) {
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
	articles, comments, code, total := Model.SearchActivitiesAndComments(title, pageSize, pageNum)
	if code == utils.SUCCESS {
		c.JSON(http.StatusOK, gin.H{
			"status":   code,
			"articles": articles,
			"comments": comments,
			"total":    total,
			"message":  utils.GetErrMsg(code),
		})
		return
	}
	c.JSON(http.StatusInternalServerError, gin.H{
		"status":   code,
		"articles": articles,
		"comments": comments,
		"total":    total,
		"message":  utils.GetErrMsg(code),
	})
}

func SearchActivities(c *gin.Context) {
	pageSize, _ := strconv.Atoi(c.Query("pageSize"))
	pageNum, _ := strconv.Atoi(c.Query("pageNum"))
	title := c.Query("title")
	//switch {
	//case pageSize <= 0:
	//	pageSize = 10
	//}
	//if pageNum == 0 {
	//	pageNum = 1
	//}
	articles, code, total := Model.SearchActivities(title, pageSize, pageNum)
	if code == utils.SUCCESS {
		c.JSON(http.StatusOK, gin.H{
			"status":   code,
			"articles": articles,
			"total":    total,
			"message":  utils.GetErrMsg(code),
		})
		return
	}
	c.JSON(http.StatusInternalServerError, gin.H{
		"status":   code,
		"articles": articles,
		"total":    total,
		"message":  utils.GetErrMsg(code),
	})
}

func SearchArticleDays(c *gin.Context) {
	pageSize, _ := strconv.Atoi(c.Query("pageSize"))
	pageNum, _ := strconv.Atoi(c.Query("pageNum"))
	title, _ := strconv.Atoi(c.Query("day"))
	//switch {
	//case pageSize <= 0:
	//	pageSize = 10
	//}
	//if pageNum == 0 {
	//	pageNum = 1
	//}
	if title <= 0 {
		title = 0
	} else if title >= 7 {
		title = 6
	}
	articles, code, total := Model.SearchArticleDays(title, pageSize, pageNum)
	if code == utils.SUCCESS {
		c.JSON(http.StatusOK, gin.H{
			"status":  code,
			"data":    articles,
			"total":   total,
			"message": utils.GetErrMsg(code),
		})
		return
	}
	c.JSON(http.StatusInternalServerError, gin.H{
		"status":  code,
		"data":    articles,
		"total":   total,
		"message": utils.GetErrMsg(code),
	})
}

func SearchTopics(c *gin.Context) {
	pageSize, _ := strconv.Atoi(c.Query("pageSize"))
	pageNum, _ := strconv.Atoi(c.Query("pageNum"))
	title := c.Query("title")
	//switch {
	//case pageSize <= 0:
	//	pageSize = 10
	//}
	//if pageNum == 0 {
	//	pageNum = 1
	//}
	articles, code, total := Model.SearchTopics(title, pageSize, pageNum)
	if code == utils.SUCCESS {
		c.JSON(http.StatusOK, gin.H{
			"status":  code,
			"data":    articles,
			"total":   total,
			"message": utils.GetErrMsg(code),
		})
		return
	}
	c.JSON(http.StatusInternalServerError, gin.H{
		"status":  code,
		"data":    articles,
		"total":   total,
		"message": utils.GetErrMsg(code),
	})
}

func SearchMusic(c *gin.Context) {
	pageSize, _ := strconv.Atoi(c.Query("pageSize"))
	pageNum, _ := strconv.Atoi(c.Query("pageNum"))
	title := c.Query("title")
	//switch {
	//case pageSize <= 0:
	//	pageSize = 10
	//}
	//if pageNum == 0 {
	//	pageNum = 1
	//}
	musics, code, total := Model.SearchMusics(title, pageSize, pageNum)
	if code == utils.SUCCESS {
		c.JSON(http.StatusOK, gin.H{
			"status":  code,
			"data":    musics,
			"total":   total,
			"message": utils.GetErrMsg(code),
		})
		return
	}
	c.JSON(http.StatusInternalServerError, gin.H{
		"status":  code,
		"data":    musics,
		"total":   total,
		"message": utils.GetErrMsg(code),
	})
}
