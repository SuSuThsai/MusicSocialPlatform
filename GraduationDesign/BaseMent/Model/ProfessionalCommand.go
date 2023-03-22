package Model

import (
	"GraduationDesign/BaseMent/Config"
	"math"
	"sort"
)

func GetAUserCommandMusic30(userId string) []Music {
	var userHabbity []UserListenTypeCount
	var allUserHabbity [][]UserListenTypeCount
	var users []User
	Config.DB.Find(&users).Model(&User{})
	for i := 0; i < len(users); i++ {
		if users[i].UserId == userId {
			continue
		}
		var data []UserListenTypeCount
		Config.DB.Where("user_id = ?", userId).Order("listen_count DESC").Find(&data).Model(&UserListenTypeCount{})
		allUserHabbity = append(allUserHabbity, data)
	}
	Config.DB.Where("user_id = ?", userId).Order("listen_count DESC").Find(&userHabbity).Model(&UserListenTypeCount{})
	var similarities []float64
	for _, userHabbity1 := range allUserHabbity {
		similarity := calculateSimilarity(userHabbity, userHabbity1)
		similarities = append(similarities, similarity)
	}
	var recommendedMusic []Music
	sort.Slice(similarities, func(i, j int) bool {
		return similarities[i] > similarities[j]
	})
	for i, similarity := range similarities {
		if similarity > 0.5 { // threshold for similarity
			for _, musicType := range allUserHabbity[i] {
				var music []UserListenMusicCount
				Config.DB.Where("user_id = ?", musicType.UserId).Order("listen_count DESC").Limit(100).Find(&music)
				for _, data := range music {
					var musicTopic []MusicTopic
					Config.DB.Where("user_id = ?", data.MusicId).Order("listen_count DESC").Limit(10).Find(&musicTopic)
					for _, topic := range musicTopic {
						if topic.Tip == musicType.Habits {
							a, _ := GetAMusic(data.MusicId)
							recommendedMusic = append(recommendedMusic, a)
						}
						if len(recommendedMusic) > 30 {
							return recommendedMusic
						}
					}
				}
			}
		}
	}
	return recommendedMusic
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
