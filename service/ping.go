package service

import (
	"fmt"
	"log"
	"themoment-team/go-hellogsm/repository"
)

func Ping() {
	mysqlPing()
	relayApiPing()
}

func mysqlPing() {
	result := repository.SelectOne()
	if result != 1 {
		panic(fmt.Sprintf("mysql ping 결과: [%d] 실패", result))
	}
	log.Printf("mysql ping 결과: [%d] 성공", result)
}

func relayApiPing() {
	err := PingRelayApi()
	if err != nil {
		panic(err.Error())
	}
	log.Printf("relay-api ping 결과: 성공")
}
