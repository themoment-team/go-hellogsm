package my_job

import (
	"errors"
	"gorm.io/gorm"
	"log"
	"math"
	e "themoment-team/go-hellogsm/error"
	"themoment-team/go-hellogsm/internal"
	"themoment-team/go-hellogsm/jobs"
	"themoment-team/go-hellogsm/repository"
	"themoment-team/go-hellogsm/types"
)

type MajorAssignmentStatus string

const (
	NORMAL_ASSIGNED     MajorAssignmentStatus = "NORMAL_ASSIGNED"
	ADDITIONAL_ASSIGNED MajorAssignmentStatus = "ADDITIONAL_ASSIGNED"
)

type ConditionalAssignMajorStep struct {
}

type ApplicantAssignMajorStep struct {
}

func getMajorAssignmentSteps() []jobs.Step {
	return []jobs.Step{&ConditionalAssignMajorStep{}, &ApplicantAssignMajorStep{}}
}

func BuildMajorAssignmentJob() *jobs.SimpleJob {
	return jobs.NewSimpleJob(internal.MajorAssignmentJob, getMajorAssignmentSteps(), nil)
}

func (a *ConditionalAssignMajorStep) Processor(context *jobs.BatchContext, db *gorm.DB) error {
	giveUpCount := repository.CountByGiveUpApplicant()

	// 각 학과 별 남은 자리, 이후 실행할 작업에 대한 정보를 context에 담는다

	// 중도포기 지원자가 없다면 일반학과배정 진행
	if giveUpCount == 0 {
		assignedMajor := makeAssignedMajor(0, 0, 0, 0, 0, 0)
		context.Put("assignedMajor", assignedMajor)

		context.Put("status", NORMAL_ASSIGNED)
		log.Println("일반학과배정이 진행됩니다.")
	} else {
		// 중도포기 지원자가 있다면 추가학과배정 진행
		normalSwCount, normalIotCount, normalAiCount := repository.QueryByScrenningsAssignedMajor(types.GeneralScreening, types.SpecialScreening)
		extraSwCount, extraIotCount, extraAiCount := repository.QueryByScrenningsAssignedMajor(types.ExtraAdmissionScreening, types.ExtraVeteransScreening)
		assignedMajor := makeAssignedMajor(
			types.SWMajor-normalSwCount, types.IOTMajor-normalIotCount, types.AIMajor-normalAiCount,
			types.ExtraMajor-extraSwCount-normalSwCount, types.ExtraMajor-extraIotCount-normalIotCount, types.ExtraMajor-extraAiCount-normalAiCount)
		context.Put("assignedMajor", assignedMajor)

		context.Put("status", ADDITIONAL_ASSIGNED)
		log.Println(
			"일반전형, 사회통합전형: ", normalSwCount+normalIotCount+normalAiCount, "개, ",
			"정원 외 특별전형: ", extraSwCount+extraIotCount+extraAiCount, "개",
			"의 남은 자리를 대상으로 추가학과배정이 진행됩니다.")
	}

	return nil
}

var updatedMemberIds []int

func (s *ApplicantAssignMajorStep) Processor(context *jobs.BatchContext, db *gorm.DB) error {

	statusInterface := context.Get("status")
	assignedMajorInterface := context.Get("assignedMajor")

	// 각 학과 별 남은 자리
	assignedMajor, ok := assignedMajorInterface.(map[string]map[types.Major]int)
	if !ok {
		log.Println("assignedMajorInterface의 타입이 올바르지 않습니다.")
		return nil
	}

	status, ok := statusInterface.(MajorAssignmentStatus)
	if !ok {
		log.Println("status의 타입이 올바르지 않습니다.")
		return nil
	}

	// 각 학과 별 정원
	maxMajor := makeMaxMajor()

	var targetApplicant []types.Applicant
	var err error

	// 학과배정 상태에 맞는 base data set 초기화 (일반학과배정, 추가모집배정)
	switch status {
	case NORMAL_ASSIGNED:
		err, targetApplicant = repository.QueryAllByFinalTestPassApplicant()
	case ADDITIONAL_ASSIGNED:
		err, targetApplicant = repository.QueryAllByAdditionalApplicant()
	}

	if err != nil {
		log.Println(err)
		return nil
	}

	for _, applicant := range targetApplicant {

		memberId := applicant.MemberID

		first := applicant.FirstDesiredMajor
		second := applicant.SecondDesiredMajor
		third := applicant.ThirdDesiredMajor

		var decideMajor types.Major
		var err error

		switch applicant.AppliedScreening {
		case types.GeneralScreening, types.SpecialScreening:
			decideMajor, err = assign(types.NORMAL, first, second, third, assignedMajor, maxMajor, status)
		case types.ExtraVeteransScreening, types.ExtraAdmissionScreening:
			decideMajor, err = assign(types.EXTRA, first, second, third, assignedMajor, maxMajor, status)
		}

		var rollBackNeededError e.RollbackNeededError
		if err != nil {
			if errors.As(err, &rollBackNeededError) {
				// 학과 배정 중 에러가 발생해 롤백
				return err
			} else {
				break
			}
		}

		// 배정된 학과를 반영
		updatedMemberIds = append(updatedMemberIds, memberId)
		updateMajorError := repository.UpdateDecideMajor(db, decideMajor, memberId)
		if updateMajorError != nil {
			// 배정된 학과 반영 중 에러가 발생해 롤백
			return updateMajorError
		}
	}

	if status == NORMAL_ASSIGNED {
		log.Println(len(updatedMemberIds), "명이 학과에 배정되었습니다. 일반학과배정 종료")
	} else {
		log.Println(len(updatedMemberIds), "명이 추가로 학과에 배정되었습니다. 추가학과배정 종료")
	}

	return nil
}

func assign(
	key string, first types.Major, second types.Major, third types.Major,
	assignedMajor map[string]map[types.Major]int, maxMajor map[string]map[types.Major]int, status MajorAssignmentStatus,
) (types.Major, error) {
	if assignedMajor[key][first] < maxMajor[key][first] {
		assignedMajor[key][first]++
		return first, nil
	} else if assignedMajor[key][second] < maxMajor[key][second] {
		assignedMajor[key][second]++
		return second, nil
	} else if assignedMajor[key][third] < maxMajor[key][third] {
		assignedMajor[key][third]++
		return third, nil
	} else {
		if status == NORMAL_ASSIGNED {
			return "", e.WrapExpectedActualIsDiffError("일반학과배정시에는 발생할 수 없는 상황입니다. 모든 최종 합격자가 학과에 배정되어야 합니다.")
		} else {
			return "", errors.New("학과가 배정되지 않았습니다")
		}
	}
}

func makeAssignedMajor(
	normalSw int, normalIot int, normalAi int,
	ExtraSw int, ExtraIot int, ExtraAi int,
) map[string]map[types.Major]int {
	assignedMajor := make(map[string]map[types.Major]int)

	assignedMajor[types.NORMAL] = make(map[types.Major]int)
	assignedMajor[types.EXTRA] = make(map[types.Major]int)

	assignedMajor[types.NORMAL][types.SW] = normalSw
	assignedMajor[types.NORMAL][types.IOT] = normalIot
	assignedMajor[types.NORMAL][types.AI] = normalAi
	assignedMajor[types.EXTRA][types.SW] = int(math.Max(0, float64(ExtraSw)))
	assignedMajor[types.EXTRA][types.IOT] = int(math.Max(0, float64(ExtraIot)))
	assignedMajor[types.EXTRA][types.AI] = int(math.Max(0, float64(ExtraAi)))

	return assignedMajor
}

func makeMaxMajor() map[string]map[types.Major]int {
	maxMajor := make(map[string]map[types.Major]int)

	maxMajor[types.NORMAL] = make(map[types.Major]int)
	maxMajor[types.EXTRA] = make(map[types.Major]int)

	maxMajor[types.NORMAL][types.SW] = types.SWMajor
	maxMajor[types.NORMAL][types.IOT] = types.IOTMajor
	maxMajor[types.NORMAL][types.AI] = types.AIMajor
	maxMajor[types.EXTRA][types.SW] = 2
	maxMajor[types.EXTRA][types.IOT] = 2
	maxMajor[types.EXTRA][types.AI] = 2

	return maxMajor
}
