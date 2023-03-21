package Config

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"log"
	"os"
	"time"
)

var err error

func InitsPSQL() {
	dsn := fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s", Conf.PostgreSQL.Host, Conf.PostgreSQL.Port, Conf.PostgreSQL.DbUser, Conf.PostgreSQL.DbName, Conf.PostgreSQL.DbPassword)
	log.Println("PsqlInfo:", dsn)
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold: time.Second, // 慢 Config 阈值
			LogLevel:      logger.Info, // Log level
			Colorful:      true,        // 禁用彩色打印
		},
	)
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		// gorm日志模式
		Logger: newLogger,
		// 外键约束
		DisableForeignKeyConstraintWhenMigrating: true,
		// 禁用默认事务（提高运行速度）
		SkipDefaultTransaction: true,
		NamingStrategy: schema.NamingStrategy{
			// 使用单数表名，启用该选项，此时，`User` 的表名应该是 `user`
			SingularTable: true,
		},
	})
	if err != nil {
		log.Fatal("无法连接到PSQL,请检查参数", err)
	}
	//DB.AutoMigrate(&Model.User{}, &Model.UserInfo{})
}
