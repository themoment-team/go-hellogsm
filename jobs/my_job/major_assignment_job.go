package my_job

import (
	"errors"
	"log"
	"math"
	"themoment-team/go-hellogsm/internal"
	"themoment-team/go-hellogsm/jobs"
	"themoment-team/go-hellogsm/repository"
)

type MajorAssignmentStatus string

const (
	NORMAL_ASSIGNED     MajorAssignmentStatus = "NORMAL_ASSIGNED"
	ADDITIONAL_ASSIGNED MajorAssignmentStatus = "ADDITIONAL_ASSIGNED"
)

type ConditionalAssignMajor struct {
}

type ApplicantAssignMajor struct {
}

func getMajorAssignmentSteps() []jobs.Step {
	return []jobs.Step{&ConditionalAssignMajor{}, &ApplicantAssignMajor{}}
}

func BuildMajorAssignmentJob() *jobs.SimpleJob {
	return jobs.NewSimpleJob(internal.MajorAssignmentJob, getMajorAssignmentSteps(), nil)
}

func (a *ConditionalAssignMajor) Processor(context *jobs.BatchContext) {
	giveUpCount := repository.CountByGiveUpApplicant()

	// 각 학과 별 남은 자리, 이후 실행할 작업에 대한 정보를 context에 담는다

	// 중도포기 지원자가 없다면 일반학과배정 진행
	if giveUpCount == 0 {
		remainingMajor := makeRemainingMajor(0, 0, 0, 0, 0, 0)
		context.Put("remainingMajor", remainingMajor)

		context.Put("status", NORMAL_ASSIGNED)
	} else {
		// 중도포기 지원자가 있다면 추가학과배정 진행
		normalSwCount, normalIotCount, normalAiCount := repository.QueryByScrenningsRemainingMajor(jobs.GeneralScreening, jobs.SpecialScreening)
		extraSwCount, extraIotCount, extraAiCount := repository.QueryByScrenningsRemainingMajor(jobs.ExtraAdmissionScreening, jobs.ExtraVeteransScreening)
		remainingMajor := makeRemainingMajor(
			jobs.SWMajor-normalSwCount, jobs.IOTMajor-normalIotCount, jobs.AIMajor-normalAiCount,
			jobs.ExtraMajor-extraSwCount-normalSwCount, jobs.ExtraMajor-extraIotCount-normalIotCount, jobs.ExtraMajor-extraAiCount-normalAiCount)
		context.Put("remainingMajor", remainingMajor)

		context.Put("status", ADDITIONAL_ASSIGNED)
	}
}

var updatedMemberIds []int

func (s *ApplicantAssignMajor) Processor(context *jobs.BatchContext) {

	statusInterface := context.Get("status")
	remainingMajorInterface := context.Get("remainingMajor")

	// 각 학과 별 남은 자리
	remainingMajor, ok := remainingMajorInterface.(map[string]map[jobs.Major]int)
	if !ok {
		log.Println("remainingDeptInterface의 타입이 올바르지 않습니다.")
		return
	}

	status, ok := statusInterface.(MajorAssignmentStatus)
	if !ok {
		log.Println("status의 타입이 올바르지 않습니다.")
		return
	}

	log.Println(status, "가 진행됩니다.")

	// 각 학과 별 정원
	maxMajor := makeMaxMajor()

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
			decideMajor, err = assign(jobs.NORMAL, first, second, third, remainingMajor, maxMajor, status)
		case jobs.ExtraVeteransScreening, jobs.ExtraAdmissionScreening:
			decideMajor, err = assign(jobs.EXTRA, first, second, third, remainingMajor, maxMajor, status)
		}

		// 배정된 학과를 반영
		if err == nil {
			updatedMemberIds = append(updatedMemberIds, memberId)
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
	remainingMajor map[string]map[jobs.Major]int, maxMajor map[string]map[jobs.Major]int, status MajorAssignmentStatus,
) (jobs.Major, error) {
	if remainingMajor[key][first] < maxMajor[key][first] {
		remainingMajor[key][first]++
		return first, nil
	} else if remainingMajor[key][second] < maxMajor[key][second] {
		remainingMajor[key][second]++
		return second, nil
	} else if remainingMajor[key][third] < maxMajor[key][third] {
		remainingMajor[key][third]++
		return third, nil
	} else {
		if status == NORMAL_ASSIGNED {
			log.Println("DecideMajor 롤백을 진행합니다.")
			repository.RollBackDecideMajor(updatedMemberIds)
			panic("일반학과배정시에는 발생할 수 없는 상황입니다. 모든 최종 합격자가 학과에 배정되어야 합니다.")
		} else {
			return "", errors.New("학과가 배정되지 않았습니다")
		}
	}
}

func makeRemainingMajor(
	normalSw int, normalIot int, normalAi int,
	ExtraSw int, ExtraIot int, ExtraAi int,
) map[string]map[jobs.Major]int {
	remainingMajor := make(map[string]map[jobs.Major]int)

	remainingMajor[jobs.NORMAL] = make(map[jobs.Major]int)
	remainingMajor[jobs.EXTRA] = make(map[jobs.Major]int)

	remainingMajor[jobs.NORMAL][jobs.SW] = normalSw
	remainingMajor[jobs.NORMAL][jobs.IOT] = normalIot
	remainingMajor[jobs.NORMAL][jobs.AI] = normalAi
	remainingMajor[jobs.EXTRA][jobs.SW] = int(math.Max(0, float64(ExtraSw)))
	remainingMajor[jobs.EXTRA][jobs.IOT] = int(math.Max(0, float64(ExtraIot)))
	remainingMajor[jobs.EXTRA][jobs.AI] = int(math.Max(0, float64(ExtraAi)))

	return remainingMajor
}

func makeMaxMajor() map[string]map[jobs.Major]int {
	maxMajor := make(map[string]map[jobs.Major]int)

	maxMajor[jobs.NORMAL] = make(map[jobs.Major]int)
	maxMajor[jobs.EXTRA] = make(map[jobs.Major]int)

	maxMajor[jobs.NORMAL][jobs.SW] = jobs.SWMajor
	maxMajor[jobs.NORMAL][jobs.IOT] = jobs.IOTMajor
	maxMajor[jobs.NORMAL][jobs.AI] = jobs.AIMajor
	maxMajor[jobs.EXTRA][jobs.SW] = 2
	maxMajor[jobs.EXTRA][jobs.IOT] = 2
	maxMajor[jobs.EXTRA][jobs.AI] = 2

	return maxMajor
}
