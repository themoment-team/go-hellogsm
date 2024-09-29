package jobs

import (
	"gorm.io/gorm"
)

// Step 은 간단히 Processor 만을 갖는다.
// 그냥 Processor 에 해당하는 함수를 실행시키면 되지만, Step 을 사용해서 객체적으로 구분한다.
// 추후 Step 에는 Reader, Writer 등의 개념을 부여할 수 있다.
type Step interface {
	Processor(batchContext *BatchContext, db *gorm.DB) error
}
