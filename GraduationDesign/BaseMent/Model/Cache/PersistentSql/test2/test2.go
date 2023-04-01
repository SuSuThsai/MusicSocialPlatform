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
	//Config.DB.Where("id = ?", 329).Delete(&Model.Music{})
	//Config.DB.Where("id = ?", 312).Delete(&Model.Music{})
	file1, _ := os.Open("BaseMent/Model/Cache/PersistentSql/test1/data3.txt")
	defer file1.Close()
	scanner := bufio.NewScanner(file1)
	for scanner.Scan() {
		data := strings.Split(scanner.Text(), " - ")
		id, _ := strconv.Atoi(data[0])
		data2 := strings.Trim(data[3], " ")
		var data3 = make(map[string]interface{})
		data3["words_source"] = data2
		Config.DB.Model(&Model.Music{}).Where("id = ?", id).Updates(&data3)
	}
	//http.Header.Set("Access-Control-Allow-Origin", "*")
	//_, err := http.Get("https://43.139.123.205:18760/down/IhEqG5Nz5OY2.lrc")
	//fmt.Println(err)
}
