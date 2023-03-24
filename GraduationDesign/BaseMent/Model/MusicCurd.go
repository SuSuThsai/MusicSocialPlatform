package Model

import (
	"GraduationDesign/BaseMent/Config"
	"GraduationDesign/BaseMent/utils"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"log"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

// GetMusicRankWeek Get Music RankWeek
func GetMusicRankWeek(y1 int, m string, week1 int) ([]Music, int, int64) {
	y := strconv.Itoa(y1)
	week := strconv.Itoa(week1)
	var musics []MusicRankWeek
	err = Config.DB.Where("year = ? and month = ? and week = ?", y, m, week).Find(&musics).Error
	var result []Music
	for i := 0; i < len(musics); i++ {
		a, _ := GetAMusic(musics[i].Id)
		result = append(result, a)
	}
	if err != nil {
		return nil, utils.ERROR, 0
	}
	return result, utils.SUCCESS, int64(len(musics))
}

// GetMusicRankMonth Get Music RankMonth
func GetMusicRankMonth(y1 int, m string) ([]Music, int, int64) {
	y := strconv.Itoa(y1)
	var musics []MusicRankMonth
	err = Config.DB.Where("year = ? and month = ?", y, m).Limit(100).Order("count DESC").Find(&musics).Error
	if err != nil {
		return nil, utils.ERROR, 0
	}
	var result []Music
	for i := 0; i < len(musics); i++ {
		a, _ := GetAMusic(musics[i].Id)
		result = append(result, a)
	}
	return result, utils.SUCCESS, int64(len(musics))
}

// GetMusicRankYear Get Music RankYear
func GetMusicRankYear(y1 int) ([]Music, int, int64) {
	y := strconv.Itoa(y1)
	var musics []MusicRankYear
	err = Config.DB.Where("year = ?", y).Limit(100).Order("count DESC").Find(&musics).Error
	if err != nil {
		return nil, utils.ERROR, 0
	}
	var result []Music
	for i := 0; i < len(musics); i++ {
		a, _ := GetAMusic(musics[i].Id)
		result = append(result, a)
	}
	return result, utils.SUCCESS, int64(len(musics))
}

// GetAMusicList Get A MusicList
func GetAMusicList(data []redis.Z) ([]Music, int, int64) {
	var music []Music
	id := 0
	for i, datum := range data {
		var data1 Music
		temp := datum.Member.(string)
		temp2 := strings.Split(temp, " ")
		id, _ = strconv.Atoi(temp2[len(temp2)-1])
		err = Config.DB.Where("id = ?",
			id).First(&data1).Error
		if err != nil {
			log.Println("更新日派第", i, "名歌曲失败", "musicId:", id)
			continue
		}
		music = append(music, data1)
	}
	return music, utils.SUCCESS, int64(len(music))
}

// GetMusicListRankWeek Get MusicList RankWeek
func GetMusicListRankWeek(y1 int, m string, week1 int) ([]MusicList, int, int64) {
	y := strconv.Itoa(y1)
	week := strconv.Itoa(week1)
	var musics []MusicListRankWeek
	err = Config.DB.Where("year = ? and month = ? and week = ?", y, m, week).Limit(100).Order("count DESC").Find(&musics).Error
	if err != nil {
		return nil, utils.ERROR, 0
	}
	var result []MusicList
	for i := 0; i < len(musics); i++ {
		a, _ := FindAMusicList(musics[i].Id)
		result = append(result, a)
	}
	return result, utils.SUCCESS, int64(len(musics))
}

// GetMusicListRankMonth Get MusicList RankMonth
func GetMusicListRankMonth(y1 int, m string) ([]MusicList, int, int64) {
	y := strconv.Itoa(y1)
	var musics []MusicRankMonth
	err = Config.DB.Where("year = ? and month = ?", y, m).Limit(100).Order("count DESC").Find(&musics).Error
	if err != nil {
		return nil, utils.ERROR, 0
	}
	var result []MusicList
	for i := 0; i < len(musics); i++ {
		a, _ := FindAMusicList(musics[i].Id)
		result = append(result, a)
	}
	return result, utils.SUCCESS, int64(len(musics))
}

// GetMusicListRankYear Get MusicList RankYear
func GetMusicListRankYear(y1 int) ([]MusicList, int, int64) {
	y := strconv.Itoa(y1)
	var musics []MusicListRankYear
	err = Config.DB.Where("year = ?", y).Limit(100).Order("count DESC").Find(&musics).Error
	if err != nil {
		return nil, utils.ERROR, 0
	}
	var result []MusicList
	for i := 0; i < len(musics); i++ {
		a, _ := FindAMusicList(musics[i].Id)
		result = append(result, a)
	}
	return result, utils.SUCCESS, int64(len(musics))
}

// GetAMusicListList Get A MusicListList
func GetAMusicListList(data []redis.Z) ([]MusicList, int, int64) {
	var musicList []MusicList
	id := 0
	for i, datum := range data {
		var data1 MusicList
		temp := datum.Member.(string)
		temp2 := strings.Split(temp, " ")
		id, _ = strconv.Atoi(temp2[len(temp2)-1])
		err := Config.DB.Where("id = ?",
			id).First(&data1).Error
		if err != nil {
			log.Println("更新日派第", i, "名歌单失败", "musicId:", id)
			continue
		}
		musicList = append(musicList, data1)
	}
	return musicList, utils.SUCCESS, int64(len(musicList))
}

func GetAMusic(musicId uint) (Music, int) {
	var music Music
	Config.DB.Where("id = ?", musicId).First(&music)
	if music.ID == 0 {
		return music, utils.ERROR
	}
	return music, utils.SUCCESS
}

func FindAMusicList(listId uint) (MusicList, int) {
	var musicList MusicList
	Config.DB.Where("id = ?", listId).First(&musicList)
	if musicList.ID == 0 {
		return musicList, utils.ERROR
	}
	return musicList, utils.SUCCESS
}

// SearchUserMusicsLike Search User The Musics Like
func SearchUserMusicsLike(userId uint, pageSize int, pageNum int) ([]MusicLike, []Music, int, int64) {
	var music []MusicLike
	var total int64
	err = Config.DB.Order("created_at DESC").Where("user_id = ?",
		userId).Find(&music).Count(&total).Error
	//err = Config.DB.Limit(pageSize).Offset((pageNum-1)*pageSize).Order("Created_At DESC").Where("user_id = ? and like = true",
	//	userId).Find(&music).Count(&total).Error
	var result []MusicLike
	var result2 []Music
	for i := 0; i < len(music); i++ {
		if music[i].Like == true {
			result = append(result, music[i])
			a, _ := GetAMusic(music[i].MusicId)
			result2 = append(result2, a)
		} else {
			total -= 1
		}
	}
	if err != nil {
		return result, result2, utils.ERROR, 0
	}
	return result, result2, utils.SUCCESS, total
}

// SearchAllMusics Search Musics
func SearchAllMusics(pageSize int, pageNum int) ([]Music, int, int64) {
	var music []Music
	var total int64
	err = Config.DB.Preload("MusicListen", "listen DESC").Find(&music).Count(&total).Error
	//err = Config.DB.Preload("MusicListen", "listen DESC").Limit(pageSize).Offset((pageNum - 1) * pageSize).Find(&music).Count(&total).Error
	if err != nil {
		return music, utils.ERROR, 0
	}
	return music, utils.SUCCESS, total
}

// SearchAllMusicList Search All MusicLists
func SearchAllMusicList(pageSize int, pageNum int) ([]MusicList, int, int64) {
	var musicList []MusicList
	var total int64
	err = Config.DB.Order("like_count DESC").Find(&musicList).Count(&total).Error
	//err = Config.DB.Limit(pageSize).Offset((pageNum - 1) * pageSize).Order("listen DESC").Find(&musicList).Count(&total).Error
	if err != nil {
		return musicList, utils.ERROR, 0
	}
	return musicList, utils.SUCCESS, total
}

// SearchUserMusicsListLike Search User The MusicList Like
func SearchUserMusicsListLike(userId uint, pageSize int, pageNum int) ([]MusicListLike, []MusicList, int, int64) {
	var musicList []MusicListLike
	var total int64
	err = Config.DB.Order("created_at DESC").Where("user_id = ?",
		userId).Find(&musicList).Count(&total).Error
	//err = Config.DB.Limit(pageSize).Offset((pageNum-1)*pageSize).Order("Created_At DESC").Where("user_id = ? and like = true",
	//	userId).Find(&music).Count(&total).Error
	var result []MusicListLike
	var result2 []MusicList
	for i := 0; i < len(musicList); i++ {
		if musicList[i].Like == true {
			result = append(result, musicList[i])
			a, _ := FindAMusicList(musicList[i].ListId)
			result2 = append(result2, a)
		}
	}
	if err != nil {
		return result, result2, utils.ERROR, 0
	}
	return result, result2, utils.SUCCESS, total
}

// SearchMusics Search The Musics
func SearchMusics(title string, pageSize int, pageNum int) ([]Music, int, int64) {
	var music []Music
	var total int64
	err = Config.DB.Order("created_at DESC").Where("name LIKE ? or singer LIKE ?",
		"%"+title+"%", "%"+title+"%").Find(&music).Count(&total).Error
	//err = Config.DB.Limit(pageSize).Offset((pageNum-1)*pageSize).Order("Created_At DESC").Where("name LIKE ?",
	//	"%"+title+"%").Find(&music).Count(&total).Error
	if err != nil {
		return music, utils.ERROR, 0
	}
	return music, utils.SUCCESS, total
}

func GetCommandMusicDays(userId string, title int) ([]Music, int) {
	m, d, n := time.Now().AddDate(0, 0, -title).Truncate(24*time.Hour).Month().String(), time.Now().AddDate(0, 0, -title).Truncate(24*time.Hour).Weekday().String(), strconv.Itoa(time.Now().AddDate(0, 0, -title).Truncate(24*time.Hour).Day())
	y1, w1 := time.Now().AddDate(0, 0, -title).Truncate(24 * time.Hour).ISOWeek()
	w := strconv.Itoa(w1)
	y := strconv.Itoa(y1)
	var data2 []CommandMusicCount
	Config.DB.Where("user_id =  ? and year = ? and month = ? and week = ? and day = ? and number = ?",
		userId, y, m, w, d, n).Find(&data2)
	var result []Music
	for i := 0; i < len(data2); i++ {
		a, _ := GetAMusic(data2[i].MusicId)
		result = append(result, a)
	}
	return result, utils.SUCCESS
}

func GetCommandMusic(userId string) ([]Music, int) {
	m, d, _ := time.Now().Month().String(), time.Now().Weekday().String(), strconv.Itoa(time.Now().Day())
	y1, w1 := time.Now().ISOWeek()
	w := strconv.Itoa(w1)
	y := strconv.Itoa(y1)
	var data2 []CommandMusicCount
	Config.DB.Where("user_id = ? and year = ? and month = ? and week = ? and day = ?",
		userId, y, m, w, d).Find(&data2)
	var result []Music
	if len(data2) != 0 {
		for i := 0; i < len(data2); i++ {
			a, _ := GetAMusic(data2[i].MusicId)
			result = append(result, a)
		}
		return result, utils.SUCCESS
	}
	return result, utils.ERROR
}

func CountCommandMusic(musics []Music, userId string) int {
	m, d, n := time.Now().Month().String(), time.Now().Weekday().String(), strconv.Itoa(time.Now().Day())
	y1, w1 := time.Now().ISOWeek()
	w := strconv.Itoa(w1)
	y := strconv.Itoa(y1)
	var data2 CommandMusicCount
	Config.DB.Where("year = ? and month = ? and week = ? and day = ?",
		y, m, w, d).Find(&data2)
	for i := 0; i < len(musics); i++ {
		var data CommandMusicCount
		data.MusicId = musics[i].Id
		data.UserId = userId
		data.Year = y
		data.Month = m
		data.Week = w
		data.Day = d
		data.Number = n
		Config.DB.Create(&data).Model(&CommandMusicCount{})
	}
	return utils.SUCCESS
}

// SearchMusicsProfessional Search The Musics Professional
func SearchMusicsProfessional(titles []string) ([]Music, int) {
	var music []MusicTopic
	for _, title := range titles {
		var data []MusicTopic
		Config.DB.Order("Created_At DESC").Where("tip LIKE ?",
			title).Find(&data)
		music = append(music, data...)
	}
	music2 := RandomSongs(music)
	var reslut []Music
	for i := 0; i < len(music2); i++ {
		a, _ := GetAMusic(music2[i].Id)
		reslut = append(reslut, a)
	}
	return reslut, utils.SUCCESS
}

func RandomSongs(music []MusicTopic) []MusicTopic {
	if len(music) <= 30 {
		return music
	}
	nums := map[int]bool{}
	nums2 := map[uint]bool{}
	var result []MusicTopic
	lenght := len(music)
	rand.Seed(time.Now().UnixNano())
	for len(result) < 30 {
		rand := rand.Intn(lenght)
		if nums[rand] == false && nums2[music[rand].Id] == false {
			nums[rand] = true
			nums2[music[rand].Id] = true
			result = append(result, music[rand])
		} else if nums2[music[rand].Id] == true {
			nums[rand] = true
		}
		if len(nums) == len(music) || len(nums2) == len(music) {
			break
		}
	}
	return result
}

// SearchMusicsSinger Search The MusicsSinger
func SearchMusicsSinger(title string, pageSize int, pageNum int) ([]Music, int, int64) {
	var music []Music
	var total int64
	err = Config.DB.Limit(pageSize).Offset((pageNum-1)*pageSize).Order("Created_At DESC").Where("singer LIKE ?",
		"%"+title+"%").Find(&music).Count(&total).Error
	if err != nil {
		return music, utils.ERROR, 0
	}
	return music, utils.SUCCESS, total
}

// MusicsListen Musics has Listen
func MusicsListen(id uint) int {
	var data MusicListen
	err = Config.DB.Where("id = ?", id).First(&data).Error
	if data.Id == 0 {
		data.Id = id
		Config.DB.Create(&data)
	} else if err != nil {
		return utils.ERROR
	}
	err = Config.DB.Model(&data).Where("id = ?", id).Update("listen", gorm.Expr("listen+ ?", 1)).Error
	if err != nil {
		return utils.ERROR
	}
	return utils.SUCCESS
}

// MusicsLikeCount MusicsLike Count
func MusicsLikeCount(id uint) (int, int64) {
	var data MusicLike
	var total int64
	err = Config.DB.Where("id = ? and like = true", id).Find(&data).Count(&total).Error
	if err != nil && total != 0 {
		log.Println("MusicLikeCountERR", err)
		return utils.ERROR, 0
	}
	return utils.SUCCESS, total
}

func MusicRankCount(id uint, y string, m string, w string, d string, v string, num int) int {
	var musicY MusicRankYear
	var musicM MusicRankMonth
	var musicW MusicRankWeek
	var musicD MusicRankDay
	err = Config.DB.Where("id = ? and year = ?", id, y).Model(&musicY).Error
	if err == gorm.ErrRecordNotFound || musicY.Id == 0 {
		musicY = MusicRankYear{
			Id:    id,
			Year:  y,
			Count: num,
		}
		err = Config.DB.Create(&musicY).Error
		if err != nil {
			log.Println("创建年记录数据失败:", err, "y", y, "musicId:", id, "count:", num)
		}
	}
	err = Config.DB.Model(&musicY).Where("id = ? and year = ?", id, y).Update("count", gorm.Expr("count+ ?", num)).Error
	if err != nil {
		log.Println("增加年记录数据失败:", err, "y", y, "musicId:", id, "count:", num)
	}
	err = Config.DB.Where("id = ? and year = ? and month = ?", id, y, m).Model(&musicM).Error
	if err == gorm.ErrRecordNotFound || musicM.Id == 0 {
		musicM = MusicRankMonth{
			Id:    id,
			Year:  y,
			Month: m,
			Count: num,
		}
		err = Config.DB.Create(&musicM).Error
		if err != nil {
			log.Println("创建年月记录数据失败:", err, "y", y, "m", m, "musicId:", id, "count:", num)
		}
	}
	err = Config.DB.Model(&musicM).Where("id = ? and year = ? and month = ?", id, y, m).Update("count", gorm.Expr("count+ ?", num)).Error
	if err != nil {
		log.Println("增加年月记录数据失败:", err, "y", y, "m", m, "musicId:", id, "count:", num)
	}
	err = Config.DB.Where("id = ? and year = ? and month = ? and week = ?", id, y, m, w).Model(&musicW).Error
	if err == gorm.ErrRecordNotFound || musicW.Id == 0 {
		musicW = MusicRankWeek{
			Id:    id,
			Year:  y,
			Month: m,
			Week:  w,
			Count: num,
		}
		err = Config.DB.Create(&musicW).Error
		if err != nil {
			log.Println("创建年月周记录数据失败:", err, "y", y, "m", m, "w", w, "musicId:", id, "count:", num)
		}
	}
	err = Config.DB.Model(&musicW).Where("id = ? and year = ? and month = ? and week = ?", id, y, m, w).Update("count", gorm.Expr("count+ ?", num)).Error
	if err != nil {
		log.Println("增加年月周记录数据失败:", err, "y", y, "m", m, "w", w, "musicId:", id, "count:", num)
	}
	err = Config.DB.Where("id = ? and year = ? and month = ? and week = ? and day = ?", id, y, m, w, d).Model(&musicD).Error
	if err == gorm.ErrRecordNotFound || musicD.Id == 0 {
		musicD = MusicRankDay{
			Id:     id,
			Year:   y,
			Month:  m,
			Week:   w,
			Day:    d,
			Number: v,
			Count:  num,
		}
		err = Config.DB.Create(&musicD).Error
		if err != nil {
			log.Println("创建年月周日记录数据失败:", err, "y", y, "m", m, "w", w, "d", d, "musicId:", id, "count:", num)
		}
	}
	err = Config.DB.Model(&musicD).Where("id = ? and year = ? and month = ? and week = ? and day = ?", id, y, m, w, d).Update("count", gorm.Expr("count+ ?", num)).Error
	if err != nil {
		log.Println("增加年月周日记录数据失败:", err, "y", y, "m", m, "w", w, "d", d, "musicId:", id, "count:", num)
	}
	return utils.SUCCESS
}

func MusicListRankCount(id uint, y string, m string, w string, d string, v string, num int) int {
	var musicListY MusicListRankYear
	var musicListM MusicListRankMonth
	var musicListW MusicListRankWeek
	var musicListD MusicListRankDay
	err = Config.DB.Where("id = ? and year = ?", id, y).Model(&musicListY).Error
	if err == gorm.ErrRecordNotFound || musicListY.Id == 0 {
		musicListY = MusicListRankYear{
			Id:    id,
			Year:  y,
			Count: num,
		}
		err = Config.DB.Create(&musicListY).Error
		if err != nil {
			log.Println("创建年记录数据失败:", err, "y", y, "musicListId:", id, "count:", num)
		}
	}
	err = Config.DB.Model(&musicListY).Where("id = ? and year = ?", id, y).Update("count", gorm.Expr("count+ ?", num)).Error
	if err != nil {
		log.Println("增加年记录数据失败:", err, "y", y, "musicListId:", id, "count:", num)
	}
	err = Config.DB.Where("id = ? and year = ? and month = ?", id, y, m).Model(&musicListM).Error
	if err == gorm.ErrRecordNotFound || musicListM.Id == 0 {
		musicListM = MusicListRankMonth{
			Id:    id,
			Year:  y,
			Month: m,
			Count: num,
		}
		err = Config.DB.Create(&musicListM).Error
		if err != nil {
			log.Println("创建年月记录数据失败:", err, "y", y, "m", m, "musicListId:", id, "count:", num)
		}
	}
	err = Config.DB.Model(&musicListM).Where("id = ? and year = ? and month = ?", id, y, m).Update("count", gorm.Expr("count+ ?", num)).Error
	if err != nil {
		log.Println("增加年月记录数据失败:", err, "y", y, "m", m, "musicListId:", id, "count:", num)
	}
	err = Config.DB.Where("id = ? and year = ? and month = ? and week = ?", id, y, m, w).Model(&musicListW).Error
	if err == gorm.ErrRecordNotFound || musicListW.Id == 0 {
		musicListW = MusicListRankWeek{
			Id:    id,
			Year:  y,
			Month: m,
			Week:  w,
			Count: num,
		}
		err = Config.DB.Create(&musicListW).Error
		if err != nil {
			log.Println("创建年月周记录数据失败:", err, "y", y, "m", m, "w", w, "musicListId:", id, "count:", num)
		}
	}
	err = Config.DB.Model(&musicListW).Where("id = ? and year = ? and month = ? and week = ?", id, y, m, w).Update("count", gorm.Expr("count+ ?", num)).Error
	if err != nil {
		log.Println("增加年月周记录数据失败:", err, "y", y, "m", m, "w", w, "musicListId:", id, "count:", num)
	}
	err = Config.DB.Where("id = ? and year = ? and month = ? and week = ? and day = ?", id, y, m, w, d).Model(&musicListD).Error
	if err == gorm.ErrRecordNotFound || musicListD.Id == 0 {
		musicListD = MusicListRankDay{
			Id:     id,
			Year:   y,
			Month:  m,
			Week:   w,
			Day:    d,
			Number: v,
			Count:  num,
		}
		err = Config.DB.Create(&musicListD).Error
		if err != nil {
			log.Println("创建年月周日记录数据失败:", err, "y", y, "m", m, "w", w, "d", d, "musicListId:", id, "count:", num)
		}
	}
	err = Config.DB.Model(&musicListD).Where("id = ? and year = ? and month = ? and week = ? and day = ?", id, y, m, w, d).Update("count", gorm.Expr("count+ ?", num)).Error
	if err != nil {
		log.Println("增加年月周日记录数据失败:", err, "y", y, "m", m, "w", w, "d", d, "musicListId:", id, "count:", num)
	}
	return utils.SUCCESS
}

// AddMusic UpLoad Music
func AddMusic(data *Music) int {
	err = Config.DB.Create(&data).Model(&Music{}).Error
	if err != nil {
		return utils.ERROR
	}
	return utils.SUCCESS
}

// UpLoadMusic UpLoad Music
func UpLoadMusic() {
}

func CreatMusicHabit(musicId uint, musicName string, tip string) int {
	var musicTopic MusicTopic
	musicTopic.Id = musicId
	musicTopic.Name = musicName
	musicTopic.Tip = tip
	err = Config.DB.Create(&musicTopic).Model(&MusicTopic{}).Error
	if err != nil {
		return utils.ERROR
	}
	return utils.SUCCESS
}

func GetMusicHabit(musicId uint, musicName string) []MusicTopic {
	var musicTopic []MusicTopic
	Config.DB.Where("id = ? or name = ?", musicId, musicName).Find(&musicTopic)
	return musicTopic
}

func GetMusicListHabit(listId uint) []Tips {
	var musicTopic []Tips
	Config.DB.Where("list_id = ?", listId).Find(&musicTopic)
	if len(musicTopic) > 5 {
		return musicTopic[:5]
	}
	return musicTopic
}
