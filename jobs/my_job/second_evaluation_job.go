package my_job

import (
	"log"
	"themoment-team/go-hellogsm/internal"
	"themoment-team/go-hellogsm/jobs"
	"themoment-team/go-hellogsm/repository"
)

// SecondEvaluationAbsenteeExclusionStep 2차 전형(직무적성검사 or 심층면접) 미응시자를 탈락 처리하는 Step 이다.
type SecondEvaluationAbsenteeExclusionStep struct {
}

// TotalEvaluationTopScoringApplicantsSelectionByScreeningStep 모든 평가(1차 + 2차)를 기반으로 상위 성적의 지원자들을 전형별로 선발하는 Step 이다.
type TotalEvaluationTopScoringApplicantsSelectionByScreeningStep struct {
}

// job 에 필요한 step 들을 반환한다.
func getSecondEvaluationSteps() []jobs.Step {
	return []jobs.Step{&SecondEvaluationAbsenteeExclusionStep{}, &TotalEvaluationTopScoringApplicantsSelectionByScreeningStep{}}
}

// BuildSecondEvaluationJob 2차 평가 배치 Job을 생성한다.
func BuildSecondEvaluationJob(properties internal.ApplicationProperties) *jobs.SimpleJob {
	return jobs.NewSimpleJob(internal.SecondEvaluationJob, getSecondEvaluationSteps(), nil)
}

func (s *SecondEvaluationAbsenteeExclusionStep) Processor(context *jobs.BatchContext) { // 하나 트랜잭션으로 묶
	// 처리 전 데이터 검증
	err := PreCheckAbsenteeExclusion()
	if err != nil {
		panic(err.Error())
	}

	// 2차 전형 미응시자 탈락 처리
	log.Printf("2차 전형(직무적성검사 or 심층면접) 미응시자를 탈락 처리합니다")
	repository.UpdateSecondTestPassStatusForAbsentees()

	// 처리 후 데이터 검증
	err = PostCheckAbsenteeExclusion()
	if err != nil {
		panic(err.Error())
	}
}

// 1차 평가를 마친 data는 applied_screening이 존재햐야 한다.
func PreCheckAbsenteeExclusion() error {
	isAllNotNull := repository.IsAllFirstPassUserHaveAppliedScreening()
	if isAllNotNull == false {
		return internal.ExtractExpectedActualIsDiffError("1차 평가를 마친 data는 전부 applied_screening이 존재햐야 한다.")
	}

	return nil
}

// 2차 전형(직무적성검사 or 심층면접) 미응시자는 탈락처리 되어있다.
func PostCheckAbsenteeExclusion() error {
	isAllAbsenteeFall := repository.IsAllAbsenteeFall()
	if isAllAbsenteeFall == false {
		return internal.ExtractExpectedActualIsDiffError("2차 전형(직무적성검사 or 심층면접) 미응시자는 전부 탈락처리 되어있다.")
	}

	return nil
}

func (s *TotalEvaluationTopScoringApplicantsSelectionByScreeningStep) Processor(context *jobs.BatchContext) {
	// 특별전형으로 2차에 응시한 원서 조회
	log.Printf("특별전형으로 2차 전형에 응시한 지원자를 조회합니다.")
	specialOneseoIds := repository.QuerySpecialOneseoIds()

	var remainingSpecialOneseos int

	// 특별전형 limit 분기처리 ( 그냥 하나로 합쳐도 될 듯? )
	if len(specialOneseoIds) <= jobs.SpecialSuccessfulApplicantOf2E {
		// 특별전형 인원이 limit 이하일때 전부 2차전형 합격처리
		log.Printf("특별전형으로 2차 전형을 진행한 인원이 %d명 이하일 때 전부 합격처리합니다.", jobs.SpecialSuccessfulApplicantOf2E)
		repository.UpdateSecondTestPassYnForSpecialPass(specialOneseoIds)
		remainingSpecialOneseos = jobs.SpecialSuccessfulApplicantOf2E - len(specialOneseoIds)

	} else {
		// 특별전형 인원이 limit 초과일때 상위 limit명은 합격처리 후 하위 n-limit명은 일반전형으로 변경
		log.Printf("특별전형으로 2차 전형을 진행한 인원이 %d명 초과일 때", jobs.SpecialSuccessfulApplicantOf2E)

		log.Printf("상위 %d명은 합격처리", jobs.SpecialSuccessfulApplicantOf2E)
		passSpecialOneseos := specialOneseoIds[:jobs.SpecialSuccessfulApplicantOf2E] // 인덱스 0~limit 까지 처리
		repository.UpdateSecondTestPassYnForSpecialPass(passSpecialOneseos)          // 상위 limit명 합격처리

		log.Printf("하위 %d명은 일반전형으로 변경합니다.", len(specialOneseoIds)-jobs.SpecialSuccessfulApplicantOf2E)
		fallSpecialOneseos := specialOneseoIds[jobs.SpecialSuccessfulApplicantOf2E:] // 인덱스 limit~n 까지 처리
		repository.UpdateAppliedScreeingForSpecialFall(fallSpecialOneseos)           // 하위 n-limit명 일반전형으로 변경
	}

	// 일반전형으로 2차에 응시한 지원자 상위 limit명 합격처리
	log.Printf("일반전형으로 2차 전형을 진행한 인원 중 성적 상위 %d명을 합격처리합니다.", jobs.GeneralSuccessfulApplicantOf2E+remainingSpecialOneseos)
	requiredGeneralOneseos := jobs.GeneralSuccessfulApplicantOf2E + remainingSpecialOneseos
	repository.UpdateSecondTestPassYnForGeneral(requiredGeneralOneseos)
}
