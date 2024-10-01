package repository

import (
	"themoment-team/go-hellogsm/configs"
	e "themoment-team/go-hellogsm/error"

	"gorm.io/gorm"
)

func UpdateSecondTestPassStatusForAbsentees(db *gorm.DB) error {
	query := (`
		UPDATE tb_entrance_test_result
		SET second_test_pass_yn = 'NO'
		WHERE first_test_pass_yn = 'YES' 
		  AND (aptitude_evaluation_score IS NULL OR interview_score IS NULL)
	`)

	return e.WrapRollbackNeededError(db.Exec(query).Error)
}

func IsAllFirstPassUserHaveAppliedScreening(db *gorm.DB) (bool, error) {
	query := (`
		SELECT COUNT(*) 
		FROM tb_oneseo 
		WHERE applied_screening IS NULL 
		  AND oneseo_id IN (
			SELECT oneseo_id 
			FROM tb_entrance_test_result 
			WHERE first_test_pass_yn = 'YES'
		)
	`)

	var result int
	err := e.WrapRollbackNeededError(db.Raw(query).Scan(&result).Error)
	if err != nil {
		return false, err
	}

	return result < 1, nil
}

func IsAllAbsenteeFall(db *gorm.DB) (bool, error) {
	// 미응시자 count query
	query := (`
		SELECT COUNT(*) 
		FROM tb_entrance_test_result 
		WHERE first_test_pass_yn = 'YES' 
		  AND (aptitude_evaluation_score IS NULL OR interview_score IS NULL)
	`)
	var absenteeCount int
	err := e.WrapRollbackNeededError(db.Raw(query).Scan(&absenteeCount).Error)
	if err != nil {
		return false, err
	}

	// 2차 전형 탈락자 count query
	query = (`
		SELECT COUNT(*) 
		FROM tb_entrance_test_result 
		WHERE second_test_pass_yn = 'NO'
	`)
	var fallCount int
	err = e.WrapRollbackNeededError(db.Raw(query).Scan(&fallCount).Error)
	if err != nil {
		return false, err
	}

	return absenteeCount == fallCount, nil
}

// extra admission oneseo id조회 쿼리 (성적순으로 order)
func QueryExtraAdOneseoIds() []int {
	var ids []int
	configs.MyDB.Raw(`
		SELECT o.oneseo_id 
		FROM tb_oneseo o
		JOIN tb_entrance_test_result tr ON o.oneseo_id = tr.oneseo_id
		JOIN tb_entrance_test_factors_detail td ON tr.entrance_test_factors_detail_id = td.entrance_test_factors_detail_id
		JOIN tb_member m ON o.member_id = m.member_id
		WHERE 
			o.applied_screening = 'EXTRA_ADMISSION' 
			AND tr.second_test_pass_yn IS NULL
		ORDER BY 
			(((tr.document_evaluation_score / 3) * 0.5) + (tr.aptitude_evaluation_score * 0.3) + (tr.interview_score * 0.2)) DESC, 
			td.total_subjects_score DESC,
			(td.score_3_2 + td.score_3_1) DESC,
			(td.score_2_2 + td.score_2_1) DESC, 
			td.score_2_2 DESC, 
			td.score_2_1 DESC, 
			td.total_non_subjects_score DESC, 
			m.birth ASC;
	`).Scan(&ids)

	return ids
}

// extra admission limit명 이하일때 second_test = pass
func UpdateSecondTestPassYnForExtraAdPass(passExtraAdOneseoIds []int) {

	// second_test_pass_yn = YES
	configs.MyDB.Exec(`
		UPDATE tb_entrance_test_result
		SET second_test_pass_yn = 'YES'
		WHERE oneseo_id IN ?
	`, passExtraAdOneseoIds)
}

// extra admission limit명 초과일때 하위 n명 applied_screening = SPECIAL 설정 쿼리
func UpdateAppliedScreeingForExtraAdFall(fallExtraAdOneseoIds []int) {
	configs.MyDB.Exec(`
		UPDATE tb_oneseo
		SET applied_screening = 'SPECIAL'
		WHERE oneseo_id IN ?
	`, fallExtraAdOneseoIds)
}

// extra veteran oneseo id조회 쿼리 (성적순으로 order)
func QueryExtraVeOneseoIds() []int {
	var ids []int
	configs.MyDB.Raw(`
		SELECT o.oneseo_id 
		FROM tb_oneseo o
		JOIN tb_entrance_test_result tr ON o.oneseo_id = tr.oneseo_id
		JOIN tb_entrance_test_factors_detail td ON tr.entrance_test_factors_detail_id = td.entrance_test_factors_detail_id
		JOIN tb_member m ON o.member_id = m.member_id
		WHERE 
			o.applied_screening = 'EXTRA_VETERANS' 
			AND tr.second_test_pass_yn IS NULL
		ORDER BY 
			(((tr.document_evaluation_score / 3) * 0.5) + (tr.aptitude_evaluation_score * 0.3) + (tr.interview_score * 0.2)) DESC, 
			td.total_subjects_score DESC,
			(td.score_3_2 + td.score_3_1) DESC,
			(td.score_2_2 + td.score_2_1) DESC, 
			td.score_2_2 DESC, 
			td.score_2_1 DESC, 
			td.total_non_subjects_score DESC, 
			m.birth ASC;
	`).Scan(&ids)

	return ids
}

// extra veteran limit명 이하일때 second_test = pass
func UpdateSecondTestPassYnForExtraVePass(passExtraVeOneseoIds []int) {

	// second_test_pass_yn = YES
	configs.MyDB.Exec(`
		UPDATE tb_entrance_test_result
		SET second_test_pass_yn = 'YES'
		WHERE oneseo_id IN ?
	`, passExtraVeOneseoIds)
}

// extra veteran limit명 초과일때 하위 n명 applied_screening = SPECIAL 설정 쿼리
func UpdateAppliedScreeingForExtraVeFall(fallExtraVeOneseoIds []int) {
	configs.MyDB.Exec(`
		UPDATE tb_oneseo
		SET applied_screening = 'SPECIAL'
		WHERE oneseo_id IN ?
	`, fallExtraVeOneseoIds)
}

// special oneseo id조회 쿼리 (성적순으로 order)
func QuerySpecialOneseoIds() []int {
	var ids []int
	configs.MyDB.Raw(`
		SELECT o.oneseo_id 
		FROM tb_oneseo o
		JOIN tb_entrance_test_result tr ON o.oneseo_id = tr.oneseo_id
		JOIN tb_entrance_test_factors_detail td ON tr.entrance_test_factors_detail_id = td.entrance_test_factors_detail_id
		JOIN tb_member m ON o.member_id = m.member_id
		WHERE 
			o.applied_screening = 'SPECIAL'
			AND tr.second_test_pass_yn IS NULL
		ORDER BY 
			(((tr.document_evaluation_score / 3) * 0.5) + (tr.aptitude_evaluation_score * 0.3) + (tr.interview_score * 0.2)) DESC, 
			td.total_subjects_score DESC,
			(td.score_3_2 + td.score_3_1) DESC,
			(td.score_2_2 + td.score_2_1) DESC, 
			td.score_2_2 DESC, 
			td.score_2_1 DESC, 
			td.total_non_subjects_score DESC, 
			m.birth ASC;
	`).Scan(&ids)

	return ids
}

// special limit명 이하일때 second_test = pass
func UpdateSecondTestPassYnForSpecialPass(passSpecialOneseoIds []int) {

	// second_test_pass_yn = YES
	configs.MyDB.Exec(`
		UPDATE tb_entrance_test_result
		SET second_test_pass_yn = 'YES'
		WHERE oneseo_id IN ?
	`, passSpecialOneseoIds)
}

// special limit명 초과일때 하위 n명 applied_screening = general 설정 쿼리
func UpdateAppliedScreeingForSpecialFall(fallSpecialOneseoIds []int) {
	configs.MyDB.Exec(`
		UPDATE tb_oneseo
		SET applied_screening = 'GENERAL'
		WHERE oneseo_id IN ?
	`, fallSpecialOneseoIds)
}

// general 성적 상위 n명(limit 호출쪽에서 능동적으로) second_test = pass & 나머지 지원자 탈락 처리
func UpdateSecondTestPassYnForGeneral(generalPassLimit int) {

	// second_test_pass_yn = YES
	configs.MyDB.Exec(`
		UPDATE tb_entrance_test_result tr
		JOIN (
			SELECT tr.entrance_test_result_id 
			FROM tb_entrance_test_result tr
			JOIN tb_entrance_test_factors_detail td ON tr.entrance_test_factors_detail_id = td.entrance_test_factors_detail_id
			JOIN tb_oneseo o ON tr.oneseo_id = o.oneseo_id
			JOIN tb_member m ON o.member_id = m.member_id
			WHERE 
				o.applied_screening = 'GENERAL'
				AND tr.second_test_pass_yn IS NULL
			ORDER BY 
				(((tr.document_evaluation_score / 3) * 0.5) + (tr.aptitude_evaluation_score * 0.3) + (tr.interview_score * 0.2)) DESC, 
				td.total_subjects_score DESC,
				(td.score_3_2 + td.score_3_1) DESC,
				(td.score_2_2 + td.score_2_1) DESC, 
				td.score_2_2 DESC, 
				td.score_2_1 DESC, 
				td.total_non_subjects_score DESC, 
				m.birth ASC
			LIMIT ?
		) AS subquery ON tr.entrance_test_result_id = subquery.entrance_test_result_id
		SET tr.second_test_pass_yn = 'YES';
	`, generalPassLimit)

	// second_test_pass_yn = NO
	configs.MyDB.Exec(`
		UPDATE tb_entrance_test_result tr
		JOIN (
			SELECT tr.entrance_test_result_id 
			FROM tb_entrance_test_result tr
			JOIN tb_oneseo o ON tr.oneseo_id = o.oneseo_id
			WHERE 
				o.applied_screening IS NOT NULL
				AND tr.second_test_pass_yn IS NULL
		) AS subquery ON tr.entrance_test_result_id = subquery.entrance_test_result_id
		SET tr.second_test_pass_yn = 'NO';
	`)
}
