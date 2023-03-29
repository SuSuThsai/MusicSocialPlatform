package main

import (
	"GraduationDesign/BaseMent/Config"
	"GraduationDesign/BaseMent/Model"
	"fmt"
	"os"
	"strconv"
)

func main() {
	//Config.InitsSQL()
	//PersistentSql.InitCache()
	//i := 0
	//for i < 2 {
	//	err := PersistentSql.TestQueue.SendScheduleMsg(strconv.Itoa(i)+"nana", time.Now().Add(0), PersistentSql.WithRetryCount(3))
	//	err = PersistentSql.TestQueue2.SendScheduleMsg(strconv.Itoa(i)+"quansi", time.Now().Add(0), PersistentSql.WithRetryCount(3))
	//	if err != nil {
	//		panic(err)
	//	}
	//	i++
	//}
	//fmt.Println("finished")
	//time.Sleep(15 * time.Second)
	Config.InitsConfig()
	Config.InitsPSQL()
	//dir, err := os.Open("H:/MP3/cn_song")
	//if err != nil {
	//	// 目录打开失败，执行错误处理逻辑
	//}
	//defer dir.Close()
	//// 读取目录中的所有文件
	//files, err := dir.Readdir(0)
	//if err != nil {
	//	// 读取目录内容失败，执行错误处理逻辑
	//}
	file1, _ := os.OpenFile("BaseMent/Model/Cache/PersistentSql/test1/data2.txt", os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	defer file1.Close()
	var music []Model.Music
	Config.DB.Where("source = ?", "").Find(&music)
	for _, m := range music {
		id := strconv.Itoa(int(m.Id))
		_, _ = fmt.Fprintln(file1, id+" - "+m.Name+" - ")
	}
	//for _, file := range files {
	//	if !file.IsDir() { // 判断是否是目录
	//		// 在这里对每个文件进行操作，比如打印文件名
	//		data := strings.Split(file.Name(), " - ")
	//		var music Model.Music
	//		if len(data) == 1 {
	//			musicname := strings.Split(data[0], ".mp3")
	//			music.Name = musicname[0]
	//			id := Model.AddMusic2(&music)
	//			id1 := strconv.Itoa(int(id))
	//			_, err = fmt.Fprintln(file1, id1+" - "+musicname[0]+" - ")
	//		} else {
	//			singer := data[0]
	//			musicname := strings.Split(data[1], ".mp3")
	//			music.Name = musicname[0]
	//			music.Singer = singer
	//			id := Model.AddMusic2(&music)
	//			id1 := strconv.Itoa(int(id))
	//			_, err = fmt.Fprintln(file1, id1+" - "+musicname[0]+" - ")
	//		}
	//	}
	//}
	//code := Model.AddMusic(&music)
}
