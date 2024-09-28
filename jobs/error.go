package jobs

import (
	"fmt"
	"log"
)

type FatalError error
type ExpectedActualIsDiffError error
type RollbackNeededError error

func WrapExpectedActualIsDiffError(expected string) RollbackNeededError {
	return WrapRollbackNeededError(fmt.Sprintf("기대 결과와 실제 결과가 다름. 기대 결과: [%s]", expected))
}

func WrapRollbackNeededError(msg string) RollbackNeededError {
	return fmt.Errorf("%s", msg)
}

func WrapFatalErr(err error) {
	log.Println("An fatal error occurred. Exiting...")
	panic(err)
}
