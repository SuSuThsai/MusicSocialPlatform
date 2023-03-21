package Config

import (
	"fmt"
	"github.com/redis/go-redis/v9"
	"log"
)

func InitRedis() {
	DBR = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", Conf.Redis.HostR, Conf.Redis.PortR),
		Password: Conf.Redis.DbPassWordR,
		DB:       Conf.Redis.DBModel,
		//PoolSize:     100,
		//MaxIdleConns: 100,
	})
	DBR2 = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", Conf.Redis.HostR, Conf.Redis.PortR),
		Password: Conf.Redis.DbPassWordR,
		DB:       Conf.Redis.DBModel + 1,
		//PoolSize:     100,
		//MaxIdleConns: 100,
	})
	log.Println("RedisInfo:", DBR)
}
