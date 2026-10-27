package repository

import (
	"fmt"
	"strings"

	"themoment-team/go-hellogsm/configs"
	e "themoment-team/go-hellogsm/error"
	"themoment-team/go-hellogsm/types"

	"gorm.io/gorm"
)

type ExaminationNumberSample struct {
	Name              string
	AppliedScreening  string
	ExaminationNumber string
}

type OverflowApplicant struct {
	Name             string
	AppliedScreening string
	OneseoSubmitCode string
}

// buildRoomCaseClauses 고사실 정원 목록으로부터 수험번호 할당용 CASE 절 본문을 생성한다.
// 고사실 정원이 바뀌어도 이 함수가 유일한 기준이 되므로, SQL의 경계값을 손으로 맞출 필요가 없다.
func buildRoomCaseClauses(capacities []int) (roomCase string, seatCase string, total int) {
	var roomSb, seatSb strings.Builder
	cumulative := 0

	for i, capacity := range capacities {
		roomNumber := i + 1
		previousCumulative := cumulative
		cumulative += capacity

		fmt.Fprintf(&roomSb, "WHEN fpa.row_num <= %d THEN %d\n", cumulative, roomNumber)
		if previousCumulative == 0 {
			fmt.Fprintf(&seatSb, "WHEN fpa.row_num <= %d THEN fpa.row_num\n", cumulative)
		} else {
			fmt.Fprintf(&seatSb, "WHEN fpa.row_num <= %d THEN fpa.row_num - %d\n", cumulative, previousCumulative)
		}
	}

	return roomSb.String(), seatSb.String(), cumulative
}

// 고사장 상황 변동 시 types.ExaminationRoomCapacities만 수정하면 된다.
func AssignExaminationNumbers(db *gorm.DB) error {
	roomCase, seatCase, total := buildRoomCaseClauses(types.ExaminationRoomCapacities)
	lastRoomNumber := len(types.ExaminationRoomCapacities)

	query := fmt.Sprintf(`
WITH FirstPassApplicants AS (
    SELECT
        o.oneseo_id,
        m.name,
        o.applied_screening,
        o.oneseo_submit_code,
        ROW_NUMBER() OVER (
            ORDER BY
                SUBSTRING_INDEX(o.oneseo_submit_code, '-', 1) ASC,
                CAST(SUBSTRING_INDEX(o.oneseo_submit_code, '-', -1) AS UNSIGNED) ASC
        ) AS row_num
    FROM tb_oneseo o
    JOIN tb_member m ON o.member_id = m.member_id
    WHERE o.applied_screening IS NOT NULL
      AND o.real_oneseo_arrived_yn = 'YES'
)
UPDATE tb_oneseo AS o
JOIN FirstPassApplicants AS fpa ON o.oneseo_id = fpa.oneseo_id
SET o.examination_number = CONCAT(
    LPAD(
        CASE
            %s            ELSE %d
        END
    , 2, '0'),
    LPAD(
        CASE
            %s            ELSE 0
        END
    , 2, '0')
)
WHERE fpa.row_num <= %d;`, roomCase, lastRoomNumber, seatCase, total)

	err := db.Exec(query).Error
	if err != nil {
		return e.WrapRollbackNeededError(err)
	}

	return nil
}

func CountFirstPassApplicants(db *gorm.DB) int {
	var count int
	db.Raw(`
		SELECT COUNT(*) 
		FROM tb_oneseo 
		WHERE applied_screening IS NOT NULL 
		  AND real_oneseo_arrived_yn = 'YES'
	`).Scan(&count)
	return count
}

func CountExistingExaminationNumbers(db *gorm.DB) int {
	var count int
	db.Raw("SELECT COUNT(*) FROM tb_oneseo WHERE examination_number IS NOT NULL").Scan(&count)
	return count
}

func CountFirstPassWithoutExaminationNumber(db *gorm.DB) int {
	var count int
	db.Raw(`
		SELECT COUNT(*) 
		FROM tb_oneseo 
		WHERE applied_screening IS NOT NULL 
		  AND real_oneseo_arrived_yn = 'YES' 
		  AND examination_number IS NULL
	`).Scan(&count)
	return count
}

func ValidateExaminationNumberFormat(db *gorm.DB) error {
	var invalidCount int
	db.Raw(`
		SELECT COUNT(*) 
		FROM tb_oneseo 
		WHERE examination_number IS NOT NULL 
		  AND examination_number NOT REGEXP '^[0-9]{4}$'
	`).Scan(&invalidCount)

	if invalidCount > 0 {
		return e.WrapExpectedActualIsDiffError("수험번호 형식이 올바르지 않은 데이터 존재")
	}
	return nil
}

func ValidateExaminationNumberUniqueness(db *gorm.DB) error {
	var duplicateCount int
	db.Raw(`
		SELECT COUNT(*) - COUNT(DISTINCT examination_number)
		FROM tb_oneseo 
		WHERE examination_number IS NOT NULL
	`).Scan(&duplicateCount)

	if duplicateCount > 0 {
		return e.WrapExpectedActualIsDiffError("중복된 수험번호 존재")
	}
	return nil
}

func GetExaminationNumberSamples(db *gorm.DB, limit int) []ExaminationNumberSample {
	var samples []ExaminationNumberSample
	db.Raw(`
		SELECT m.name, o.applied_screening, o.examination_number
		FROM tb_oneseo o
		JOIN tb_member m ON o.member_id = m.member_id
		WHERE o.examination_number IS NOT NULL
		ORDER BY o.examination_number ASC
		LIMIT ?
	`, limit).Scan(&samples)

	return samples
}

// GetOverflowApplicants 고사실 정원 합계(types.ExaminationRoomCapacities)를 초과한 지원자 목록을
// 정렬 순서대로 반환한다. AssignExaminationNumbers와 동일한 정렬 기준을 사용해야
// validateAfterExaminationNumberAssignment의 미할당자 수 대조가 성립한다.
func GetOverflowApplicants(db *gorm.DB) []OverflowApplicant {
	total := types.TotalExaminationRoomCapacity()

	var list []OverflowApplicant
	db.Raw(fmt.Sprintf(`
		WITH FirstPassApplicants AS (
		    SELECT
		        m.name,
		        o.applied_screening,
		        o.oneseo_submit_code,
		        ROW_NUMBER() OVER (
		            ORDER BY
		                SUBSTRING_INDEX(o.oneseo_submit_code, '-', 1) ASC,
		                CAST(SUBSTRING_INDEX(o.oneseo_submit_code, '-', -1) AS UNSIGNED) ASC
		        ) AS row_num
		    FROM tb_oneseo o
		    JOIN tb_member m ON o.member_id = m.member_id
		    WHERE o.applied_screening IS NOT NULL
		      AND o.real_oneseo_arrived_yn = 'YES'
		)
		SELECT name, applied_screening, oneseo_submit_code
		FROM FirstPassApplicants
		WHERE row_num > %d
		ORDER BY row_num ASC
	`, total)).Scan(&list)
	return list
}

func HasFirstPassApplicants() bool {
	var count int
	configs.MyDB.Raw(`
		SELECT COUNT(*) 
		FROM tb_oneseo 
		WHERE applied_screening IS NOT NULL 
		  AND real_oneseo_arrived_yn = 'YES'
	`).Scan(&count)
	return count > 0
}

func HasOneseoWithExaminationNumber() bool {
	var count int
	configs.MyDB.Raw(`
		SELECT COUNT(*) 
		FROM tb_oneseo 
		WHERE examination_number IS NOT NULL
	`).Scan(&count)
	return count > 0
}

func GetExaminationNumberStats() (totalFirstPass int, assignedExamNumber int, rooms int) {
	configs.MyDB.Raw(`
		SELECT COUNT(*) 
		FROM tb_oneseo 
		WHERE applied_screening IS NOT NULL 
		  AND real_oneseo_arrived_yn = 'YES'
	`).Scan(&totalFirstPass)

	configs.MyDB.Raw(`
		SELECT COUNT(*) 
		FROM tb_oneseo 
		WHERE examination_number IS NOT NULL
	`).Scan(&assignedExamNumber)

	if totalFirstPass > 0 {
		remain := totalFirstPass
		rooms = 0
		for _, c := range types.ExaminationRoomCapacities {
			if remain <= 0 {
				break
			}
			rooms++
			remain -= c
		}
		if remain > 0 {
			// 정의된 방 수를 초과하는 경우 남는 인원은 미집계
		}
	}

	return totalFirstPass, assignedExamNumber, rooms
}
