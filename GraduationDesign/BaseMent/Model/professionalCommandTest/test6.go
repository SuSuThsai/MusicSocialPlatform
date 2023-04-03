package main

import (
	"GraduationDesign/BaseMent/Config"
	"GraduationDesign/BaseMent/Model"
	"fmt"
	"math"
)

func main() {
	Config.InitsConfig()
	Config.InitsPSQL()
	userId := "1906100044"
	var userHabbity []Model.UserListenTypeCount
	var allUserHabbity [][]Model.UserListenTypeCount
	var users []Model.User
	Config.DB.Model(&Model.User{}).Find(&users)
	for i := 0; i < len(users); i++ {
		if users[i].UserId == userId {
			continue
		}
		fmt.Println(users[i].UserId)
		var data []Model.UserListenTypeCount
		Config.DB.Model(&Model.UserListenTypeCount{}).Where("user_id = ?", users[i].UserId).Order("listen_count DESC").Find(&data)
		allUserHabbity = append(allUserHabbity, data)
	}

	Config.DB.Model(&Model.UserListenTypeCount{}).Where("user_id = ?", userId).Order("listen_count DESC").Find(&userHabbity)
	var similarities []float64
	var similarities2 []float64
	for _, userHabbity1 := range allUserHabbity {
		similarity := calculateSimilarity(userHabbity, userHabbity1)
		similarity2 := calculateSimilarity2(userHabbity, userHabbity1)
		similarities = append(similarities, similarity)
		similarities2 = append(similarities2, similarity2)
	}
	fmt.Println("皮尔逊相关系数：", similarities)
	fmt.Println("余弦相似度：", similarities2)
}

// 皮尔逊相关系数来计算用户之间的相似度
func calculateSimilarity(userHabbity []Model.UserListenTypeCount, userHabbity1 []Model.UserListenTypeCount) float64 {
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
func calculateSimilarity2(userHabbity []Model.UserListenTypeCount, userHabbity1 []Model.UserListenTypeCount) float64 {
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
