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
			tr.first_test_pass_yn = 'YES'
			AND tr.second_test_pass_yn IS NULL
			AND o.wanted_screening = 'EXTRA_ADMISSION' 
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

// extra admission limit명 이하일때 second_test = pass & applied_screening = extra_admission 설정 쿼리
func UpdateSecondTestPassYnForExtraAdPass(passExtraAdOneseoIds []int) {

	// second_test_pass_yn = YES
	configs.MyDB.Raw(`
		UPDATE tb_entrance_test_result
		SET second_test_pass_yn = 'YES'
		WHERE oneseo_id IN ?
	`, &passExtraAdOneseoIds)

	// applied_screening = EXTRA_ADMISSION
	configs.MyDB.Raw(`
		UPDATE tb_oneseo
		SET applied_screening = 'EXTRA_ADMISSION'
		WHERE oneseo_id IN ?
	`, &passExtraAdOneseoIds)
}

// extra admission limit명 초과일때 하위 n명 applied_screening = SPECIAL 설정 쿼리
func UpdateAppliedScreeingForExtraAdFall(fallExtraAdOneseoIds []int) {
	configs.MyDB.Raw(`
		UPDATE tb_oneseo
		SET applied_screening = 'SPECIAL'
		WHERE oneseo_id IN ?
	`, &fallExtraAdOneseoIds)
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
			tr.first_test_pass_yn = 'YES'
			AND tr.second_test_pass_yn IS NULL
			AND o.wanted_screening = 'EXTRA_VETERANS' 
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

// extra veteran limit명 이하일때 second_test = pass & applied_screening = extra_veteran 설정 쿼리
func UpdateSecondTestPassYnForExtraVePass(passExtraVeOneseoIds []int) {

	// second_test_pass_yn = YES
	configs.MyDB.Raw(`
		UPDATE tb_entrance_test_result
		SET second_test_pass_yn = 'YES'
		WHERE oneseo_id IN ?
	`, &passExtraVeOneseoIds)

	// applied_screening = EXTRA_VETERANS
	configs.MyDB.Raw(`
		UPDATE tb_oneseo
		SET applied_screening = 'EXTRA_VETERANS'
		WHERE oneseo_id IN ?
	`, &passExtraVeOneseoIds)
}

// extra veteran limit명 초과일때 하위 n명 applied_screening = SPECIAL 설정 쿼리
func UpdateAppliedScreeingForExtraVeFall(fallExtraVeOneseoIds []int) {
	configs.MyDB.Raw(`
		UPDATE tb_oneseo
		SET applied_screening = 'SPECIAL'
		WHERE oneseo_id IN ?
	`, &fallExtraVeOneseoIds)
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
			tr.first_test_pass_yn = 'YES'
			AND tr.second_test_pass_yn IS NULL
			AND (o.wanted_screening = 'SPECIAL' OR o.applied_screening = 'SPECIAL')
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

// special limit명 이하일때 second_test = pass & applied_screening = special 설정 쿼리
func UpdateSecondTestPassYnForSpecialPass(passSpecialOneseoIds []int) {

	// second_test_pass_yn = YES
	configs.MyDB.Raw(`
		UPDATE tb_entrance_test_result
		SET second_test_pass_yn = 'YES'
		WHERE oneseo_id IN ?
	`, &passSpecialOneseoIds)

	// applied_screening = special
	configs.MyDB.Raw(`
		UPDATE tb_oneseo
		SET applied_screening = 'SPECIAL'
		WHERE oneseo_id IN ?
	`, &passSpecialOneseoIds)
}

// special limit명 초과일때 하위 n명 applied_screening = general 설정 쿼리
func UpdateAppliedScreeingForSpecialFall(fallSpecialOneseoIds []int) {
	configs.MyDB.Raw(`
		UPDATE tb_oneseo
		SET applied_screening = 'GENERAL'
		WHERE oneseo_id IN ?
	`, &fallSpecialOneseoIds)
}

// general 성적 상위 n명(limit 호출쪽에서 능동적으로) second_test = pass & applied_screening = general 설정 쿼리
func UpdateSecondTestPassYnForGeneral(generalPassLimit int) {

	// second_test_pass_yn = YES
	configs.MyDB.Raw(`
		UPDATE tb_entrance_test_result tr
		JOIN (
			SELECT tr.entrance_test_result_id 
			FROM tb_entrance_test_result tr
			JOIN tb_entrance_test_factors_detail td ON tr.entrance_test_factors_detail_id = td.entrance_test_factors_detail_id
			JOIN tb_oneseo o ON tr.oneseo_id = o.oneseo_id
			JOIN tb_member m ON o.member_id = m.member_id
			WHERE 
				(o.wanted_screening = 'GENERAL' OR o.applied_screening = 'GENERAL')
				AND tr.first_test_pass_yn = 'YES'
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
	`, &generalPassLimit)

	// applied_screening = general
	configs.MyDB.Raw(`
		UPDATE tb_oneseo
		JOIN (
			SELECT o.oneseo_id
			FROM tb_oneseo o
			JOIN tb_entrance_test_result tr ON o.oneseo_id = tr.oneseo_id
			JOIN tb_entrance_test_factors_detail td ON tr.entrance_test_factors_detail_id = td.entrance_test_factors_detail_id
			JOIN tb_member m ON o.member_id = m.member_id
			WHERE 
				o.wanted_screening = 'GENERAL'
				AND tr.first_test_pass_yn = 'YES'
				AND tr.second_test_pass_yn = 'YES'
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
		) AS subquery ON o.oneseo_id = subquery.oneseo_id
		SET o.applied_screening = 'GENERAL'
`, &generalPassLimit)
}
