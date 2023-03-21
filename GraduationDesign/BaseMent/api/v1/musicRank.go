package v1

import (
	"GraduationDesign/BaseMent/Model"
	"GraduationDesign/BaseMent/Model/Cache"
	"GraduationDesign/BaseMent/Model/Cache/TopRankCache"
	"GraduationDesign/BaseMent/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

func MusicDayRankAdd(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	code := TopRankCache.MusicRankAdd(uint(id))
	if code == utils.ERROR {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status:":  code,
			"message:": utils.GetErrMsg(code),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status:":  code,
		"message:": utils.GetErrMsg(code),
	})
}

func GetMusicRankList(c *gin.Context) {
	musics, code, total := TopRankCache.GetMusicRankList()
	if code == utils.ERROR {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status:":  code,
			"message:": utils.GetErrMsg(code),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  code,
		"data":    musics,
		"total":   total,
		"message": utils.GetErrMsg(code),
	})
}

func GetMusicRankWeekList(c *gin.Context) {
	musics, code, total := Cache.GetACacheMusicRankWeek()
	if code == utils.ERROR || len(musics) <= 2 {
		y, w1 := time.Now().ISOWeek()
		m := utils.GetCNTimeMonth(time.Now().Month().String())
		musics, code, total = Model.GetMusicRankWeek(y, m, w1)
	}
	if code == utils.ERROR {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  code,
			"message": utils.GetErrMsg(code),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  code,
		"data":    musics,
		"total":   total,
		"message": utils.GetErrMsg(code),
	})
}

func GetMusicRankMonthList(c *gin.Context) {
	musics, code, total := Cache.GetACacheMusicRankMonth()
	if code == utils.ERROR || len(musics) <= 2 {
		y, _ := time.Now().ISOWeek()
		m := utils.GetCNTimeMonth(time.Now().Month().String())
		musics, code, total = Model.GetMusicRankMonth(y, m)
	}
	if code == utils.ERROR {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  code,
			"message": utils.GetErrMsg(code),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  code,
		"data":    musics,
		"total":   total,
		"message": utils.GetErrMsg(code),
	})
}

func GetMusicRankYearList(c *gin.Context) {
	musics, code, total := Cache.GetACacheMusicRankYear()
	if code == utils.ERROR {
		y, _ := time.Now().ISOWeek()
		musics, code, total = Model.GetMusicRankYear(y)
	}
	if code == utils.ERROR {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  code,
			"message": utils.GetErrMsg(code),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  code,
		"data":    musics,
		"total":   total,
		"message": utils.GetErrMsg(code),
	})
}

func MusicListDayRankAdd(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	code := TopRankCache.MusicListRankAdd(uint(id))
	if code == utils.ERROR {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  code,
			"message": utils.GetErrMsg(code),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  code,
		"message": utils.GetErrMsg(code),
	})
}

func MusicListRankDayList(c *gin.Context) {
	Cache.PersistentMusicListDayRank()
	musicList, code, total := TopRankCache.GetMusicListRankList()
	if code == utils.ERROR {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  code,
			"message": utils.GetErrMsg(code),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  code,
		"data":    musicList,
		"total":   total,
		"message": utils.GetErrMsg(code),
	})
}

func GetMusicListRankWeekList(c *gin.Context) {
	musics, code, total := Cache.GetACacheMusicListRankWeek()
	if code == utils.ERROR {
		y, w1 := time.Now().ISOWeek()
		m := utils.GetCNTimeMonth(time.Now().Month().String())
		musics, code, total = Model.GetMusicListRankWeek(y, m, w1)
	}
	if code == utils.ERROR {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  code,
			"message": utils.GetErrMsg(code),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  code,
		"data":    musics,
		"total":   total,
		"message": utils.GetErrMsg(code),
	})
}

func GetMusicListRankMonthList(c *gin.Context) {
	musics, code, total := Cache.GetACacheMusicListRankMonth()
	if code == utils.ERROR {
		y, _ := time.Now().ISOWeek()
		m := utils.GetCNTimeMonth(time.Now().Month().String())
		musics, code, total = Model.GetMusicListRankMonth(y, m)
	}
	if code == utils.ERROR {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  code,
			"message": utils.GetErrMsg(code),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  code,
		"data":    musics,
		"total":   total,
		"message": utils.GetErrMsg(code),
	})
}

func GetMusicListRankYearList(c *gin.Context) {
	musics, code, total := Cache.GetACacheMusicListRankYear()
	if code == utils.ERROR {
		y, _ := time.Now().ISOWeek()
		musics, code, total = Model.GetMusicListRankYear(y)
	}
	if code == utils.ERROR {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  code,
			"message": utils.GetErrMsg(code),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  code,
		"data":    musics,
		"total":   total,
		"message": utils.GetErrMsg(code),
	})
}
