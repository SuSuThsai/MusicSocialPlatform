package TopRankCache

import (
	"GraduationDesign/BaseMent/Config"
	"GraduationDesign/BaseMent/Model"
	"GraduationDesign/BaseMent/Model/Cache"
	"GraduationDesign/BaseMent/utils"
)

func InitTopRankCacheBasement() {
	specD := "0 0 0 * * ?"
	specW := "0 10 0 * * 1/2"
	specM := "0 15 0 1/10 * ?"
	specY := "0 0 3 1 1/2 ?"
	utils.ScheduledUpdateTask(func() {
		Cache.PersistentMusicDayRank()
		Cache.UpdateMusicRankDay()
	}, specD)
	utils.ScheduledUpdateTask(func() {
		Config.InitGlobalUserCommandListen()
		Model.DeleteCommandMusicIsListen()
	}, specD)
	utils.ScheduledUpdateTask(func() {
		Cache.PersistentMusicListDayRank()
		Cache.UpdateMusicListRankDay()
	}, specD)
	utils.ScheduledUpdateTask(func() {
		Cache.UpdateMusicListRankWeek()
	}, specW)
	utils.ScheduledUpdateTask(func() {
		Cache.UpdateMusicListRankMonth()
	}, specM)
	utils.ScheduledUpdateTask(func() {
		Cache.UpdateMusicListRankYear()
	}, specY)
	utils.ScheduledUpdateTask(func() {
		Cache.UpdateMusicRankWeek()
	}, specW)
	utils.ScheduledUpdateTask(func() {
		Cache.UpdateMusicRankMonth()
	}, specM)
	utils.ScheduledUpdateTask(func() {
		Cache.UpdateMusicRankYear()
	}, specY)
	Model.ScheduledArticleTask()
}
