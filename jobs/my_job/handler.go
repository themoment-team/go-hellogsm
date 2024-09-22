package my_job

import (
	"log"
	"themoment-team/go-hellogsm/internal"
)

func Run(applicationProperties internal.ApplicationProperties, jobs []string) {
	for _, job := range jobs {
		var err error
		switch job {
		case internal.FirstEvaluationJob:
			firstEvaluationJob := BuildFirstEvaluationJob(applicationProperties)
			err = firstEvaluationJob.Start()
		case internal.SecondEvaluationJob:
			secondEvaluationJob := BuildSecondEvaluationJob(applicationProperties)
			err = secondEvaluationJob.Start()
		case internal.DepartmentAssignmentJob:
			departmentAssignmentJob := BuildDepartmentAssignmentJob(applicationProperties)
			err = departmentAssignmentJob.Start()
		default:
			doNothing()
		}

		if err != nil {
			log.Println(err.Error())
			return
		}
	}
}

func doNothing() {

}
