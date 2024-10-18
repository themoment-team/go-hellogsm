package error

import (
	"fmt"
)

type RollbackNeededError struct {
	err error
}

func WrapExpectedActualIsDiffError(expected string) error {
	return WrapRollbackNeededError(fmt.Errorf("기대 결과와 실제 결과가 다름. 기대 결과: [%s]", expected))
}

func WrapRollbackNeededError(err error) *RollbackNeededError {
	return &RollbackNeededError{err: err}
}

func (e RollbackNeededError) Error() string {
	return e.err.Error()
}
