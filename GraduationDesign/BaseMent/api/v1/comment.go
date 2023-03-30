package v1

import (
	"GraduationDesign/BaseMent/Model"
	"GraduationDesign/BaseMent/Model/Cache"
	"GraduationDesign/BaseMent/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func AddAComment(c *gin.Context) {
	var data Model.Comment
	var code int
	var msg string
	_ = c.ShouldBind(&data)
	msg, code = utils.Validate(&Model.Comment{UserId: data.UserId})
	if code != utils.SUCCESS {
		c.JSON(http.StatusNotAcceptable, gin.H{
			"status":  code,
			"message": msg,
		})
		c.Abort()
		return
	}
	code = Cache.UpdateCaCheArticleComment(&data)
	if code == utils.ERROR {
		code = Model.CreatAComment(&data)
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

func CheckAComment(c *gin.Context) {
	id, _ := strconv.Atoi(c.Query("id"))
	rId, _ := strconv.Atoi(c.Query("r_id"))
	pageSize, _ := strconv.Atoi(c.Query("pageSize"))
	pageNum, _ := strconv.Atoi(c.Query("pageNum"))
	switch {
	case pageSize > 30:
		pageSize = 30
	case pageSize <= 0:
		pageSize = 10
	}
	if pageNum == 0 {
		pageNum = 1
	}
	comments, _, total := Cache.GetCacheArticleCommentOne(uint(id), uint(rId), pageSize, pageNum)
	//comments, total := Model.GetCommentListInAComment(uint(id), uint(rId), pageSize, pageNum)
	c.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
		"data":   comments,
		"total":  total,
	})
}

func GetCommentListByArticleId(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("article_id"))
	pageSize, _ := strconv.Atoi(c.Query("pageSize"))
	pageNum, _ := strconv.Atoi(c.Query("pageNum"))
	//switch {
	//case pageSize <= 0:
	//	pageSize = 10
	//}
	//if pageNum == 0 {
	//	pageNum = 1
	//}
	comments, _, total := Cache.GetCacheArticleComments(uint(id), pageSize, pageNum)
	//comments, total := Model.GetCommentListByArticleId(uint(id), pageSize, pageNum)
	c.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
		"data":   comments,
		"total":  total,
	})
}

func GetCommentListByTypeId(c *gin.Context) {
	id, _ := strconv.Atoi(c.Query("id"))
	id2, _ := strconv.Atoi(c.Query("type"))
	pageSize, _ := strconv.Atoi(c.Query("pageSize"))
	pageNum, _ := strconv.Atoi(c.Query("pageNum"))
	//switch {
	//case pageSize <= 0:
	//	pageSize = 10
	//}
	//if pageNum == 0 {
	//	pageNum = 1
	//}
	comments, _, total := Cache.GetCacheComments(uint(id), uint(id2), pageSize, pageNum)
	//comments, total := Model.GetCommentListByArticleId(uint(id), pageSize, pageNum)
	c.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
		"data":   comments,
		"total":  total,
	})
}

func GetUserAllComments(c *gin.Context) {
	id := c.Param("user_id")
	pageSize, _ := strconv.Atoi(c.Query("pageSize"))
	pageNum, _ := strconv.Atoi(c.Query("pageNum"))
	switch {
	case pageSize <= 0:
		pageSize = 10
	}
	if pageNum == 0 {
		pageNum = 1
	}
	comments, total := Model.GetUserAllComments(id, pageSize, pageNum)
	c.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
		"data":   comments,
		"total":  total,
	})
}

func DeleteComment(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	code := Model.DeleteComment(uint(id))
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

func LikeComment(c *gin.Context) {
	articleId, _ := strconv.Atoi(c.PostForm("article_id"))
	commentId, _ := strconv.Atoi(c.PostForm("comment_id"))
	userId := c.GetString("user_id")
	code := Cache.UpdateCacheCommentLike(uint(articleId), uint(commentId), userId)
	if code == utils.ERROR {
		code = Model.LikeComment(uint(articleId), uint(commentId), userId)
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

func DisLikeComment(c *gin.Context) {
	articleId, _ := strconv.Atoi(c.PostForm("article_id"))
	commentId, _ := strconv.Atoi(c.PostForm("comment_id"))
	userId := c.GetString("user_id")
	code := Cache.UpdateCacheCommentDisLike(uint(articleId), uint(commentId), userId)
	if code == utils.ERROR {
		code = Model.DisLikeComment(uint(articleId), uint(commentId), userId)
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

func CheckCommentLike(c *gin.Context) {
	commentId, _ := strconv.Atoi(c.Param("id"))
	userId := c.GetString("user_id")
	status := Cache.UpdateCacheCommentIsLike(uint(commentId), userId)
	//status := Model.CheckCommentLike(uint(commentId), userId)
	c.JSON(http.StatusOK, gin.H{
		"status": status,
	})
}

func SearchComments(c *gin.Context) {
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
	comments, code, total := Model.SearchComments(title, pageSize, pageNum)
	if code == utils.SUCCESS {
		c.JSON(http.StatusOK, gin.H{
			"status":   code,
			"articles": comments,
			"total":    total,
			"message":  utils.GetErrMsg(code),
		})
		return
	}
	c.JSON(http.StatusInternalServerError, gin.H{
		"status":   code,
		"articles": comments,
		"total":    total,
		"message":  utils.GetErrMsg(code),
	})
}
