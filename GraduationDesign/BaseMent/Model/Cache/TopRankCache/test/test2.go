package main

import (
	"context"
)

var C = context.Background()

func main() {
	//y, m, d, w := strconv.Itoa(time.Now().Year()), time.Now().Month().String(), strconv.Itoa(time.Now().Day()), time.Now().Weekday().String()
	//_, b := time.Now().ISOWeek()
	//fmt.Println(y, m, d, w, b)
	//utils.GetMonthWeek()
	//var id uint
	//id = 0
	//id2 := strconv.Itoa(int(id))
	////key:=y+" "+m+" "+d+" "+w
	//num := Config.DBR.ZAdd(C, Config.Conf.Cache.CacheRankMusic, redis.Z{Member: y + id2, Score: 1}).Err()
	//num := Config.DBR.ZAdd(C, Config.Conf.Cache.CacheRankMusic, redis.Z{Member: y + m + id2, Score: 1}).Err()
	//num := Config.DBR.ZAdd(C, Config.Conf.Cache.CacheRankMusic, redis.Z{Member: y + m + d + w, Score: 1}).Err()
	//num := Config.DBR.ZAdd(C, Config.Conf.Cache.CacheRankMusic, redis.Z{Member: y, Score: 1}).Err()
	//num := Config.DBR.ZAdd(C, Config.Conf.Cache.CacheRankMusic, redis.Z{Member: y, Score: 1}).Err()
	//Config.InitsSQL()
	//TopRankCache.MusicRankAdd(2)
	//TopRankCache.GetMusicRankList()
}
