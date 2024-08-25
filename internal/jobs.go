package internal

import (
	"fmt"
	"gorm.io/gorm/utils"
	"os"
	"strings"
)

const (
	FirstEvaluationJob      = "firstEvaluationJob"
	SecondEvaluationJob     = "secondEvaluationJob"
	DepartmentAssignmentJob = "departmentAssignmentJob"
)

var (
	MyJobs = []string{FirstEvaluationJob, SecondEvaluationJob, DepartmentAssignmentJob}
)

func SetJobs(jobsAsString string) {
	os.Setenv("jobs", strings.Join(validateJobs(strings.Split(jobsAsString, ",")), ","))
}

func GetJobs() []string {
	return strings.Split(os.Getenv("jobs"), ",")
}

func validateJobs(jobs []string) []string {
	var validatedJobs []string
	for _, job := range jobs {
		if isAvailableJob(job) {
			validatedJobs = append(validatedJobs, job)
		} else {
			fmt.Println(fmt.Sprintf("[%s] 에 해당하는 Job 은 존재하지 않습니다. 무시됨.", job))
		}
	}

	if 0 > len(validatedJobs) {
		panic("실행 가능한 Job이 없습니다.")
	}

	return validatedJobs
}

func isAvailableJob(job string) bool {
	return utils.Contains(MyJobs, job)
}
