package model

import (
	"log"
	"os"
	"time"

	"gorm.io/gorm/logger"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func NewSqliteConn(dsn string) *gorm.DB {
	var err error
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // 设置日志输出
		logger.Config{
			SlowThreshold: 10 * time.Second, // 慢查询阈值，超过该阈值的查询将被认为是慢查询
			LogLevel:      logger.Info,      // 日志级别
			Colorful:      true,             // 是否启用彩色日志
		},
	)
	conn, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{
		Logger: newLogger, // 设置 Logger
	})
	if err != nil {
		log.Fatal()
	}

	return conn
}
