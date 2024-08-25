package jobs

import "themoment-team/go-hellogsm/internal"

type Job interface {
	Execute(properties internal.ApplicationProperties) error
}
