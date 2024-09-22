package my_job

import (
	"fmt"
	"log"
	"themoment-team/go-hellogsm/internal"
	"themoment-team/go-hellogsm/jobs"
	"themoment-team/go-hellogsm/repository"

	"gorm.io/gorm"
)

// DecideAppliedScreeningStep 적용 전형 평가를 하는 Step 이다.
type DecideAppliedScreeningStep struct {
}

// job 에 필요한 step 들을 반환한다.
func getSteps() []jobs.Step {
	return []jobs.Step{&DecideAppliedScreeningStep{}}
}

// BuildFirstEvaluationJob 1차 평가 배치 Job을 생성한다.
func BuildFirstEvaluationJob(properties internal.ApplicationProperties) *jobs.SimpleJob {
	return jobs.NewSimpleJob(internal.FirstEvaluationJob, getSteps(), nil)
}

func (s *DecideAppliedScreeningStep) Processor(context *jobs.BatchContext, tx *gorm.DB) error {
	// 시작 데이터 검증
	err := canNextEvaluation(jobs.ExtraAdmissionScreening, jobs.ExtraAdmissionScreening, tx)
	if err != nil {
		return err
	}

	// 정원 외 특별전형 평가
	// 특례 대상
	extraAdCount, err := repository.CountOneseoByAppliedScreening(string(jobs.ExtraAdmissionScreening), tx)
	if err != nil {
		return err
	}
	logAppliedScreeningResult(jobs.ExtraAdmissionScreening, jobs.ExtraAdmissionSuccessfulApplicantOf1E, extraAdCount)
	err = applyExtraAdScreening(tx)
	if err != nil {
		return err
	}

	// 국가 보훈 대상
	extraVeCount, err := repository.CountOneseoByAppliedScreening(string(jobs.ExtraVeteransScreening), tx)
	if err != nil {
		return err
	}
	logAppliedScreeningResult(jobs.ExtraVeteransScreening, jobs.ExtraVeteransSuccessfulApplicantOf1E, extraVeCount)
	err = applyExtraVeScreening(tx)
	if err != nil {
		return err
	}

	// 특별전형 평가 전 데이터 검증
	err = canNextEvaluation(jobs.SpecialScreening, jobs.ExtraVeteransScreening, tx)
	if err != nil {
		return err
	}

	// 특별전형 평가
	specialWantedCount, err := repository.CountOneseoByWantedScreening(string(jobs.SpecialScreening), tx)
	if err != nil {
		return err
	}
	logAppliedScreeningResult(jobs.SpecialScreening, jobs.SpecialSuccessfulApplicantOf1E, specialWantedCount)
	err = applySpecialScreening(tx)
	if err != nil {
		return err
	}

	// 일반전형 평가 전 데이터 검증
	err = canNextEvaluation(jobs.GeneralScreening, jobs.SpecialScreening, tx)
	if err != nil {
		return err
	}

	// 일반전형 평가
	generalWantedCount, err := repository.CountOneseoByWantedScreening(string(jobs.GeneralScreening), tx)
	if err != nil {
		return err
	}
	logAppliedScreeningResult(jobs.GeneralScreening, jobs.GeneralSuccessfulApplicantOf1E, generalWantedCount)
	err = applyGeneralScreening(tx)
	if err != nil {
		return err
	}

	// 평가 끝 데이터 검증
	err = canNextEvaluation(jobs.GeneralScreening, jobs.GeneralScreening, tx)
	if err != nil {
		return err
	}

	// 합격/불합격자 구분 처리
	err = decideFailedApplicants(tx)
	if err != nil {
		return err
	}

	return nil
}

// 잘못된 평가 방향에 대한 검증을 진행한다.
// 데이터 꼬임을 방지하기 위한 검증을 진행한다.
func canNextEvaluation(to jobs.Screening, from jobs.Screening, tx *gorm.DB) error {
	switch {
	case to == jobs.ExtraAdmissionScreening && from == jobs.ExtraAdmissionScreening:
		return beforeAll(tx)
	case to == jobs.SpecialScreening && from == jobs.ExtraVeteransScreening:
		return validateToScreening(to, tx)
	case to == jobs.GeneralScreening && from == jobs.SpecialScreening:
		return validateToScreening(to, tx)
	case to == jobs.GeneralScreening && from == jobs.GeneralScreening:
		return afterAll()
	default:
		return fmt.Errorf("정의되지 않은 방향입니다. 상세: to: [%s] from: [%s] 으로는 불가능", to, from)
	}
}

// from -> to 검증의 경우,
// 희망전형(to) 원서의 적용전형은 모두 null 을 보장 해야 함.
func validateToScreening(to jobs.Screening, tx *gorm.DB) error {
	isNull, err := repository.IsAppliedScreeningAllNullBy(string(to), tx)
	if err != nil {
		return err
	}
	if isNull == false {
		return internal.ExtractExpectedActualIsDiffError(fmt.Sprintf("희망전형의 [%s] 적용전형이 모두 null 상태", to))
	}
	return nil
}

// 1차 평가 시작 검증의 경우,
// 모든 적용전형은 null 을 보장 해야 함.
func beforeAll(tx *gorm.DB) error {
	isAllNull, err := repository.IsAppliedScreeningAllNull(tx)
	if err != nil {
		return err
	}
	if isAllNull == false {
		return internal.ExtractExpectedActualIsDiffError("적용전형은 모두 null인 상태")
	}

	return nil
}

// 1차 평가 끝 검증
func afterAll() error {
	log.Println("1차 applied_screening 결정 끝")
	return nil
}

// 정원 외 특별전형 / 특례 대상 적용전형 처리.
func applyExtraAdScreening(tx *gorm.DB) error {
	return repository.SaveAppliedScreening(
		convertScreeningToStrArr([]jobs.Screening{jobs.ExtraAdmissionScreening}),
		string(jobs.ExtraAdmissionScreening),
		jobs.ExtraAdmissionSuccessfulApplicantOf1E,
		tx,
	)
}

// 정원 외 특별전형 / 국가 보훈 대상 적용전형 처리.
func applyExtraVeScreening(tx *gorm.DB) error {
	return repository.SaveAppliedScreening(
		convertScreeningToStrArr([]jobs.Screening{jobs.ExtraVeteransScreening}),
		string(jobs.ExtraVeteransScreening),
		jobs.ExtraVeteransSuccessfulApplicantOf1E,
		tx,
	)
}

// 특별전형 대상 적용전형 처리.
func applySpecialScreening(tx *gorm.DB) error {
	return repository.SaveAppliedScreening(
		convertScreeningToStrArr([]jobs.Screening{jobs.ExtraAdmissionScreening, jobs.ExtraVeteransScreening, jobs.SpecialScreening}),
		string(jobs.SpecialScreening),
		jobs.SpecialSuccessfulApplicantOf1E,
		tx,
	)
}

// 일반전형 대상 적용전형 처리.
func applyGeneralScreening(tx *gorm.DB) error {
	return repository.SaveAppliedScreening(
		convertScreeningToStrArr([]jobs.Screening{jobs.ExtraAdmissionScreening, jobs.ExtraVeteransScreening, jobs.SpecialScreening, jobs.GeneralScreening}),
		string(jobs.GeneralScreening),
		jobs.GeneralSuccessfulApplicantOf1E,
		tx,
	)
}

// 불합격자 처리.
func decideFailedApplicants(tx *gorm.DB) error {
	return repository.SaveFirstTestPassYn(tx)
}

func logAppliedScreeningResult(wantedScreening jobs.Screening, success1E int, applicantCount int) {
	log.Printf("[%s]전형 희망 지원자 수 [%d]명 중 [%d]명을 [%s]전형으로 적용합니다.",
		wantedScreening,
		applicantCount,
		success1E,
		wantedScreening)

	// 최대 합격자 수를 초과할 경우 다음 평가 프로세스가 적용됨을 알린다.
	if applicantCount > success1E {
		log.Printf("하위 [%d]명은 다음 평가 프로세스가 적용됩니다.",
			applicantCount-success1E)
	}
}

func convertScreeningToStrArr(jobScreening []jobs.Screening) []string {
	strScreening := make([]string, len(jobScreening))
	for i, screening := range jobScreening {
		strScreening[i] = string(screening)
	}
	return strScreening
}
