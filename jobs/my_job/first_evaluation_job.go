package my_job

import (
	"fmt"
	"gorm.io/gorm"
	"log"
	e "themoment-team/go-hellogsm/error"
	"themoment-team/go-hellogsm/internal"
	"themoment-team/go-hellogsm/jobs"
	"themoment-team/go-hellogsm/repository"
	"themoment-team/go-hellogsm/types"
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

func (s *DecideAppliedScreeningStep) Processor(batchContext *jobs.BatchContext, db *gorm.DB) error {
	// 시작 데이터 검증
	err := canNextEvaluation(types.ExtraAdmissionScreening, types.ExtraAdmissionScreening)
	if err != nil {
		return err
	}

	// 정원 외 특별전형 평가
	// 특례 대상
	extraAdCount := repository.CountOneseoByAppliedScreening(string(types.ExtraAdmissionScreening))
	logAppliedScreeningResult(types.ExtraAdmissionScreening, types.ExtraAdmissionSuccessfulApplicantOf1E, extraAdCount)
	applyExtraAdScreening(db)
	// 국가 보훈 대상
	extraVeCount := repository.CountOneseoByAppliedScreening(string(types.ExtraVeteransScreening))
	logAppliedScreeningResult(types.ExtraVeteransScreening, types.ExtraVeteransSuccessfulApplicantOf1E, extraVeCount)
	applyExtraVeScreening(db)

	// 특별전형 평가 전 데이터 검증
	err = canNextEvaluation(types.SpecialScreening, types.ExtraVeteransScreening)
	if err != nil {
		return err
	}

	// 특별전형 평가
	specialWantedCount := repository.CountOneseoByWantedScreening(string(types.SpecialScreening))
	logAppliedScreeningResult(types.SpecialScreening, types.SpecialSuccessfulApplicantOf1E, specialWantedCount)
	applySpecialScreening(db)

	// 일반전형 평가 전 데이터 검증
	err = canNextEvaluation(types.GeneralScreening, types.SpecialScreening)
	if err != nil {
		return err
	}

	// 일반전형 평가
	generalWantedCount := repository.CountOneseoByWantedScreening(string(types.GeneralScreening))
	logAppliedScreeningResult(types.GeneralScreening, types.GeneralSuccessfulApplicantOf1E, generalWantedCount)
	applyGeneralScreening(db)

	// 평가 끝 데이터 검증
	err = canNextEvaluation(types.GeneralScreening, types.GeneralScreening)
	if err != nil {
		return err
	}

	// 합격/불합격자 구분 처리
	decideFailedApplicants(db)
	return nil
}

// 잘못된 평가 방향에 대한 검증을 진행한다.
// 데이터 꼬임을 방지하기 위한 검증을 진행한다.
func canNextEvaluation(to types.Screening, from types.Screening) error {
	switch {
	case to == types.ExtraAdmissionScreening && from == types.ExtraAdmissionScreening:
		return beforeAll()
	case to == types.SpecialScreening && from == types.ExtraVeteransScreening:
		return validateToScreening(to)
	case to == types.GeneralScreening && from == types.SpecialScreening:
		return validateToScreening(to)
	case to == types.GeneralScreening && from == types.GeneralScreening:
		return afterAll()
	default:
		return fmt.Errorf("정의되지 않은 방향입니다. 상세: to: [%s] from: [%s] 으로는 불가능", to, from)
	}
}

// from -> to 검증의 경우,
// 희망전형(to) 원서의 적용전형은 모두 null 을 보장 해야 함.
func validateToScreening(to types.Screening) error {
	isNull := repository.IsAppliedScreeningAllNullBy(string(to))
	if isNull == false {
		return e.WrapExpectedActualIsDiffError(fmt.Sprintf("희망전형의 [%s] 적용전형이 모두 null 상태", to))
	}
	return nil
}

// 1차 평가 시작 검증의 경우,
// 모든 적용전형은 null 을 보장 해야 함.
func beforeAll() error {
	isAllNull := repository.IsAppliedScreeningAllNull()
	if isAllNull == false {
		return e.WrapExpectedActualIsDiffError("적용전형은 모두 null인 상태")
	}

	return nil
}

// 1차 평가 끝 검증
func afterAll() error {
	log.Println("1차 applied_screening 결정 끝")
	return nil
}

// 정원 외 특별전형 / 특례 대상 적용전형 처리.
func applyExtraAdScreening(db *gorm.DB) {
	repository.SaveAppliedScreening(
		db,
		convertScreeningToStrArr([]types.Screening{types.ExtraAdmissionScreening}),
		string(types.ExtraAdmissionScreening),
		types.ExtraAdmissionSuccessfulApplicantOf1E,
	)
}

// 정원 외 특별전형 / 국가 보훈 대상 적용전형 처리.
func applyExtraVeScreening(db *gorm.DB) {
	repository.SaveAppliedScreening(
		db,
		convertScreeningToStrArr([]types.Screening{types.ExtraVeteransScreening}),
		string(types.ExtraVeteransScreening),
		types.ExtraVeteransSuccessfulApplicantOf1E,
	)
}

// 특별전형 대상 적용전형 처리.
func applySpecialScreening(db *gorm.DB) {
	repository.SaveAppliedScreening(
		db,
		convertScreeningToStrArr([]types.Screening{types.ExtraAdmissionScreening, types.ExtraVeteransScreening, types.SpecialScreening}),
		string(types.SpecialScreening),
		types.SpecialSuccessfulApplicantOf1E,
	)
}

// 일반전형 대상 적용전형 처리.
func applyGeneralScreening(db *gorm.DB) {
	repository.SaveAppliedScreening(
		db,
		convertScreeningToStrArr([]types.Screening{types.ExtraAdmissionScreening, types.ExtraVeteransScreening, types.SpecialScreening, types.GeneralScreening}),
		string(types.GeneralScreening),
		types.GeneralSuccessfulApplicantOf1E,
	)
}

// 불합격자 처리.
func decideFailedApplicants(db *gorm.DB) {
	repository.SaveFirstTestPassYn()
}

func logAppliedScreeningResult(wantedScreening types.Screening, success1E int, applicantCount int) {
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

func convertScreeningToStrArr(jobScreening []types.Screening) []string {
	strScreening := make([]string, len(jobScreening))
	for i, screening := range jobScreening {
		strScreening[i] = string(screening)
	}
	return strScreening
}
