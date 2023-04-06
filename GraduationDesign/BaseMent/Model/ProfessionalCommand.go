package Model

import (
	"GraduationDesign/BaseMent/Config"
	"math"
	"sort"
)

func GetAUserCommandMusic30(userId string) []Music {
	var userHabbity []UserListenTypeCount
	type X struct {
		allUserHabbity []UserListenTypeCount
		userid         int
	}
	var allUserHabbity []X
	var users []User
	Config.DB.Model(&User{}).Find(&users)
	for i := 0; i < len(users); i++ {
		if users[i].UserId == userId {
			continue
		}
		var data []UserListenTypeCount
		Config.DB.Model(&UserListenTypeCount{}).Where("user_id = ?", users[i].UserId).Order("listen_count DESC").Find(&data)
		allUserHabbity = append(allUserHabbity, X{allUserHabbity: data, userid: i})
	}

	Config.DB.Model(&UserListenTypeCount{}).Where("user_id = ?", userId).Order("listen_count DESC").Find(&userHabbity)
	type Y struct {
		similarities float64
		userid       int
	}
	var similarities []Y
	for _, userHabbity1 := range allUserHabbity {
		similarity := calculateSimilarity2(userHabbity, userHabbity1.allUserHabbity)
		similarities = append(similarities, Y{similarities: similarity, userid: userHabbity1.userid})
	}
	var recommendedMusic []uint
	flag1 := make(map[uint]bool)
	sort.Slice(similarities, func(i, j int) bool {
		return similarities[i].similarities > similarities[j].similarities
	})
	for _, similarity := range similarities {
		if similarity.similarities > 0.3 { // threshold for similarity
			//fmt.Println(similarity.userid, "6666")
			var music []UserListenMusicCount
			Config.DB.Where("user_id = ?", users[similarity.userid].UserId).Order("listen_count DESC").Limit(100).Find(&music)
			//fmt.Println(music, "7777")
			for _, data := range music {
				var musicTopic []MusicTopic
				Config.DB.Where("id = ?", data.MusicId).Limit(10).Find(&musicTopic)
				for _, musicType := range userHabbity {
					for _, topic := range musicTopic {
						if topic.Tip == musicType.Habits {
							//a, _ := GetAMusic(data.MusicId)
							if !flag1[data.MusicId] {
								recommendedMusic = append(recommendedMusic, data.MusicId)
								flag1[data.MusicId] = true
							}
							break
						}
						//if len(recommendedMusic) > 30 {
						//	return recommendedMusic
						//}
					}
				}
			}
		}
	}
	//fmt.Println(len(recommendedMusic), recommendedMusic, "222222")
	//TODO 改变推荐逻辑 改为优先级推荐
	var userMusicListenCount []UserListenMusicCount
	Config.DB.Model(&UserListenMusicCount{}).Where("user_id = ?", userId).Find(&userMusicListenCount)
	flag := map[int]int{}
	for i := 0; i < len(userMusicListenCount); i++ {
		flag[int(userMusicListenCount[i].MusicId)] = userMusicListenCount[i].ListenCount
	}
	//三类 <5 5~10 10~20 >20
	type musicCount struct {
		music uint
		num   int
	}
	countPic := make([][]musicCount, 4)
	for i := 0; i < len(recommendedMusic); i++ {
		if flag[int(recommendedMusic[i])] < 5 {
			countPic[0] = append(countPic[0], musicCount{music: recommendedMusic[i], num: flag[int(recommendedMusic[i])]})
		} else if 5 <= flag[int(recommendedMusic[i])] && flag[int(recommendedMusic[i])] < 10 {
			countPic[1] = append(countPic[1], musicCount{music: recommendedMusic[i], num: flag[int(recommendedMusic[i])]})
		} else if 10 <= flag[int(recommendedMusic[i])] && flag[int(recommendedMusic[i])] < 20 {
			countPic[2] = append(countPic[2], musicCount{music: recommendedMusic[i], num: flag[int(recommendedMusic[i])]})
		} else {
			countPic[3] = append(countPic[3], musicCount{music: recommendedMusic[i], num: flag[int(recommendedMusic[i])]})
		}
	}
	for i := 0; i < len(countPic); i++ {
		sort.Slice(countPic[i], func(j, k int) bool {
			return countPic[i][j].num < countPic[i][k].num
		})
	}
	var result []Music
	flag2 := make(map[uint]bool)
	for i := 0; i < len(countPic); i++ {
		for j := 0; j < len(countPic[i]); j++ {
			if !flag2[countPic[i][j].music] {
				a, _ := GetAMusic(countPic[i][j].music)
				result = append(result, a)
				flag2[countPic[i][j].music] = true
			}
			if len(result) >= 30 {
				return result
			}
		}
	}
	return result
}

// 皮尔逊相关系数来计算用户之间的相似度
func calculateSimilarity(userHabbity []UserListenTypeCount, userHabbity1 []UserListenTypeCount) float64 {
	var sum1, sum2, sum3 float64
	var avg1, avg2 float64
	for _, targetMusicType := range userHabbity {
		avg1 += float64(targetMusicType.ListenCount)
	}
	avg1 /= float64(len(userHabbity))
	for _, musicType := range userHabbity1 {
		avg2 += float64(musicType.ListenCount)
	}
	avg2 /= float64(len(userHabbity1))
	for _, targetMusicType := range userHabbity {
		for _, musicType := range userHabbity1 {
			if targetMusicType.Habits == musicType.Habits {
				sum1 += (float64(targetMusicType.ListenCount) - avg1) * (float64(musicType.ListenCount) - avg2)
				sum2 += (float64(targetMusicType.ListenCount) - avg1) * (float64(targetMusicType.ListenCount) - avg1)
				sum3 += (float64(musicType.ListenCount) - avg2) * (float64(musicType.ListenCount) - avg2)
				break
			}
		}
	}
	if sum2 == 0 || sum3 == 0 {
		return 0
	}
	return sum1 / (math.Sqrt(sum2) * math.Sqrt(sum3))
}

// 余弦相似度来计算两个用户之间的相似度
func calculateSimilarity2(userHabbity []UserListenTypeCount, userHabbity1 []UserListenTypeCount) float64 {
	var sum, sum1, sum2 float64
	for _, targetMusicType := range userHabbity {
		for _, musicType := range userHabbity1 {
			if targetMusicType.Habits == musicType.Habits {
				sum += float64(targetMusicType.ListenCount * musicType.ListenCount)
				break
			}
		}
		sum1 += float64(targetMusicType.ListenCount * targetMusicType.ListenCount)
	}
	for _, musicType := range userHabbity1 {
		sum2 += float64(musicType.ListenCount * musicType.ListenCount)
	}
	if sum1 == 0 || sum2 == 0 {
		return 0
	}
	return sum / (math.Sqrt(sum1) * math.Sqrt(sum2))
}
