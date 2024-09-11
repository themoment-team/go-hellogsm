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

func (s *SecondEvaluationAbsenteeExclusionStep) Processor(context *jobs.BatchContext) {
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
}
