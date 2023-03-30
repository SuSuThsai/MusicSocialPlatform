package main

import (
	"GraduationDesign/BaseMent/Config"
	"GraduationDesign/BaseMent/Model"
	"bufio"
	"os"
	"strconv"
	"strings"
)

func main() {
	Config.InitsConfig()
	Config.InitsPSQL()
	//var result []x
	//TODO 打入标签
	file1, _ := os.Open("C:/Users/Hasee/Desktop/毕设开发进程/MusicSocialPlatform/GraduationDesign/BaseMent/Model/Cache/PersistentSql/test1/data2.txt")
	defer file1.Close()
	scanner := bufio.NewScanner(file1)
	for scanner.Scan() {
		data := strings.Split(scanner.Text(), " -@")
		data2 := strings.Split(data[4], "+")
		id, _ := strconv.Atoi(data[0])
		m, _ := Model.GetAMusic(uint(id))
		for i := 0; i < len(data2); i++ {
			if data2[i] != "" {
				Model.CreatMusicHabit(m.Id, m.Name, data2[i])
			}
		}
	}
	//TODO 获取没有标签的歌曲
	//file1, _ := os.OpenFile("BaseMent/Model/Cache/PersistentSql/test1/data2.txt", os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	//defer file1.Close()
	//var music []Model.Music
	//Config.DB.Find(&music)
	//type x struct {
	//	name   string
	//	id     uint
	//	lenght int
	//}
	//total := 0
	//for j := 0; j < len(music); j++ {
	//	var data []Model.MusicTopic
	//	data = Model.GetMusicHabit(music[j].Id, "")
	//	id := strconv.Itoa(int(music[j].Id))
	//	y := strconv.Itoa(len(data))
	//	if len(data) == 0 {
	//		_, _ = fmt.Fprintln(file1, id+" -@"+music[j].Name+" -@"+" -@"+y+" -@")
	//		total++
	//	}
	//}
	//fmt.Println(total)
	//file, err := excelize.OpenFile("C:/Users/Hasee/Desktop/毕设开发进程/MusicSocialPlatform/GraduationDesign/BaseMent/Model/Cache/PersistentSql/test3/final_emo_v2.xlsx")
	//if err != nil {
	//	fmt.Println("Error:", err)
	//	return
	//}
	//defer file.Close()
	//total := 0
	//total1 := 0
	//total2 := 0
	//rows, _ := file.GetRows("Sheet1")
	//for i := 0; i < len(rows); i++ {
	//	if i < 1 {
	//		continue
	//	}
	//	//TODO 歌词的标签
	//	////a := rows[i][0] // 第几行第几列，行号，列号都从0开始
	//	//a := rows[i][1]
	//	//a1 := strings.Split(a, " - ")
	//	//a2 := ""
	//	//a3 := ""
	//	//if len(a1) == 1 {
	//	//	a2 = strings.Trim(a1[0], ".lrc")
	//	//} else {
	//	//	a2 = strings.Trim(a1[1], ".lrc")
	//	//}
	//	//a3 = strings.Trim(a2, " ")
	//	//b := rows[i][2]
	//	//c := rows[i][3]
	//	//d := rows[i][4]
	//	//var music []Model.Music
	//	//if len(a1) == 1 {
	//	//	Config.DB.Where("name = ?", a3).Find(&music)
	//	//} else {
	//	//	data3 := strings.Trim(a1[0], " ")
	//	//	Config.DB.Where("name = ? and singer = ?", a3, data3).Find(&music)
	//	//}
	//	//for j := 0; j < len(music); j++ {
	//	//	if music[j].Id > 0 {
	//	//		fmt.Println(music[j].Id, music[j].Name, music[j].Singer)
	//	//		total1++
	//	//		Model.CreatMusicHabit(music[j].Id, music[j].Name, b)
	//	//		Model.CreatMusicHabit(music[j].Id, music[j].Name, c)
	//	//		Model.CreatMusicHabit(music[j].Id, music[j].Name, d)
	//	//	}
	//	//}
	//	//total++
	//
	//	//	//TODO 旋律的标签
	//	//	//a := rows[i][0] // 第几行第几列，行号，列号都从0开始
	//	//	emotion := rows[i][7]
	//	//	x := emotion
	//	//	x1 := strings.Trim(x, "25%")
	//	//	x2 := strings.Trim(x1, "50%")
	//	//	x3 := strings.Trim(x2, "75%")
	//	//
	//	//	name := rows[i][8]
	//	//	singer := rows[i][9]
	//	//	name1 := strings.Trim(name, " ")
	//	//	singer1 := strings.Trim(singer, " ")
	//	//	var music []Model.Music
	//	//	Config.DB.Where("name = ? and singer = ?", name1, singer1).Find(&music)
	//	//	for j := 0; j < len(music); j++ {
	//	//		if music[j].Id > 0 {
	//	//			fmt.Println(music[j].Id, music[j].Name, music[j].Singer, x3)
	//	//			total1++
	//	//			Model.CreatMusicHabit(music[j].Id, music[j].Name, x)
	//	//			//Model.CreatMusicHabit(music[j].Id, music[j].Name, c)
	//	//			//Model.CreatMusicHabit(music[j].Id, music[j].Name, d)
	//	//		}
	//	//		total2++
	//	//	}
	//	//	total++
	//}
	//fmt.Println(total, total1, total2)
}
