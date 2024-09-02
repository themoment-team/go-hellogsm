package my_job

import (
	"fmt"
	"log"
	"themoment-team/go-hellogsm/internal"
	"themoment-team/go-hellogsm/jobs"
	"themoment-team/go-hellogsm/repository"
)

// ApplicationScreeningEvaluationStep 적용 전형 평가를 하는 Step 이다.
type ApplicationScreeningEvaluationStep struct {
}

// CalculationOfSuccessfulApplicantByScreeningProcessStep 전형별 합격자를 산출하는 Step 이다.
type CalculationOfSuccessfulApplicantByScreeningProcessStep struct {
}

// job 에 필요한 step 들을 반환한다.
func getSteps() []jobs.Step {
	return []jobs.Step{&ApplicationScreeningEvaluationStep{}, &CalculationOfSuccessfulApplicantByScreeningProcessStep{}}
}

// BuildFirstEvaluationJob 1차 평가 배치 Job을 생성한다.
func BuildFirstEvaluationJob(properties internal.ApplicationProperties) *jobs.SimpleJob {
	return jobs.NewSimpleJob(internal.FirstEvaluationJob, getSteps(), nil)
}

func (s *ApplicationScreeningEvaluationStep) Processor(context *jobs.BatchContext) {
	// 시작 전 데이터 검증
	err := canNextProcess(jobs.SpecialScreening, jobs.SpecialScreening)
	if err != nil {
		panic(err.Error())
	}

	// 특별전형을 희망하는 지원자의 적용 전형 처리
	specialWantedCount := repository.CountOneseoByWantedScreening(string(jobs.SpecialScreening))
	if specialWantedCount > jobs.SpecialSuccessfulApplicantOf1E {
		log.Printf("특별전형 희망 지원자 수가 [%d]명을 초과하여 하위 [%d]명을 일반전형으로 적용합니다.", jobs.SpecialSuccessfulApplicantOf1E, specialWantedCount-jobs.SpecialSuccessfulApplicantOf1E)
		applyScreeningSpecialPartial()
	} else {
		log.Printf("특별전형 희망 지원자 수 [%d]명을 모두 특별전형으로 적용합니다.", specialWantedCount)
		applyScreeningSpecialAll()
	}

	err = canNextProcess(jobs.GeneralScreening, jobs.SpecialScreening)
	if err != nil {
		panic(err.Error())
		// tx rollback
	}

	// 일반전형을 희망하는 지원자의 적용 전형 처리
}

// 비정상적인 상황에 대한 방어와 데이터 꼬임을 방지하기 위해 다음 스텝으로 넘어가기전 검증한다.
func canNextProcess(to jobs.Screening, from jobs.Screening) error {
	if (to == from) && (to == jobs.SpecialScreening) {
		return validateAllAppliedScreening()
	} else if (to == jobs.GeneralScreening) && (from == jobs.SpecialScreening) {
		return validateToFromScreening(to, from)
	} else {
		return fmt.Errorf("정의되지 않은 방향입니다. 상세: to: [%s] from: [%s] 으로는 불가능", to, from)
	}
}

// from -> to 검증의 경우
// 희망전형(from) 원서의 적용전형은 모두 not null 을 보장 해야 함.
// 희망전형(to) 원서의 적용전형은 모두 null 을 보장 해야 함.
func validateToFromScreening(to jobs.Screening, from jobs.Screening) error {
	isNotNull := repository.IsAppliedScreeningAllNotNullBy(from)
	if isNotNull == false {
		return internal.ExtractExpectedActualIsDiffError(fmt.Sprintf("희망전형의 [%s] 적용전형이 모두 not-null 상태", from))
	}

	isNull := repository.IsAppliedScreeningAllNullBy(to)
	if isNull == false {
		return internal.ExtractExpectedActualIsDiffError(fmt.Sprintf("희망전형의 [%s] 적용전형이 모두 null 상태", to))
	}
	return nil
}

// 모든 적용전형은 null 을 보장 해야 함.
func validateAllAppliedScreening() error {
	isAllNull := repository.IsAppliedScreeningAllNull()
	if isAllNull == false {
		return internal.ExtractExpectedActualIsDiffError("적용전형은 모두 null인 상태")
	}

	return nil
}

// 특별전형으로 지원한 모든 지원자를 특별전형으로 적용한다.
func applyScreeningSpecialAll() {
	repository.SaveAppliedScreening(jobs.SpecialScreening, jobs.SpecialScreening, jobs.SpecialSuccessfulApplicantOf1E)
}

// 특별전형으로 지원한 상위 8명의 지원자를 특별전형으로 적용한다.
// 남은 X 명의 지원자는 일반 전형으로 적용하여 평가한다.
func applyScreeningSpecialPartial() {
	repository.SaveAppliedScreening(jobs.SpecialScreening, jobs.SpecialScreening, jobs.SpecialSuccessfulApplicantOf1E)
	repository.SaveAppliedScreening(jobs.SpecialScreening, jobs.GeneralScreening, jobs.JustAll)
}

func (s *CalculationOfSuccessfulApplicantByScreeningProcessStep) Processor(context *jobs.BatchContext) {
	log.Println("hello reader 2")
}
