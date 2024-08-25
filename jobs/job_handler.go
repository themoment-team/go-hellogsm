package jobs

import (
	"fmt"
	"themoment-team/go-hellogsm/internal"
)

func Run(applicationProperties internal.ApplicationProperties, jobs []string) {
	for _, job := range jobs {
		switch job {
		case internal.FirstEvaluationJob:
			targetJob := FirstEvaluationJob{}
			err := targetJob.Execute(applicationProperties)
			if err != nil {
				fmt.Println(err)
			}
		case internal.SecondEvaluationJob:
			targetJob := SecondEvaluationJob{}
			err := targetJob.Execute(applicationProperties)
			if err != nil {
				fmt.Println(err)
			}
		case internal.DepartmentAssignmentJob:
			targetJob := DepartmentAssignmentJob{}
			err := targetJob.Execute(applicationProperties)
			if err != nil {
				fmt.Println(err)
			}
		default:
			doNothing()
		}
	}
}

func doNothing() {

}
