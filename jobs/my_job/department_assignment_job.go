package my_job

import (
	"themoment-team/go-hellogsm/internal"
	"themoment-team/go-hellogsm/jobs"
	"themoment-team/go-hellogsm/repository"
)

type ConditionalAssignDepartment struct {
}

type ApplicantAssignDepartment struct {
}

func getDepartmentAssignmentSteps() []jobs.Step {
	return []jobs.Step{&ConditionalAssignDepartment{}, &ApplicantAssignDepartment{}}
}

func BuildDepartmentAssignmentJob(properties internal.ApplicationProperties) *jobs.SimpleJob {
	return jobs.NewSimpleJob(internal.DepartmentAssignmentJob, getDepartmentAssignmentSteps(), nil)
}

func (a *ConditionalAssignDepartment) Processor(context *jobs.BatchContext) {
	giveUpCount := repository.CountByGiveUpApplicant()

	var totalCapacity = 0

	// 중도포기 지원자가 없다면 일반학과배정 진행
	if giveUpCount == 0 {

		totalCapacity = repository.CountByFinalPassApplicant()
		context.Put("totalCapacity", totalCapacity)

		context.Put("SW", 0)
		context.Put("IOT", 0)
		context.Put("AI", 0)

		context.Put("status", "일반학과배정")
	} else {
		// 중도포기 지원자가 있다면 추가학과배정 진행

		totalCapacity = giveUpCount
		context.Put("totalCapacity", totalCapacity)

		sw, iot, ai := repository.QueryByRemainingDepartment()
		context.Put("SW", sw)
		context.Put("IOT", iot)
		context.Put("AI", ai)

		context.Put("status", "추가학과배정")
	}
}

func (s *ApplicantAssignDepartment) Processor(context *jobs.BatchContext) {
}
