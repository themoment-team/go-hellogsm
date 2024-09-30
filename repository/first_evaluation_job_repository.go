package repository

import (
	"fmt"
	"gorm.io/gorm"
	"log"
	"themoment-team/go-hellogsm/configs"
	e "themoment-team/go-hellogsm/error"
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
update tb_oneseo tbo
    join (select tbo_inner.oneseo_id
          from tb_oneseo tbo_inner
                   join tb_entrance_test_result tbe
                        on tbo_inner.oneseo_id = tbe.oneseo_id
          where tbo_inner.wanted_screening in ?
            and tbo_inner.applied_screening is null 
    		and real_oneseo_arrived_yn = 'YES'
          order by tbe.document_evaluation_score
          LIMIT ?) as limited_tbo
    on tbo.oneseo_id = limited_tbo.oneseo_id
set tbo.applied_screening = ?
where tbo.oneseo_id is not null
`)
	return e.WrapRollbackNeededError(db.Exec(query, evaluateScreening, top, appliedScreening).Error)
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
set tbe.first_test_pass_yn = IF(tbo.applied_screening is not null and tbo.real_oneseo_arrived_yn = 'YES', 'YES', 'NO')
where tbo.oneseo_id is not null;
`
	return e.WrapRollbackNeededError(db.Exec(query).Error)
}
