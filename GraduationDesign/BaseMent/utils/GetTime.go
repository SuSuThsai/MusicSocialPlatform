package utils

import (
	"fmt"
	"github.com/robfig/cron/v3"
	"time"
)

func GetMonthWeek() int {
	y := time.Now().Format("2006-01-02")
	date, _ := time.ParseInLocation("2006-01-02", y, time.Local)
	// week of the year
	fmt.Println(date, time.Now().GoString(), y)
	_, week := date.ISOWeek()

	weekday := date.Weekday()
	if weekday == 0 {
		weekday = 7
	}
	// 获取日期所在周的周四日期
	thursday := date.AddDate(0, 0, int(4-weekday))
	//// 所属月
	//month := int(thursday.Month())
	//// 所属年
	//year := thursday.Year()
	// 获取日期所在周的周四所在的月是第几周
	_, week1 := time.Date(thursday.Year(), thursday.Month(), 4, 0, 0, 0, 0, time.Local).ISOWeek()
	// 日期所在周数减去日期所在那个周的周四的月份的4号所在的周数再加一即为本月第几周
	week = week - week1 + 1
	return week
}
func Explain() {
	//┌─────────────second 范围 (0 - 60)
	//│ ┌───────────── min (0 - 59)
	//│ │ ┌────────────── hour (0 - 23)
	//│ │ │ ┌─────────────── day of month (1 - 31)
	//│ │ │ │ ┌──────────────── month (1 - 12)
	//│ │ │ │ │ ┌───────────────── day of week (0 - 6) (0 to 6 are Sunday to
	//│ │ │ │ │ │                  Saturday)
	//│ │ │ │ │ │
	//│ │ │ │ │ │
	//* * * * * *
	//@yearly (or @annually)	每年1月1日午夜跑步一次	0 0 0 1 1 *
	//@monthly	每个月第一天的午夜跑一次	0 0 0 1 * *
	//@weekly	每周周六的午夜运行一次	0 0 0 * * 0
	//@daily (or @midnight)	每天午夜跑一次	0 0 0 * * *
	//@hourly	每小时运行一次	0 0 * * * *
	//@every <duration>	every duration
	//	每隔5秒执行一次：*/5 * * * * ?
	//
	//	每隔1分钟执行一次：0 */1 * * * ?
	//
	//	每天23点执行一次：0 0 23 * * ?
	//
	//	每天凌晨1点执行一次：0 0 1 * * ?
	//
	//	每月1号凌晨1点执行一次：0 0 1 1 * ?
	//
	//每周一和周三晚上22:30: 00 30 22 * * 1,3
	//
	//	在26分、29分、33分执行一次：0 26,29,33 * * * ?
	//
	//	每天的0点、13点、18点、21点都执行一次：0 0 0,13,18,21 * * ?
	//
	//每年三月的星期四的下午14:10和14:40:  00 10,40 14 ? 3 4
}

func ScheduledUpdateTask(a func(), spec string) {
	data, _ := time.LoadLocation("Asia/Shanghai")
	crontab := cron.New(cron.WithSeconds(), cron.WithLocation(data))
	crontab.AddFunc(spec, a)
	go crontab.Start()
	//defer crontab.Stop()
}

func GetCNTimeWeek(data string) string {
	d := ""
	switch data {
	case "Monday":
		d = "星期一"
	case "Tuesday":
		d = "星期二"
	case "Wednesday":
		d = "星期三"
	case "Thursday":
		d = "星期四"
	case "Friday":
		d = "星期五"
	case "Saturday":
		d = "星期六"
	case "Sunday":
		d = "星期日"
	}
	return d
}

func GetCNTimeMonth(data string) string {
	d := ""
	switch data {
	case "January":
		d = "一月"
	case "February":
		d = "二月"
	case "March":
		d = "三月"
	case "April":
		d = "四月"
	case "May":
		d = "五月"
	case "June":
		d = "六月"
	case "July":
		d = "七月"
	case "August":
		d = "八月"
	case "September":
		d = "九月"
	case "October":
		d = "十月"
	case "November":
		d = "十一月"
	case "December":
		d = "十二月"
	}
	return d
}
