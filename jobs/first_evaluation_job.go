package jobs

import (
	"fmt"
	"themoment-team/go-hellogsm/internal"
)

type FirstEvaluationJob struct{}

func (f *FirstEvaluationJob) Execute(properties internal.ApplicationProperties) error {
	// 구현 내용
	fmt.Println("hello first evaluation job :)")
	fmt.Println(fmt.Sprintf("%s arrived !", properties.Mysql.Username))
	return nil
}
