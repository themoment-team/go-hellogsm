package internal

import (
	"fmt"
)

func ExtractExpectedActualIsDiffError(expected string) error {
	return fmt.Errorf("기대 결과와 실제 결과가 다름. 기대 결과: [%s]", expected)
}
