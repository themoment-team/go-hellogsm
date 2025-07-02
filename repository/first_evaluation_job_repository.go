package repository

import (
	"fmt"
	"log"
	"themoment-team/go-hellogsm/configs"
	e "themoment-team/go-hellogsm/error"
	"themoment-team/go-hellogsm/repository/shared"

	"gorm.io/gorm"
)

func CountOneseoByWantedScreening(wantedScreening string) int {
	var result int
	tx := configs.MyDB.Raw("select count(*) from tb_oneseo where wanted_screening = ? and real_oneseo_arrived_yn = 'YES'", wantedScreening).Scan(&result)
	if tx.Error != nil {
		log.Println(tx.Error.Error())
	}
	return result
}

func SaveAppliedScreening(db *gorm.DB, evaluateScreening []string, appliedScreening string, top int) error {
	query := fmt.Sprintf(`
		UPDATE tb_oneseo tbo
		JOIN (
			SELECT tbo_inner.oneseo_id
			FROM tb_oneseo tbo_inner
			JOIN tb_member tbm 
				ON tbo_inner.member_id = tbm.member_id
			JOIN tb_entrance_test_result tbe 
				ON tbo_inner.oneseo_id = tbe.oneseo_id
			JOIN tb_entrance_test_factors_detail tbd 
				ON tbe.entrance_test_result_id = tbd.entrance_test_factors_detail_id
			WHERE tbo_inner.wanted_screening IN ?
			  AND tbo_inner.applied_screening IS NULL
			  AND tbo_inner.real_oneseo_arrived_yn = 'YES'
			ORDER BY 
				tbe.document_evaluation_score DESC, -- 1차전형 점수
				? -- 동점자 처리기준 (TieBreakerQuery)
			LIMIT ?
		) AS limited_tbo
		ON tbo.oneseo_id = limited_tbo.oneseo_id
		SET tbo.applied_screening = ?
		WHERE tbo.oneseo_id IS NOT NULL;
`)
	err := db.Exec(query, evaluateScreening, shared.CommonTieBreakerQuery, top, appliedScreening).Error
	if err != nil {
		return e.WrapRollbackNeededError(err)
	}

	return nil
}

func IsAppliedScreeningAllNull() bool {
	var result int
	configs.MyDB.Raw("select count(*) from tb_oneseo where applied_screening is not null and real_oneseo_arrived_yn = 'YES'").Scan(&result)
	return result < 1
}

func IsAppliedScreeningAllNullBy(wantedScreening string) bool {
	var totalCount int
	var nullCount int
	configs.MyDB.Raw("select count(*) from tb_oneseo where wanted_screening = ? and real_oneseo_arrived_yn = 'YES'", wantedScreening).Scan(&totalCount)
	configs.MyDB.Raw("select count(*) from tb_oneseo where wanted_screening = ? and applied_screening is null and real_oneseo_arrived_yn = 'YES'", wantedScreening).Scan(&nullCount)
	return totalCount == nullCount
}

func SaveFirstTestPassYn(db *gorm.DB) error {
	query := `
update tb_entrance_test_result tbe
    join tb_oneseo tbo on tbe.oneseo_id = tbo.oneseo_id
set tbe.first_test_pass_yn = IF(tbo.applied_screening is not null and tbo.real_oneseo_arrived_yn = 'YES', 'YES', 'NO'),
    tbo.pass_yn = IF(tbo.applied_screening is not null and tbo.real_oneseo_arrived_yn = 'YES', null, 'NO')
where tbo.oneseo_id is not null;
`
	err := db.Exec(query).Error
	if err != nil {
		return e.WrapRollbackNeededError(err)
	}

	return nil
}
