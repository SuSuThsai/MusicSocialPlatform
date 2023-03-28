package Model

import (
	"GraduationDesign/BaseMent/Config"
	"GraduationDesign/BaseMent/utils"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// CheckSongLike Check A Song Is Like or Not
func CheckSongLike(MusicId uint, UserId uint) bool {
	var like MusicLike
	Config.DB.Where("music_id = ? and user_id = ?", MusicId, UserId).First(&like)
	if like.ID == 0 || like.Like == false {
		return false
	}
	return true
}

// UserSongLike User The Songs Like
func UserSongLike(MusicId uint, UserId uint) int {
	var Like MusicLike
	Config.DB.Where("music_id = ? and user_id = ?", MusicId, UserId).First(&Like)
	if Like.ID == 0 {
		Like = MusicLike{MusicId: MusicId, UserId: UserId, Like: true}
		err = Config.DB.Create(&Like).Error
		if err != nil {
			return utils.ERROR
		}
		return utils.SUCCESS
	} else {
		if Like.Like == true {
			return utils.ERROR
		}
		Config.DB.Model(&Like).Update("like", true)
	}
	return utils.SUCCESS
}

// UserSongDisLike User The Songs DisLike
func UserSongDisLike(MusicId uint, userId uint) int {
	var like MusicLike
	Config.DB.Where("user_id = ? and music_id = ?", userId, MusicId).First(&like)
	if like.ID == 0 {
		return utils.SUCCESS
	} else {
		if like.Like == false {
			return utils.ERROR
		}
		err = Config.DB.Model(&like).Update("like", false).Error
		if err != nil {
			return utils.ERROR
		}
	}
	return utils.SUCCESS
}

// CheckFollowerLike Check A Follow Is Like or Not
func CheckFollowerLike(FollowerId string, UserId uint) bool {
	var like FollowList
	Config.DB.Where("follower_id = ? and user_id = ?", FollowerId, UserId).First(&like)
	if like.ID == 0 || like.Like == false {
		return false
	}
	return true
}

// GetUserFollow  GetUser Follow
func GetUserFollow(UserId string) []ConcernList {
	var likes []ConcernList
	Config.DB.Where("concern_id = ?", UserId).Find(&likes)
	var likes1 []ConcernList
	for _, like := range likes {
		if like.Like == true {
			likes1 = append(likes1, like)
		}
	}
	return likes1
}

// GetUserConcern Getuser concern
func GetUserConcern(UserId string) []ConcernList {
	var concerns []ConcernList
	Config.DB.Where("user_id = ?", UserId).Find(&concerns)
	var concerns1 []ConcernList
	for _, concern := range concerns {
		if concern.Like == true {
			concerns1 = append(concerns1, concern)
		}
	}
	return concerns1
}

// UserFollowerLike Follower The User Like
func UserFollowerLike(FollowerId string, UserId uint) int {
	var Like FollowList
	Config.DB.Where("follower_id = ? and user_id = ?", FollowerId, UserId).First(&Like)
	if Like.ID == 0 {
		Like = FollowList{FollowerId: FollowerId, UserId: UserId, Like: true}
		err = Config.DB.Create(&Like).Error
		if err != nil {
			return utils.ERROR
		}
		return utils.SUCCESS
	} else {
		if Like.Like == true {
			return utils.ERROR
		}
		err = Config.DB.Model(&Like).Update("like", true).Error
		if err != nil {
			return utils.ERROR
		}
	}
	return utils.SUCCESS
}

// UserFollowerDisLike Follower The User Like
func UserFollowerDisLike(FollowerId string, UserId uint) int {
	var Like FollowList
	Config.DB.Where("follower_id = ? and user_id = ?", FollowerId, UserId).First(&Like)
	if Like.ID == 0 {
		return utils.SUCCESS
	} else {
		if Like.Like == false {
			return utils.ERROR
		}
		err = Config.DB.Model(&Like).Update("like", false).Error
		if err != nil {
			return utils.ERROR
		}
	}
	return utils.SUCCESS
}

// CheckConcernLike Check A Concern Is Like or Not
func CheckConcernLike(ConcernId string, UserId string) bool {
	var like ConcernList
	Config.DB.Where("concern_id = ? and user_id = ?", ConcernId, UserId).First(&like)
	if like.ID == 0 || like.Like == false {
		return false
	}
	return true
}

// UserConcernLike Concern The User Like
func UserConcernLike(ConcernId string, UserId string) int {
	var Like ConcernList
	Config.DB.Where("concern_id = ? and user_id = ?", ConcernId, UserId).First(&Like)
	if Like.ID == 0 {
		Like = ConcernList{ConcernId: ConcernId, UserId: UserId, Like: true}
		err = Config.DB.Create(&Like).Error
		if err != nil {
			return utils.ERROR
		}
		return utils.SUCCESS
	} else {
		if Like.Like == true {
			return utils.ERROR
		}
		err = Config.DB.Model(&Like).Update("like", true).Error
		if err != nil {
			return utils.ERROR
		}
	}
	return utils.SUCCESS
}

// UserConcernDisLike Concern The User DisLike
func UserConcernDisLike(ConcernId string, UserId string) int {
	var Like ConcernList
	Config.DB.Where("concern_id = ? and user_id = ?", ConcernId, UserId).First(&Like)
	if Like.ID == 0 {
		return utils.SUCCESS
	} else {
		if Like.Like == false {
			return utils.ERROR
		}
		err = Config.DB.Model(&Like).Update("like", false).Error
		if err != nil {
			return utils.ERROR
		}
	}
	return utils.SUCCESS
}

// CheckMusicListLike Check A MusicListLike Is Like or Not
func CheckMusicListLike(MusicListLikeId uint, UserId uint) bool {
	var like MusicListLike
	Config.DB.Where("list_id = ? and user_id = ?", MusicListLikeId, UserId).First(&like)
	if like.ID == 0 || like.Like == false {
		return false
	}
	return true
}

// LikeMusicList Like A MusicList
func LikeMusicList(MusicListId uint, userId uint) int {
	var musicList MusicList
	var like MusicListLike
	Config.DB.Where("user_id = ? and list_id = ?", userId, MusicListId).First(&like)
	if like.ID == 0 {
		like = MusicListLike{ListId: MusicListId, UserId: userId, Like: true}
		err = Config.DB.Create(&like).Error
		if err != nil {
			return utils.ERROR
		}
	} else {
		if like.Like == true {
			return utils.ERROR
		}
		err = Config.DB.Model(&like).Update("like", true).Error
		if err != nil {
			return utils.ERROR
		}
	}
	err = Config.DB.Model(&musicList).Where("id = ?", MusicListId).Update("like_count", gorm.Expr("like_count+ ?", 1)).Error
	if err != nil {
		return utils.ERROR
	}
	return utils.SUCCESS
}

// DisLikeMusicList DisLike A MusicList
func DisLikeMusicList(MusicListId uint, userId uint) int {
	var musicList MusicList
	var like MusicListLike
	Config.DB.Where("user_id = ? and list_id = ?", userId, MusicListId).First(&like)
	if like.ID == 0 {
		return utils.SUCCESS
	} else {
		if like.Like == false {
			return utils.ERROR
		}
		err = Config.DB.Model(&like).Update("like", false).Error
		if err != nil {
			return utils.ERROR
		}
	}
	err = Config.DB.Model(&musicList).Where("id = ?", MusicListId).Update("like_count", gorm.Expr("like_count- ?", 1)).Error
	if err != nil {
		return utils.ERROR
	}
	return utils.SUCCESS
}

// CreatMusicList Creat A MusicList
func CreatMusicList(data *MusicList) (int, uint) {
	var musicList MusicList
	Config.DB.Where("user_id = ? and l_name = ?", data.UserId, data.LName).First(&musicList)
	if musicList.ID > 0 {
		return utils.ERROR, musicList.ID
	}
	err = Config.DB.Create(&data).Error
	if err != nil {
		return utils.ERROR, 0
	}
	return utils.SUCCESS, data.ID
}

// CheckAMusicList Check A MusicList is exit or not
func CheckAMusicList(id uint) int {
	var listId MusicList
	Config.DB.Where("id = ?", id).First(&listId)
	if listId.ID == 0 {
		return utils.ErrorMusicListNotExist
	}
	return utils.SUCCESS
}

// DeleteMusicList Delete A MusicList
func DeleteMusicList(id uint) int {
	var musicList MusicList
	Config.DB.Where("id = ?", id).First(&musicList)
	if musicList.ID == 0 {
		return utils.ErrorMusicListNotExist
	}
	err = Config.DB.Select(clause.Associations).Delete(&musicList).Error
	if err != nil {
		return utils.ERROR
	}
	return utils.SUCCESS
}

// EditMusicList Edit A MusicList
func EditMusicList(id uint, data *MusicList) int {
	var musicList MusicList
	var maps = make(map[string]interface{})
	maps["l_name"] = data.LName
	maps["desc"] = data.Desc
	err = Config.DB.Model(&musicList).Where("id = ?", id).Updates(&maps).Error
	if err != nil {
		return utils.ERROR
	}
	return utils.SUCCESS
}

// SearchMusicLists Search MusicLists
func SearchMusicLists(Name string, pageSize int, pageNum int) ([]MusicList, int, int64) {
	var musicList []MusicList
	var total int64
	err = Config.DB.Order("created_at DESC").Where("l_name LIKE ?",
		"%"+Name+"%").Find(&musicList).Count(&total).Error
	//err = Config.DB.Limit(pageSize).Offset((pageNum-1)*pageSize).Order("created_at DESC").Where("l_name LIKE ?",
	//	"%"+Name+"%").Find(&musicList).Count(&total).Error
	if err != nil {
		return musicList, utils.ERROR, 0
	}
	return musicList, utils.SUCCESS, total
}

// GetUserMusicList Get User MusicLists
func GetUserMusicList(userId uint, pageSize int, pageNum int) ([]MusicList, int, int64) {
	var musicList []MusicList
	var total int64
	err = Config.DB.Order("created_at DESC").Where("user_id = ?",
		userId).Find(&musicList).Count(&total).Error
	//err = Config.DB.Limit(pageSize).Offset((pageNum-1)*pageSize).Order("Created_At DESC").Where("user_id = ?",
	//	userId).Find(&musicList).Count(&total).Error
	if err != nil {
		return musicList, utils.ERROR, 0
	}
	return musicList, utils.SUCCESS, total
}

// GetUserMusicListTips Get Use rMusicList Tips
func GetUserMusicListTips(listId uint) ([]Tips, int, int64) {
	var tips []Tips
	var total int64
	err = Config.DB.Where("list_id = ?",
		listId).Find(&tips).Count(&total).Error
	//err = Config.DB.Limit(pageSize).Offset((pageNum-1)*pageSize).Order("Created_At DESC").Where("user_id = ?",
	//	userId).Find(&musicList).Count(&total).Error
	if err != nil {
		return tips, utils.ERROR, 0
	}
	return tips, utils.SUCCESS, total
}
