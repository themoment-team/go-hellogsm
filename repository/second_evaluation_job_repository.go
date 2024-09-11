package repository

import (
	"themoment-team/go-hellogsm/configs"
)

func UpdateSecondTestPassStatusForAbsentees() {
	configs.MyDB.Raw(`
		UPDATE tb_entrance_test_result
		SET second_test_pass_yn = 'NO'
		WHERE first_test_pass_yn = 'YES' 
		  AND (aptitude_evaluation_score IS NULL OR interview_score IS NULL)
	`)
}

func IsAllFirstPassUserHaveAppliedScreening() bool {
	var result int
	configs.MyDB.Raw(`
		SELECT COUNT(*) 
		FROM tb_oneseo 
		WHERE applied_screening IS NULL 
		  AND oneseo_id IN (
			SELECT oneseo_id 
			FROM tb_entrance_test_result 
			WHERE first_test_pass_yn = 'YES'
		)
	`).Scan(&result)
	return result < 1
}

func IsAllAbsenteeFall() bool {
	var absenteeCount int
	var fallCount int

	// 미응시자 count query
	configs.MyDB.Raw(`
		SELECT COUNT(*) 
		FROM tb_entrance_test_result 
		WHERE first_test_pass_yn = 'YES' 
		  AND (aptitude_evaluation_score IS NULL OR interview_score IS NULL)
	`).Scan(&absenteeCount)

	// 2차 전형 탈락자 count query
	configs.MyDB.Raw(`
		SELECT COUNT(*) 
		FROM tb_entrance_test_result 
		WHERE second_test_pass_yn = 'NO'
	`).Scan(&fallCount)

	return absenteeCount == fallCount
}
