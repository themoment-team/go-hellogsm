package main

import (
	"themoment-team/go-hellogsm/configs"
	"themoment-team/go-hellogsm/internal"
	"themoment-team/go-hellogsm/jobs/my_job"
	"themoment-team/go-hellogsm/service"
)

func main() {

	// 애플리케이션 실행 인자들을 확인해서 알맞게 초기화한다.
	internal.ApplicationArgsProcessor()
	// 프로그램 p 사용할 property들을 초기화한다.
	properties := internal.InitApplicationProperties(internal.GetActiveProfile())
	// db 싱글톤 인스턴스를 생성한다.
	configs.CreateMysqlDB(configs.CreateMysqlDsn(properties.Mysql))
	// job 을 실행하기 전에 필요한 third-party 를 ping 한다.
	service.Ping()
	// job 을 실행한다.
	my_job.Run(properties, internal.GetJobs())

}
