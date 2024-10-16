package repository

import (
	"fmt"
	"gorm.io/gorm"
	"log"
	"themoment-team/go-hellogsm/configs"
	e "themoment-team/go-hellogsm/error"
	"themoment-team/go-hellogsm/types"
)

func CountByGiveUpApplicant() int {
	result := 0
	tx := configs.MyDB.Raw(`
		SELECT COALESCE(COUNT(*), 0) 
		FROM tb_entrance_test_result tr JOIN tb_oneseo o ON tr.oneseo_id = o.oneseo_id 
		WHERE tr.second_test_pass_yn = 'YES' AND 
			  o.entrance_intention_yn = 'NO' 
	`).Scan(&result)
	if tx.Error != nil {
		log.Println(tx.Error.Error())
	}
	return result
}

func CountFinalTestPassNormalApplicant() int {
	result := 0
	tx := configs.MyDB.Raw(`
		SELECT COALESCE(COUNT(*), 0) 
		FROM tb_entrance_test_result tr JOIN tb_oneseo o ON tr.oneseo_id = o.oneseo_id 
		WHERE tr.second_test_pass_yn = 'YES' AND 
		      (o.applied_screening = 'GENERAL' OR o.applied_screening = 'SPECIAL')
	`).Scan(&result)
	if tx.Error != nil {
		log.Println(tx.Error.Error())
	}
	return result
}

func CountFinalTestPassExtraApplicant() int {
	result := 0
	tx := configs.MyDB.Raw(`
		SELECT COALESCE(COUNT(*), 0) 
		FROM tb_entrance_test_result tr JOIN tb_oneseo o ON tr.oneseo_id = o.oneseo_id 
		WHERE tr.second_test_pass_yn = 'YES' AND 
		      (o.applied_screening = 'EXTRA_VETERANS' OR o.applied_screening = 'EXTRA_ADMISSION')
	`).Scan(&result)
	if tx.Error != nil {
		log.Println(tx.Error.Error())
	}
	return result
}

func QueryByScreeningsAssignedMajor(firstScreening types.Screening, secondScreening types.Screening) (int, int, int) {
	sw := 0
	iot := 0
	ai := 0

	rows, err := configs.MyDB.Raw(`
        SELECT decided_major, COALESCE(COUNT(*), 0) 
        FROM tb_oneseo
        WHERE decided_major IS NOT NULL AND 
              entrance_intention_yn = 'NO' AND
              applied_screening IN (?, ?) 
        GROUP BY decided_major
    `, firstScreening, secondScreening).Rows()
	if err != nil {
		log.Println(err)
	}
	defer rows.Close()

	for rows.Next() {
		var major string
		var count int
		if err := rows.Scan(&major, &count); err != nil {
			log.Println(err)
			continue
		}

		switch major {
		case "SW":
			sw = count
		case "IOT":
			iot = count
		case "AI":
			ai = count
		}
	}

	return sw, iot, ai
}

func QueryAllByFinalTestPassApplicant() (error, []types.Applicant) {
	var applicants []types.Applicant

	rows, err := configs.MyDB.Raw(`
		SELECT m.member_id, o.applied_screening, o.first_desired_major, o.second_desired_major, o.third_desired_major 
		FROM tb_member m 
		JOIN tb_oneseo o ON m.member_id = o.member_id
		JOIN tb_entrance_test_result tr ON tr.oneseo_id = o.oneseo_id
		JOIN tb_entrance_test_factors_detail td ON tr.entrance_test_factors_detail_id = td.entrance_test_factors_detail_id
		WHERE tr.second_test_pass_yn = 'YES' AND 
		      o.entrance_intention_yn IS NULL AND 
		      m.role = 'APPLICANT'
		ORDER BY 
		(((tr.document_evaluation_score / 3) * 0.5) + (tr.aptitude_evaluation_score * 0.3) + (tr.interview_score * 0.2)) DESC, 
		td.total_subjects_score DESC, 
		(td.score_3_2 + td.score_3_1) DESC,
		(td.score_2_2 + td.score_2_1) DESC, 
		td.score_2_2 DESC, 
		td.score_2_1 DESC, 
		td.total_non_subjects_score DESC, 
		m.birth ASC;
	`).Rows()

	if err != nil {
		log.Println(err)
		return err, nil
	}
	defer rows.Close()

	for rows.Next() {
		var applicant types.Applicant
		if err := rows.Scan(&applicant.MemberID, &applicant.AppliedScreening, &applicant.FirstDesiredMajor, &applicant.SecondDesiredMajor, &applicant.ThirdDesiredMajor); err != nil {
			log.Println(err)
			return err, nil
		}
		applicants = append(applicants, applicant)
	}

	if err := rows.Err(); err != nil {
		log.Println(err)
		return err, nil
	}

	return nil, applicants
}

func QueryAllByAdditionalApplicant() (error, []types.Applicant) {
	var applicants []types.Applicant

	rows, err := configs.MyDB.Raw(`
		SELECT m.member_id, o.applied_screening, o.first_desired_major, o.second_desired_major, o.third_desired_major 
		FROM tb_member m 
		JOIN tb_oneseo o ON m.member_id = o.member_id
		JOIN tb_entrance_test_result tr ON tr.oneseo_id = o.oneseo_id
		JOIN tb_entrance_test_factors_detail td ON tr.entrance_test_factors_detail_id = td.entrance_test_factors_detail_id
		WHERE tr.second_test_pass_yn = 'NO' AND
		      o.entrance_intention_yn IS NULL AND 
		      m.role = 'APPLICANT'
		ORDER BY 
		(((tr.document_evaluation_score / 3) * 0.5) + (tr.aptitude_evaluation_score * 0.3) + (tr.interview_score * 0.2)) DESC, 
		td.total_subjects_score DESC, 
		(td.score_3_2 + td.score_3_1) DESC,
		(td.score_2_2 + td.score_2_1) DESC, 
		td.score_2_2 DESC, 
		td.score_2_1 DESC, 
		td.total_non_subjects_score DESC, 
		m.birth ASC;
	`).Rows()

	if err != nil {
		log.Println(err)
		return err, nil
	}
	defer rows.Close()

	for rows.Next() {
		var applicant types.Applicant
		if err := rows.Scan(&applicant.MemberID, &applicant.AppliedScreening, &applicant.FirstDesiredMajor, &applicant.SecondDesiredMajor, &applicant.ThirdDesiredMajor); err != nil {
			log.Println(err)
			return err, nil
		}
		applicants = append(applicants, applicant)
	}

	if err := rows.Err(); err != nil {
		log.Println(err)
		return err, nil
	}

	return nil, applicants
}

func UpdateDecideMajor(db *gorm.DB, decideMajor types.Major, memberId int) error {
	query := fmt.Sprintf(`
		UPDATE tb_oneseo 
		SET decided_major = ?
		WHERE member_id = ?;
	`)

	err := db.Exec(query, decideMajor, memberId).Error
	if err != nil {
		return e.WrapRollbackNeededError(err)
	}

	return nil
}
