package Model

import (
	"GraduationDesign/BaseMent/Config"
	"GraduationDesign/BaseMent/utils"
	"encoding/base64"
	"fmt"
	"golang.org/x/crypto/scrypt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"log"
	"strconv"
	"time"
)

var err error

// CheckUser CheckExit
func CheckUser(id string) (code int) {
	var user User
	Config.DB.Select("id").Where("user_id = ?", id).First(&user)
	if user.ID > 0 {
		return utils.ErrorUsernameUsed //1001
	}
	return utils.SUCCESS
}

// CheckUpUser UpdateCheckExit
func CheckUpUser(id string) (code int) {
	var user User
	Config.DB.Select("id,user_id").Where("user_id = ?", id).First(&user)
	if user.ID <= 0 || user.UserId == "" {
		return utils.ErrorUserNotExist //10003
	}
	return utils.SUCCESS
}

// CheckUpUserId CheckUpUserId
func CheckUpUserId(id uint) (user User, code int) {
	Config.DB.Select("id,user_id").Where("id = ?", id).First(&user)
	if user.ID <= 0 || user.UserId == "" {
		return user, utils.ErrorUserNotExist //10003
	}
	return user, utils.SUCCESS
}

// CheckUpUserUserid CheckUpUserUserid
func CheckUpUserUserid(id string) (user User, code int) {
	Config.DB.Where("user_id = ?", id).First(&user)
	if user.ID <= 0 || user.UserId == "" {
		return user, utils.ErrorUserNotExist //10003
	}
	return user, utils.SUCCESS
}

func CountUserTypeListened(userId string, habiity string) int {
	var data UserListenTypeCount
	err = Config.DB.Where("user_id = ? and habits = ?", userId, habiity).Find(&data).Error
	if err != nil || data.UserId == "" || err == gorm.ErrRecordNotFound {
		data.UserId = userId
		data.Habits = habiity
		data.ListenCount = 0
		err = Config.DB.Create(&data).Model(&UserListenTypeCount{}).Error
		if err != nil {
			return utils.ERROR
		}
		err = Config.DB.Model(&UserListenTypeCount{}).Where("user_id = ? and habits = ?", userId, habiity).Update("listen_count", gorm.Expr("listen_count+ ?", 1)).Error
		return utils.SUCCESS
	}
	Config.DB.Model(&UserListenTypeCount{}).Where("user_id = ? and habits = ?", userId, habiity).Update("listen_count", gorm.Expr("listen_count+ ?", 1))
	return utils.SUCCESS
}

func AddMusicListSong(listId uint, musicId uint, userId uint) int {
	var musicList SongList
	err = Config.DB.Where("list_id = ? and music_id = ?", listId, musicId).First(&musicList).Error
	if musicList.ID == 0 {
		musicList.ListId = listId
		musicList.MusicId = musicId
		musicList.UserId = userId
		err = Config.DB.Create(&musicList).Model(&SongList{}).Error
		mTopic := GetMusicHabit(musicId, "")
		for i := 0; i < len(mTopic); i++ {
			var listTopic Tips
			err = Config.DB.Where("list_id = ? and tip_name = ?", listId, mTopic[i].Tip).Find(&listTopic).Error
			if err == gorm.ErrRecordNotFound || listTopic.ListId == 0 {
				listTopic.ListId = listId
				listTopic.TipName = mTopic[i].Tip
				Config.DB.Create(&listTopic).Model(&Tips{})
			}
		}
	}
	return utils.SUCCESS
}

func DeleteMusicListSong(listId uint, musicId uint, userId uint) int {
	var musicList SongList
	err = Config.DB.Where("list_id = ? and music_id = ?", listId, musicId).First(&musicList).Error
	if musicList.ID == 0 {
		return utils.SUCCESS
	} else {
		Config.DB.Where("list_id = ? and music_id = ?", listId, musicId).Delete(&musicList)
	}
	return utils.SUCCESS
}

func GetMusicListSong(listId uint) ([]Music, int) {
	var musicList []SongList
	err = Config.DB.Where("list_id = ?", listId).Find(&musicList).Error
	var music []Music
	for i := 0; i < len(musicList); i++ {
		var data Music
		data, _ = GetAMusic(musicList[i].MusicId)
		music = append(music, data)
	}
	return music, utils.SUCCESS
}

func CountUserMusicListened(userId string, musicId uint) int {
	var data UserListenMusicCount
	err = Config.DB.Where("user_id = ? and music_id = ?", userId, musicId).Find(&data).Error
	if err != nil || data.UserId == "" || err == gorm.ErrRecordNotFound {
		data.UserId = userId
		data.MusicId = musicId
		data.ListenCount = 0
		err = Config.DB.Create(&data).Model(&UserListenMusicCount{}).Error
		if err != nil {
			return utils.ERROR
		}
		err = Config.DB.Model(&UserListenMusicCount{}).Where("user_id = ? and music_id = ?", userId, musicId).Update("listen_count", gorm.Expr("listen_count+ ?", 1)).Error
		return utils.SUCCESS
	}
	Config.DB.Model(&UserListenMusicCount{}).Where("user_id = ? and music_id = ?", userId, musicId).Update("listen_count", gorm.Expr("listen_count+ ?", 1))
	return utils.SUCCESS
}

func DeleteCommandMusicIsListen() {
	Config.DB.Delete(&MusicCommandIsListen{})
}

func GetUserCommandMusicCount(userId string, musicId uint) int {
	var data2 MusicCommandIsListen
	var data CommandMusicListenCount
	m, d, n := time.Now().Month().String(), time.Now().Weekday().String(), strconv.Itoa(time.Now().Day())
	y1, w1 := time.Now().ISOWeek()
	w := strconv.Itoa(w1)
	y := strconv.Itoa(y1)
	err = Config.DB.Where("user_id = ? and music_id = ?", userId, musicId).Find(&data2).Error
	if data2.IsListen == 0 || err == gorm.ErrRecordNotFound {
		data2.UserId = userId
		data2.MusicId = musicId
		data2.IsListen = 1
		Config.DB.Create(&data2)
		err = Config.DB.Where("user_id = ? and year = ? and month = ? and week = ? and day = ?", userId, y, m, w, d).Find(&data).Error
		fmt.Println("222222", data.UserId)
		if data.UserId == "" || err == gorm.ErrRecordNotFound {
			data.UserId = userId
			data.Year = y
			data.Week = w
			data.Month = m
			data.Week = w
			data.Day = d
			data.Number = n
			data.IsListen = 1
			Config.DB.Create(&data)
		} else {
			Config.DB.Model(&data).Where("user_id = ? and year = ? and month = ? and week = ? and day = ?", userId, y, m, w, d).Update("is_listen", gorm.Expr("is_listen+ ?", 1))
		}
	}
	return utils.SUCCESS
}

func GetUserMusicListened(userId string) ([]UserListenMusicCount, []Music, int) {
	var data []UserListenMusicCount
	var data2 []Music
	err = Config.DB.Where("user_id = ?", userId).Limit(30).Order("updated_at DESC").Find(&data).Error
	for i := 0; i < len(data); i++ {
		var data3 Music
		Config.DB.Where("id = ?", data[i].MusicId).First(&data3)
		data2 = append(data2, data3)
	}
	//err = Config.DB.Where("user_id = ?", userId).Find(&data).Error
	return data, data2, utils.SUCCESS
}

func GetUserTypeListened(userId string) []UserListenTypeCount {
	var data []UserListenTypeCount
	err = Config.DB.Where("user_id = ?", userId).Order("listen_count DESC").Find(&data).Error
	return data
}

// CheckUserHabitty Check User Habitty
func CheckUserHabitty(id string) ([]UserListeningHabits, int) {
	var data []UserListeningHabits
	Config.DB.Where("user_id = ?", id).Find(&data)
	if len(data) >= 5 {
		return data[:5], utils.SUCCESS
	}
	return data, utils.SUCCESS
}

// UpdateUserHabitty UpdateUserHabitty
func UpdateUserHabitty(id string) ([]UserListeningHabits, int) {
	var data []UserListeningHabits
	Config.DB.Where("user_id = ?", id).Find(&data)
	Config.DB.Delete("user_id = ?", id).Model(&UserListeningHabits{})
	//data2 := GetUserTypeListened(id)
	//todo 接受返回数据重新创建用户
	return data, utils.SUCCESS
}

// GetUser FindUserById
func GetUser(id string) (User, UserInfo, []ConcernList, []ConcernList, int) {
	var user User
	var userInfo UserInfo
	var concerns []ConcernList
	var follows []ConcernList
	err = Config.DB.Limit(1).Where("user_id = ?", id).Find(&user).Error
	if err != nil {
		return user, userInfo, concerns, follows, utils.ERROR
	}
	err = Config.DB.Limit(1).Where("user_id = ?", user.ID).Find(&userInfo).Error
	if err != nil {
		return user, userInfo, concerns, follows, utils.ERROR
	}
	concerns = GetUserConcern(user.UserId)
	follows = GetUserFollow(user.UserId)
	return user, userInfo, concerns, follows, utils.SUCCESS
}

// GetUsers FindUserByName
func GetUsers(username string, pageSize int, pageNum int, userid string) ([]User, []UserInfo, int64, [][]ConcernList, [][]ConcernList, []bool) {
	var users []User
	var userInfos []UserInfo
	var concerns [][]ConcernList
	var follows [][]ConcernList
	var isconcern []bool
	var total int64
	if username != "" {
		Config.DB.Model(&users).Where(
			"username LIKE ?", "%"+username+"%",
		).Count(&total)
		//Config.DB.Where(
		//	"username LIKE ?", "%"+username+"%",
		//).Limit(pageSize).Offset((pageNum - 1) * pageSize).Find(&users)
		Config.DB.Where(
			"username LIKE ?", "%"+username+"%",
		).Find(&users)
		userInfos = make([]UserInfo, len(users))
		concerns = make([][]ConcernList, len(users))
		follows = make([][]ConcernList, len(users))
		isconcern = make([]bool, len(users))
		for i, user := range users {
			var userinfo UserInfo
			Config.DB.Limit(1).Where(
				"user_id = ?", user.ID).Find(&userinfo)
			userInfos[i] = userinfo
			concerns[i] = GetUserConcern(user.UserId)
			follows[i] = GetUserFollow(user.UserId)
			isconcern[i] = CheckConcernLike(user.UserId, userid)
		}
		return users, userInfos, total, concerns, follows, isconcern
	}
	Config.DB.Model(&users).Count(&total)
	err = Config.DB.Find(&users).Error
	//err = Config.DB.Limit(pageSize).Offset((pageNum - 1) * pageSize).Find(&users).Error
	userInfos = make([]UserInfo, len(users))
	concerns = make([][]ConcernList, len(users))
	follows = make([][]ConcernList, len(users))
	isconcern = make([]bool, len(users))
	for i, user := range users {
		var userinfo UserInfo
		Config.DB.Select("pfp,sex").Limit(1).Where(
			"user_id", user.ID).Find(&userinfo)
		userInfos[i] = userinfo
		concerns[i] = GetUserConcern(user.UserId)
		follows[i] = GetUserFollow(user.UserId)
		isconcern[i] = CheckConcernLike(user.UserId, userid)
	}
	if err != nil {
		return users, userInfos, 0, concerns, follows, isconcern
	}
	return users, userInfos, total, concerns, follows, isconcern
}

// DeleteUser deleteUser
func DeleteUser(id string) int {
	var user User
	Config.DB.Where("user_id = ?", id).First(&user)
	if user.ID == 0 {
		return utils.ErrorUserNotExist
	}
	err = Config.DB.Select(clause.Associations).Delete(&user).Error
	//对不存在id没有进行处理
	if err != nil {
		return utils.ERROR
	}
	return utils.SUCCESS
}

// EditUserFace Edit userFace
func EditUserFace(id string, Pfp string) int {
	var user User
	var userInfo UserInfo
	var maps = make(map[string]interface{})
	Config.DB.Select("id").Where("user_id = ?", id).First(&user)
	if err != nil {
		return utils.ERROR
	}
	maps["pfp"] = Pfp
	CheckUserInfo(user.ID)
	err = Config.DB.Model(&userInfo).Where("user_id = ?", user.ID).Updates(&maps).Error
	if err != nil {
		return utils.ERROR
	}
	return utils.SUCCESS
}

// EditUserBKGD Edit UserBKGD
func EditUserBKGD(id string, BKGD string) int {
	var user User
	var userInfo UserInfo
	var maps = make(map[string]interface{})
	Config.DB.Select("id").Where("user_id = ?", id).First(&user)
	if err != nil {
		return utils.ERROR
	}
	maps["background"] = BKGD
	CheckUserInfo(user.ID)
	err = Config.DB.Model(&userInfo).Where("user_id = ?", user.ID).Updates(&maps).Error
	if err != nil {
		return utils.ERROR
	}
	return utils.SUCCESS
}

// EditMusicListImg Edit MusicList Img
func EditMusicListImg(listId uint, img string) int {
	var musicList MusicList
	var maps = make(map[string]interface{})
	Config.DB.Where("id = ?", listId).First(&musicList)
	if err != nil || musicList.ID == 0 {
		return utils.ERROR
	}
	maps["img"] = img
	err = Config.DB.Model(&musicList).Where("id = ?", musicList.ID).Updates(&maps).Error
	if err != nil {
		return utils.ERROR
	}
	return utils.SUCCESS
}

// EditUserName Edit UserName
func EditUserName(id string, data *User) int {
	var user User
	var maps = make(map[string]interface{})
	maps["username"] = data.Username
	err = Config.DB.Model(&user).Where("user_id = ?", id).Updates(&maps).Error
	if err != nil {
		return utils.ERROR
	}
	return utils.SUCCESS
}

// EditUser Edit user
func EditUser(id string, data2 *UserInfo) int {
	var user User
	var userInfo UserInfo
	var maps = make(map[string]interface{})
	Config.DB.Select("id").Where("user_id = ?", id).First(&user)
	if err != nil {
		return utils.ERROR
	}
	maps["sex"] = data2.Sex
	maps["desc"] = data2.Desc
	CheckUserInfo(user.ID)
	err = Config.DB.Model(&userInfo).Where("user_id = ?", user.ID).Updates(&maps).Error
	if err != nil {
		return utils.ERROR
	}
	return utils.SUCCESS
}

// CheckUserInfo check
func CheckUserInfo(id uint) bool {
	var userInfo UserInfo
	Config.DB.Select("id,user_id").Where("user_id = ?", id).First(&userInfo)
	if userInfo.ID == 0 {
		userInfo.UserId = id
		Config.DB.Create(&userInfo)
		return false
	}
	return true
}

// ChangePassword changePassword
func ChangePassword(data *User) int {
	err = Config.DB.Select("password").Where("user_id = ?", data.UserId).Updates(&data).Error
	if err != nil {
		return utils.ERROR
	}
	return utils.SUCCESS
}

// CreatUser Create
func CreatUser(data *User) (code int) {
	//data.Password = ScryptPW(data.UserId + data.Password)
	err = Config.DB.Create(&data).Error
	if err != nil {
		return utils.ERROR //500
	}
	return utils.SUCCESS
}

// BeforeCreate Password encryption methods：bcrypt，scrypt and salt hash
func (u *User) BeforeCreate(_ *gorm.DB) (err error) {
	u.Password = ScryptPW(u.UserId + u.Password)
	u.Role = 2
	return nil
}

func (u *User) BeforeUpdate(_ *gorm.DB) (err error) {
	u.Password = ScryptPW(u.UserId + u.Password)
	return nil
}

func ScryptPW(password string) string {
	const KeyLen = 10
	salt := make([]byte, 8)
	salt = []byte{255, 22, 18, 33, 99, 66, 25, 11}
	HashPW, err := scrypt.Key([]byte(password), salt, 16384, 8, 1, KeyLen)
	//HashPW, err :=bcrypt.GenerateFromPassword([]byte(password),KeyLen)

	if err != nil {
		log.Fatal(err)
	}
	fpw := base64.StdEncoding.EncodeToString(HashPW)
	return fpw
}

// ValidateLogin BackLogin return User for the future to expand
func ValidateLogin(id string, password string) (User, int) {
	var user User

	Config.DB.Where("user_id = ?", id).First(&user)

	if user.ID == 0 {
		return user, utils.ErrorUserNotExist
	}

	if ScryptPW(id+password) != user.Password {
		return user, utils.ErrorPasswordWrong
	}
	//if user.Role != 1 {
	//	return user, utils.ErrorUserNotRight
	//}
	return user, utils.SUCCESS
}
