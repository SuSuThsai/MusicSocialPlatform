package TopRankCache

import (
	"GraduationDesign/BaseMent/Model"
	"GraduationDesign/BaseMent/Model/Cache"
	"GraduationDesign/BaseMent/utils"
)

func InitTopRankCacheBasement() {
	specD := "0 0 0 * * ?"
	specW := "0 10 0 * * 1/2"
	specM := "0 15 0 1/10 * ?"
	specY := "0 0 3 1 1/2 ?"
	utils.ScheduledUpdateTask(Cache.PersistentMusicDayRank, specD)
	utils.ScheduledUpdateTask(Model.DeleteCommandMusicIsListen, specD)
	utils.ScheduledUpdateTask(Cache.PersistentMusicListDayRank, specD)
	utils.ScheduledUpdateTask(Cache.UpdateMusicListRankWeek, specW)
	utils.ScheduledUpdateTask(Cache.UpdateMusicListRankMonth, specM)
	utils.ScheduledUpdateTask(Cache.UpdateMusicListRankYear, specY)
	utils.ScheduledUpdateTask(Cache.UpdateMusicRankWeek, specW)
	utils.ScheduledUpdateTask(Cache.UpdateMusicRankMonth, specM)
	utils.ScheduledUpdateTask(Cache.UpdateMusicRankYear, specY)
	Model.ScheduledArticleTask()
}
