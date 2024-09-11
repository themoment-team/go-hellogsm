package my_job

import (
	"themoment-team/go-hellogsm/internal"
	"themoment-team/go-hellogsm/jobs"
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
}

func (s *TotalEvaluationTopScoringApplicantsSelectionByScreeningStep) Processor(context *jobs.BatchContext) {
}
