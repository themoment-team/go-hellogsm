package my_job

import (
	"log"
	"themoment-team/go-hellogsm/internal"
	"themoment-team/go-hellogsm/jobs"
)

// ApplicationScreeningEvaluationStep 적용 전형 평가를 하는 Step 이다.
type ApplicationScreeningEvaluationStep struct {
}

// CalculationOfSuccessfulApplicantByScreeningProcessStep 전형별 합격자를 산출하는 Step 이다.
type CalculationOfSuccessfulApplicantByScreeningProcessStep struct {
}

type FirstEvaluationJobListener struct {
}

// job 에 필요한 step 들을 반환한다.
func getSteps() []jobs.Step {
	return []jobs.Step{&ApplicationScreeningEvaluationStep{}, &CalculationOfSuccessfulApplicantByScreeningProcessStep{}}
}

// BuildFirstEvaluationJob 1차 평가 배치 Job을 생성한다.
func BuildFirstEvaluationJob(properties internal.ApplicationProperties) *jobs.SimpleJob {
	return jobs.NewSimpleJob(internal.FirstEvaluationJob, getSteps(), nil)
}

func (s *ApplicationScreeningEvaluationStep) Reader() {
	log.Println("hello reader 1")
}

func (s *ApplicationScreeningEvaluationStep) Writer() {
	log.Println("hello writer 1")
}

func (s *CalculationOfSuccessfulApplicantByScreeningProcessStep) Reader() {
	log.Println("hello reader 2")
}

func (s *CalculationOfSuccessfulApplicantByScreeningProcessStep) Writer() {
	log.Println("hello writer 2")
}
