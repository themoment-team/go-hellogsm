package main

import (
	"themoment-team/go-hellogsm/internal"
	"themoment-team/go-hellogsm/jobs/my_job"
	"themoment-team/go-hellogsm/services"
)

func main() {

	// 필수적인 meta 들을 프로그램에 초기화하고 검증한다.
	internal.ProgramEssentialMeta()
	// profile 에 알맞는 application property를 세트한다.
	properties := internal.ExtractApplicationProperties(internal.GetActiveProfile())
	// job 을 실행하기 전에 필요한 third-party 를 ping 한다.
	services.Ping(properties)
	// job 을 실행한다.
	my_job.Run(properties, internal.GetJobs())

}
