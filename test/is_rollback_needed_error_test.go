package test

import (
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	myError "themoment-team/go-hellogsm/error"
)

func TestIsRollbackNeededError(t *testing.T) {
	rollbackNeededError := myError.WrapRollbackNeededError(fmt.Errorf("helloworld"))

	actual := errTypeCheck(rollbackNeededError)
	assert.Equal(t, "same", actual)
}

func errTypeCheck(err error) string {
	var rollbackNeededError *myError.RollbackNeededError

	if errors.As(err, &rollbackNeededError) {
		return fmt.Sprintf("same")
	} else {
		return fmt.Sprintf("different")
	}
}
