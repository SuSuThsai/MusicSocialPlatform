package main

import (
	"GraduationDesign/BaseMent/Config"
	"GraduationDesign/BaseMent/Model/Cache/PersistentSql"
	"fmt"
	"time"
)

func main() {
	Config.InitsSQL()
	PersistentSql.InitCache()
	//i := 0
	//for i < 2 {
	//	err := PersistentSql.TestQueue.SendScheduleMsg(strconv.Itoa(i)+"nana", time.Now().Add(0), PersistentSql.WithRetryCount(3))
	//	err = PersistentSql.TestQueue2.SendScheduleMsg(strconv.Itoa(i)+"quansi", time.Now().Add(0), PersistentSql.WithRetryCount(3))
	//	if err != nil {
	//		panic(err)
	//	}
	//	i++
	//}
	fmt.Println("finished")
	time.Sleep(15 * time.Second)
}
