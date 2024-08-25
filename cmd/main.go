package main

import (
	"themoment-team/go-hellogsm/internal"
	"themoment-team/go-hellogsm/jobs"
)

func main() {

	// 필수적인 meta 들을 프로그램에 초기화하고 검증한다.
	internal.ProgramEssentialMeta()
	// profile 에 알맞는 application property를 세트한다.
	properties := internal.ExtractApplicationProperties(internal.GetActiveProfile())
	// job 을 실행한다.
	jobs.Run(properties, internal.GetJobs())

	//dsn := "root:1234@tcp(127.0.0.1:3306)/hellogsm_tmp?charset=utf8mb4&parseTime=True&loc=Local"
	//
	//db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{Logger: configs.GetMyDbLoggerConfig()})
	//
	//if err != nil {
	//	panic("failed to connect mysql")
	//}
	//
	//// 테스트로 select 1
	//var result int
	//tx := db.Raw("select 1").Scan(&result)
	//if tx.Error != nil {
	//	panic("sql execute error")
	//}
	//fmt.Println(result)
	//
	//var employees Employees
	//db.First(&employees)
	//
	//fmt.Println(employees.FirstName)
}

type Employees struct {
	ID        int64
	FirstName string
	LastName  string
	Email     string
	HireDate  string
	Salary    float32
}
