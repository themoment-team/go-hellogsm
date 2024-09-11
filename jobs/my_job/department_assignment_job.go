package my_job

import (
	"log"
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

	// 모집할 총 인원, 각 학과 별 남은 자리, 이후 실행할 작업에 대한 정보를 context에 담는다

	// 중도포기 지원자가 없다면 일반학과배정 진행
	if giveUpCount == 0 {

		totalCapacity = repository.CountByFinalPassApplicant()
		context.Put("totalCapacity", totalCapacity)

		remainingDepartment := makeRemainingDepartment(0, 0, 0, 0, 0, 0)
		context.Put("remainingDepartment", remainingDepartment)

		context.Put("status", "일반학과배정")
	} else {
		// 중도포기 지원자가 있다면 추가학과배정 진행

		totalCapacity = giveUpCount
		context.Put("totalCapacity", totalCapacity)

		normalSwCount, normalIotCount, normalAiCount := repository.QueryByScrenningsRemainingDepartment(jobs.GeneralScreening, jobs.SpecialScreening)
		extraSwCount, extraIotCount, extraAiCount := repository.QueryByScrenningsRemainingDepartment(jobs.ExtraAdmissionScreening, jobs.ExtraVeteransScreening)
		remainingDepartment := makeRemainingDepartment(
			jobs.SWDepartment-normalSwCount, jobs.IOTDepartment-normalIotCount, jobs.AIDepartment-normalAiCount,
			jobs.ExtraDepartment-extraSwCount, jobs.ExtraDepartment-extraIotCount, jobs.ExtraDepartment-extraAiCount)
		context.Put("remainingDepartment", remainingDepartment)

		context.Put("status", "추가학과배정")
	}
}

func makeRemainingDepartment(
	normalSw int, normalIot int, normalAi int,
	ExtraSw int, ExtraIot int, ExtraAi int,
) map[string]map[jobs.Major]int {
	remainingDepartment := make(map[string]map[jobs.Major]int)

	remainingDepartment[jobs.NORMAL] = make(map[jobs.Major]int)
	remainingDepartment[jobs.EXTRA] = make(map[jobs.Major]int)

	remainingDepartment[jobs.NORMAL][jobs.SW] = normalSw
	remainingDepartment[jobs.NORMAL][jobs.IOT] = normalIot
	remainingDepartment[jobs.NORMAL][jobs.AI] = normalAi
	remainingDepartment[jobs.EXTRA][jobs.SW] = ExtraSw
	remainingDepartment[jobs.EXTRA][jobs.IOT] = ExtraIot
	remainingDepartment[jobs.EXTRA][jobs.AI] = ExtraAi

	return remainingDepartment
}

func (s *ApplicantAssignDepartment) Processor(context *jobs.BatchContext) {
	//totalCapacity := context.Get("totalCapacity")
	//status := context.Get("status")

	remainingDeptInterface := context.Get("remainingDepartment")

	remainingDepartment, ok := remainingDeptInterface.(map[string]map[jobs.Major]int)
	if !ok {
		log.Println("remainingDeptInterface의 타입이 올바르지 않습니다.")
		return
	}

	maxDepartment := make(map[jobs.Major]int)
	maxDepartment[jobs.SW] = jobs.SWDepartment
	maxDepartment[jobs.IOT] = jobs.IOTDepartment
	maxDepartment[jobs.AI] = jobs.AIDepartment

	_, finalTestPassApplicants := repository.QueryAllByFinalTestPassApplicant()

	for _, applicant := range finalTestPassApplicants {

		first := applicant.FirstDesiredMajor
		second := applicant.SecondDesiredMajor
		third := applicant.ThirdDesiredMajor

		switch applicant.AppliedScreening {
		case jobs.GeneralScreening, jobs.SpecialScreening:
			assign(jobs.NORMAL, first, second, third, remainingDepartment, maxDepartment)
		case jobs.ExtraVeteransScreening, jobs.ExtraAdmissionScreening:
			assign(jobs.EXTRA, first, second, third, remainingDepartment, maxDepartment)
		}

	}

}

func assign(
	key string, first jobs.Major, second jobs.Major, third jobs.Major,
	remainingDepartment map[string]map[jobs.Major]int, maxDepartment map[jobs.Major]int,
) {

}
