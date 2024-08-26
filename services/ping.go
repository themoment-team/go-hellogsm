package services

import (
	"fmt"
	"log"
	"themoment-team/go-hellogsm/repository"
)

func Ping() {
	mysqlPing()
}

func mysqlPing() {
	result := repository.SelectOne()
	if result == 1 {
		log.Println(fmt.Sprintf("mysql ping 결과: [%d] 성공", result))
	} else {
		panic(fmt.Sprintf("mysql ping 결과: [%d] 실패", result))
	}
}
