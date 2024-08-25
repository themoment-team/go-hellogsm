package jobs

import "themoment-team/go-hellogsm/internal"

type DepartmentAssignmentJob struct{}

func (d *DepartmentAssignmentJob) Execute(internal.ApplicationProperties) error {
	return nil
}
