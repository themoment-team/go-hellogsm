package my_job

import (
	"themoment-team/go-hellogsm/internal"
)

func Run(applicationProperties internal.ApplicationProperties, jobs []string) {
	for _, job := range jobs {
		switch job {
		case internal.FirstEvaluationJob:
			firstEvaluationJob := BuildFirstEvaluationJob(applicationProperties)
			firstEvaluationJob.Start()
		case internal.SecondEvaluationJob:
			secondEvaluationJob := BuildSecondEvaluationJob(applicationProperties)
			secondEvaluationJob.Start()
		case internal.DepartmentAssignmentJob:
			departmentAssignmentJob := BuildDepartmentAssignmentJob(applicationProperties)
			departmentAssignmentJob.Start()
		default:
			doNothing()
		}
	}
}

func doNothing() {

}
