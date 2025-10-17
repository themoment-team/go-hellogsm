package my_job

import (
	"log"
	e "themoment-team/go-hellogsm/error"
	"themoment-team/go-hellogsm/jobs"
	"themoment-team/go-hellogsm/repository"

	"gorm.io/gorm"
)

type AssignExaminationNumberStep struct {
}

func (s *AssignExaminationNumberStep) Processor(batchContext *jobs.BatchContext, db *gorm.DB) error {
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
		return e.WrapExpectedActualIsDiffError("수험번호가 이미 할당되어 있어 중복 실행을 방지합니다")
	}
	return nil
}

// 수험번호 할당 후 검증
func validateAfterExaminationNumberAssignment(db *gorm.DB) error {
	// 모든 1차 합격자에게 수험번호가 할당되었는지 확인
	unassignedCount := repository.CountFirstPassWithoutExaminationNumber(db)
	overflow := repository.GetOverflowApplicants(db)
	if unassignedCount > 0 {
		if unassignedCount == len(overflow) {
			log.Printf("좌석 수 초과로 수험번호 미부여된 인원: [%d]명", len(overflow))
			for i, o := range overflow {
				log.Printf("초과자 #%d: 이름=[%s], 적용전형=[%s], 제출코드=[%s]", i+1, o.Name, o.AppliedScreening, o.OneseoSubmitCode)
			}
		} else {
			log.Printf("ERROR: 수험번호 미할당자 수 [%d]명, 초과자 수 [%d]명 불일치", unassignedCount, len(overflow))
			return e.WrapExpectedActualIsDiffError("일부 1차 합격자에게 수험번호가 할당되지 않았습니다")
		}
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

// 수험번호 할당 통계 반환
func getExaminationNumberStatistics(db *gorm.DB) (totalFirstPass int, assignedExamNumber int, rooms int) {
	totalFirstPass = repository.CountFirstPassApplicants(db)
	assignedExamNumber = repository.CountExistingExaminationNumbers(db)

	if totalFirstPass > 0 {
		// 18,18,18,18,12,11 기준으로 필요한 고사실 수 산정, 추후 고사실 상황 변동 시 수정 필요
		caps := []int{18, 18, 18, 18, 12, 11}
		remain := totalFirstPass
		rooms = 0
		for _, c := range caps {
			if remain <= 0 {
				break
			}
			rooms++
			remain -= c
		}
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
