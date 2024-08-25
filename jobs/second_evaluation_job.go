package jobs

import "themoment-team/go-hellogsm/internal"

type SecondEvaluationJob struct{}

func (f *SecondEvaluationJob) Execute(internal.ApplicationProperties) error {
	return nil
}
