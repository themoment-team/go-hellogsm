package services

import (
	"fmt"
	"themoment-team/go-hellogsm/configs"
	"themoment-team/go-hellogsm/internal"
	"themoment-team/go-hellogsm/repository"
)

func Ping(properties internal.ApplicationProperties) {
	mysqlPing(properties.Mysql)
}

func mysqlPing(properties internal.MysqlProperties) {
	db, err := configs.CreateMysqlDB(configs.GetMysqlDsn(properties))
	if err != nil {
		panic("DB 초기화 실패.")
	}

	result := repository.SelectOne(db)
	if result == 1 {
		fmt.Println(fmt.Sprintf("mysql ping 결과: [%d] 성공", result))
	} else {
		panic(fmt.Sprintf("mysql ping 결과: [%d] 실패", result))
	}
}
