package repository

import (
	"gorm.io/gorm"
	"themoment-team/go-hellogsm/configs"
	e "themoment-team/go-hellogsm/error"
)

type ExaminationNumberSample struct {
	Name              string
	AppliedScreening  string
	ExaminationNumber string
}

func AssignExaminationNumbers(db *gorm.DB) error {
	query := `
WITH FirstPassApplicants AS (
    SELECT 
        o.oneseo_id,
        m.name,
        o.applied_screening,
        ROW_NUMBER() OVER (ORDER BY m.name ASC, o.oneseo_id ASC) AS row_num
    FROM tb_oneseo o
    JOIN tb_member m ON o.member_id = m.member_id
    WHERE o.applied_screening IS NOT NULL 
      AND o.real_oneseo_arrived_yn = 'YES'
)
UPDATE tb_oneseo AS o
JOIN FirstPassApplicants AS fpa ON o.oneseo_id = fpa.oneseo_id
SET o.examination_number = CONCAT(
    LPAD(FLOOR((fpa.row_num - 1) / 16) + 1, 2, '0'),
    LPAD(((fpa.row_num - 1) % 16) + 1, 2, '0')
);`

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

func HasFirstPassApplicantsWithDB(db *gorm.DB) bool {
	return CountFirstPassApplicants(db) > 0
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
		rooms = (totalFirstPass + 15) / 16
	}

	return totalFirstPass, assignedExamNumber, rooms
}
