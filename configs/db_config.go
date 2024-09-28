package configs

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"sync"
	"themoment-team/go-hellogsm/internal"
	"time"
)

// MyDB 는 CreateMysqlDB 에서 singleton 인스턴스를 생성한다.
var MyDB gorm.DB
var once sync.Once

func CreateMysqlDB(dsn string) {
	once.Do(func() {
		db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{Logger: GetMyDbLogger(), SkipDefaultTransaction: true})
		if err != nil {
			panic("DB 인스턴스화 실패")
		}
		MyDB = *db
	})
}

func CreateMysqlDsn(properties internal.MysqlProperties) string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		properties.Username,
		properties.Password,
		properties.Host,
		properties.Port,
		properties.Database,
	)
}

func GetMyDbLogger() logger.Interface {
	return logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second, // Slow SQL threshold
			LogLevel:                  logger.Info, // Log level
			IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound error for logger
			ParameterizedQueries:      true,        // Don't include params in the SQL log
			Colorful:                  false,       // Disable color
		},
	)
}
