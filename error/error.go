package error

import (
	"fmt"
	"log"
)

type FatalError error
type ExpectedActualIsDiffError error
type RollbackNeededError struct {
	err error
}

func WrapExpectedActualIsDiffError(expected string) *RollbackNeededError {
	return WrapRollbackNeededError(fmt.Errorf("기대 결과와 실제 결과가 다름. 기대 결과: [%s]", expected))
}

func WrapRollbackNeededError(err error) *RollbackNeededError {
	return &RollbackNeededError{err: err}
}

func WrapFatalErr(err error) {
	log.Println("An fatal error occurred. Exiting...")
	panic(err)
}

func (e RollbackNeededError) Error() string {
	return e.err.Error()
}
