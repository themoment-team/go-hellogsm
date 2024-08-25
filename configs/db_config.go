package configs

import (
	"gorm.io/gorm/logger"
	"log"
	"os"
	"themoment-team/go-hellogsm/internal"
	"time"
)

func GetMysqlDsn(properties internal.MysqlProperties) string {
	return ""
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
