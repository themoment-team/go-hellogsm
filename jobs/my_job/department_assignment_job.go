package my_job

import (
	"log"
	"themoment-team/go-hellogsm/internal"
	"themoment-team/go-hellogsm/jobs"
	"themoment-team/go-hellogsm/repository"
)

type DepartmentAssignmentStatus string

const (
	NORMAL_ASSIGNED     DepartmentAssignmentStatus = "NORMAL_ASSIGNED"
	ADDITIONAL_ASSIGNED DepartmentAssignmentStatus = "ADDITIONAL_ASSIGNED"
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

		context.Put("status", NORMAL_ASSIGNED)
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

		context.Put("status", ADDITIONAL_ASSIGNED)
	}
}

func (s *ApplicantAssignDepartment) Processor(context *jobs.BatchContext) {
	//totalCapacity := context.Get("totalCapacity")
	status := context.Get("status")

	remainingDeptInterface := context.Get("remainingDepartment")

	remainingDepartment, ok := remainingDeptInterface.(map[string]map[jobs.Major]int)
	if !ok {
		log.Println("remainingDeptInterface의 타입이 올바르지 않습니다.")
		return
	}

	maxDepartment := makeMaxDepartment()

	var finalTestPassApplicants []repository.Applicant
	var err error

	switch status {
	case NORMAL_ASSIGNED:
		err, finalTestPassApplicants = repository.QueryAllByFinalTestPassApplicant()
	case ADDITIONAL_ASSIGNED:
		err, finalTestPassApplicants = repository.QueryAllByAdditionalApplicant()
	}

	if err != nil {
		log.Println(err)
		return
	}

	for _, applicant := range finalTestPassApplicants {

		memberId := applicant.MemberID

		first := applicant.FirstDesiredMajor
		second := applicant.SecondDesiredMajor
		third := applicant.ThirdDesiredMajor

		var decideMajor jobs.Major

		switch applicant.AppliedScreening {
		case jobs.GeneralScreening, jobs.SpecialScreening:
			decideMajor = assign(jobs.NORMAL, first, second, third, remainingDepartment, maxDepartment)
		case jobs.ExtraVeteransScreening, jobs.ExtraAdmissionScreening:
			decideMajor = assign(jobs.EXTRA, first, second, third, remainingDepartment, maxDepartment)
		}

		// 배정된 학과를 반영
		repository.UpdateDecideMajor(decideMajor, memberId)
	}

}

func assign(
	key string, first jobs.Major, second jobs.Major, third jobs.Major,
	remainingDepartment map[string]map[jobs.Major]int, maxDepartment map[string]map[jobs.Major]int,
) jobs.Major {
	if remainingDepartment[key][first] < maxDepartment[key][first] {
		remainingDepartment[key][first]++
		return first
	} else if remainingDepartment[key][second] < maxDepartment[key][second] {
		remainingDepartment[key][second]++
		return second
	} else if remainingDepartment[key][third] < maxDepartment[key][third] {
		remainingDepartment[key][third]++
		return third
	} else {
		panic("발생할 수 없는 상황입니다. 모든 최종 합격자가 학과에 배정되어야 합니다.")
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

func makeMaxDepartment() map[string]map[jobs.Major]int {
	maxDepartment := make(map[string]map[jobs.Major]int)

	maxDepartment[jobs.NORMAL] = make(map[jobs.Major]int)
	maxDepartment[jobs.EXTRA] = make(map[jobs.Major]int)

	maxDepartment[jobs.NORMAL][jobs.SW] = jobs.SWDepartment
	maxDepartment[jobs.NORMAL][jobs.IOT] = jobs.IOTDepartment
	maxDepartment[jobs.NORMAL][jobs.AI] = jobs.AIDepartment
	maxDepartment[jobs.EXTRA][jobs.SW] = 2
	maxDepartment[jobs.EXTRA][jobs.IOT] = 2
	maxDepartment[jobs.EXTRA][jobs.AI] = 2

	return maxDepartment
}
