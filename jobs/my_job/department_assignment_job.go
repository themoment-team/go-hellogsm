package my_job

import (
	"errors"
	"log"
	"math"
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

	// 각 학과 별 남은 자리, 이후 실행할 작업에 대한 정보를 context에 담는다

	// 중도포기 지원자가 없다면 일반학과배정 진행
	if giveUpCount == 0 {
		remainingDepartment := makeRemainingDepartment(0, 0, 0, 0, 0, 0)
		context.Put("remainingDepartment", remainingDepartment)

		context.Put("status", NORMAL_ASSIGNED)
	} else {
		// 중도포기 지원자가 있다면 추가학과배정 진행
		normalSwCount, normalIotCount, normalAiCount := repository.QueryByScrenningsRemainingDepartment(jobs.GeneralScreening, jobs.SpecialScreening)
		extraSwCount, extraIotCount, extraAiCount := repository.QueryByScrenningsRemainingDepartment(jobs.ExtraAdmissionScreening, jobs.ExtraVeteransScreening)
		remainingDepartment := makeRemainingDepartment(
			jobs.SWDepartment-normalSwCount, jobs.IOTDepartment-normalIotCount, jobs.AIDepartment-normalAiCount,
			jobs.ExtraDepartment-extraSwCount-normalSwCount, jobs.ExtraDepartment-extraIotCount-normalIotCount, jobs.ExtraDepartment-extraAiCount-normalAiCount)
		context.Put("remainingDepartment", remainingDepartment)

		context.Put("status", ADDITIONAL_ASSIGNED)
	}
}

func (s *ApplicantAssignDepartment) Processor(context *jobs.BatchContext) {

	statusInterface := context.Get("status")
	remainingDeptInterface := context.Get("remainingDepartment")

	// 각 학과 별 남은 자리
	remainingDepartment, ok := remainingDeptInterface.(map[string]map[jobs.Major]int)
	if !ok {
		log.Println("remainingDeptInterface의 타입이 올바르지 않습니다.")
		return
	}

	status, ok := statusInterface.(DepartmentAssignmentStatus)
	if !ok {
		log.Println("status의 타입이 올바르지 않습니다.")
		return
	}

	log.Println(status, "가 진행됩니다.")

	// 각 학과 별 정원
	maxDepartment := makeMaxDepartment()

	var finalTestPassApplicants []repository.Applicant
	var err error

	// 학과배정 상태에 맞는 base data set 초기화 (일반학과배정, 추가모집배정)
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
		var err error

		switch applicant.AppliedScreening {
		case jobs.GeneralScreening, jobs.SpecialScreening:
			decideMajor, err = assign(jobs.NORMAL, first, second, third, remainingDepartment, maxDepartment, status)
		case jobs.ExtraVeteransScreening, jobs.ExtraAdmissionScreening:
			decideMajor, err = assign(jobs.EXTRA, first, second, third, remainingDepartment, maxDepartment, status)
		}

		// 배정된 학과를 반영
		if err == nil {
			repository.UpdateDecideMajor(decideMajor, memberId)
		} else {
			return
		}
	}

	if status == NORMAL_ASSIGNED {
		log.Println("일반학과배정 종료")
	} else {
		log.Println("추가학과배정 종료")
	}

}

func assign(
	key string, first jobs.Major, second jobs.Major, third jobs.Major,
	remainingDepartment map[string]map[jobs.Major]int, maxDepartment map[string]map[jobs.Major]int, status DepartmentAssignmentStatus,
) (jobs.Major, error) {
	if remainingDepartment[key][first] < maxDepartment[key][first] {
		remainingDepartment[key][first]++
		return first, nil
	} else if remainingDepartment[key][second] < maxDepartment[key][second] {
		remainingDepartment[key][second]++
		return second, nil
	} else if remainingDepartment[key][third] < maxDepartment[key][third] {
		remainingDepartment[key][third]++
		return third, nil
	} else {
		if status == NORMAL_ASSIGNED {
			panic("일반학과배정시에는 발생할 수 없는 상황입니다. 모든 최종 합격자가 학과에 배정되어야 합니다.")
		} else {
			return "", errors.New("학과가 배정되지 않았습니다")
		}
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
	remainingDepartment[jobs.EXTRA][jobs.SW] = int(math.Max(0, float64(ExtraSw)))
	remainingDepartment[jobs.EXTRA][jobs.IOT] = int(math.Max(0, float64(ExtraIot)))
	remainingDepartment[jobs.EXTRA][jobs.AI] = int(math.Max(0, float64(ExtraAi)))

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
