package my_job

import (
	"gorm.io/gorm"
	"log"
	e "themoment-team/go-hellogsm/error"
	"themoment-team/go-hellogsm/jobs"
	"themoment-team/go-hellogsm/repository"
)

// AssignExaminationNumberStep 수험번호를 할당하는 Step이다.
type AssignExaminationNumberStep struct {
}

func (s *AssignExaminationNumberStep) Processor(batchContext *jobs.BatchContext, db *gorm.DB) error {
	log.Println("수험번호 할당을 시작합니다.")

	// 수험번호 할당 전 검증
	err := validateBeforeExaminationNumberAssignment(db)
	if err != nil {
		return err
	}

	// 할당 통계 로그
	totalFirstPass, _, rooms := getExaminationNumberStatistics(db)
	logExaminationNumberAssignmentStart(totalFirstPass, rooms)

	// 수험번호 할당 실행
	err = repository.AssignExaminationNumbers(db)
	if err != nil {
		log.Printf("수험번호 할당 중 오류 발생: %v", err)
		return err
	}

	// 수험번호 할당 후 검증
	err = validateAfterExaminationNumberAssignment(db)
	if err != nil {
		return err
	}

	// 최종 결과 로그
	_, assignedCount, finalRooms := getExaminationNumberStatistics(db)
	logExaminationNumberAssignmentComplete(assignedCount, finalRooms)

	// 할당 결과 샘플 로그
	logExaminationNumberSamples(db)

	log.Println("수험번호 할당이 완료되었습니다.")
	return nil
}

// 수험번호 할당 전 검증
func validateBeforeExaminationNumberAssignment(db *gorm.DB) error {
	log.Println("수험번호 할당 전 검증을 시작합니다.")

	// 1차 합격자 존재 여부 확인
	firstPassCount := repository.CountFirstPassApplicants(db)
	if firstPassCount == 0 {
		return e.WrapExpectedActualIsDiffError("1차 합격자가 존재하지 않아 수험번호 할당이 불가능합니다")
	}
	log.Printf("1차 합격자 [%d]명을 확인했습니다.", firstPassCount)

	// 기존 수험번호 할당 여부 확인 - 중복 실행 방지
	existingCount := repository.CountExistingExaminationNumbers(db)
	if existingCount > 0 {
		log.Printf("ERROR: 이미 [%d]개의 수험번호가 할당되어 있습니다.", existingCount)
		log.Printf("수험번호 할당 작업이 이미 완료된 상태입니다. 중복 실행을 방지하기 위해 작업을 중단합니다.")
		return e.WrapRollbackNeededError(e.WrapExpectedActualIsDiffError("수험번호가 이미 할당되어 있어 중복 실행을 방지합니다"))
	}
	log.Println("기존 수험번호 할당이 없음을 확인했습니다.")

	log.Println("수험번호 할당 전 검증이 완료되었습니다.")
	return nil
}

// 수험번호 할당 후 검증
func validateAfterExaminationNumberAssignment(db *gorm.DB) error {
	log.Println("수험번호 할당 후 검증을 시작합니다.")

	// 모든 1차 합격자에게 수험번호가 할당되었는지 확인
	unassignedCount := repository.CountFirstPassWithoutExaminationNumber(db)
	if unassignedCount > 0 {
		log.Printf("ERROR: [%d]명의 1차 합격자에게 수험번호가 할당되지 않았습니다.", unassignedCount)
		return e.WrapRollbackNeededError(e.WrapExpectedActualIsDiffError("일부 1차 합격자에게 수험번호가 할당되지 않았습니다"))
	}

	// 수험번호 형식 검증
	err := repository.ValidateExaminationNumberFormat(db)
	if err != nil {
		log.Printf("ERROR: 수험번호 형식 검증 실패: %v", err)
		return e.WrapRollbackNeededError(err)
	}

	// 수험번호 중복 검증
	err = repository.ValidateExaminationNumberUniqueness(db)
	if err != nil {
		log.Printf("ERROR: 수험번호 중복 검증 실패: %v", err)
		return e.WrapRollbackNeededError(err)
	}

	log.Println("수험번호 할당 후 검증이 완료되었습니다.")
	return nil
}

// 수험번호 할당 통계를 반환한다
func getExaminationNumberStatistics(db *gorm.DB) (totalFirstPass int, assignedExamNumber int, rooms int) {
	totalFirstPass = repository.CountFirstPassApplicants(db)
	assignedExamNumber = repository.CountExistingExaminationNumbers(db)

	// 필요한 고사실 수 (16명당 1개 고사실)
	if totalFirstPass > 0 {
		rooms = (totalFirstPass + 15) / 16 // 올림 계산
	}

	return totalFirstPass, assignedExamNumber, rooms
}

// 수험번호 할당 시작 로그
func logExaminationNumberAssignmentStart(totalFirstPass int, rooms int) {
	log.Printf("1차 합격자 [%d]명에 대해 수험번호를 할당합니다. (필요 고사실: [%d]개)", totalFirstPass, rooms)
}

// 수험번호 할당 완료 로그
func logExaminationNumberAssignmentComplete(assignedCount int, rooms int) {
	log.Printf("수험번호 할당이 완료되었습니다. 할당된 수험번호: [%d]개, 사용된 고사실: [%d]개", assignedCount, rooms)
}

// 할당된 수험번호 샘플 로그
func logExaminationNumberSamples(db *gorm.DB) {
	samples := repository.GetExaminationNumberSamples(db, 10)

	log.Println("=== 할당된 수험번호 샘플 (첫 10명) ===")
	for _, sample := range samples {
		log.Printf("이름: [%s], 전형: [%s], 수험번호: [%s]",
			sample.Name, sample.AppliedScreening, sample.ExaminationNumber)
	}
}
