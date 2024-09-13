package repository

import (
	"log"
	"themoment-team/go-hellogsm/configs"
	"themoment-team/go-hellogsm/jobs"
)

func CountOneseoByWantedScreening(wantedScreening string) int {
	var result int
	tx := configs.MyDB.Raw("select count(*) from tb_oneseo where wanted_screening = ?", wantedScreening).Scan(&result)
	if tx.Error != nil {
		log.Println(tx.Error.Error())
	}
	return result
}

func SaveAppliedScreening(wantedScreening jobs.Screening, appliedScreening jobs.Screening, top int) {
	configs.MyDB.Raw(`
update tb_oneseo
set applied_screening = ?
where oneseo_id in (select tbo.oneseo_id
                    from tb_oneseo tbo
                             join tb_entrance_test_result tbe on tbo.oneseo_id = tbe.oneseo_id
                    where tbo.wanted_screening = ?
                      and tbo.applied_screening is null
                    order by document_evaluation_score
                    limit ?)
`, &appliedScreening, &wantedScreening, &top)
}

func IsAppliedScreeningAllNull() bool {
	var result int
	configs.MyDB.Raw("select count(*) from tb_oneseo where applied_screening is not null").Scan(&result)
	return result < 1
}

func IsAppliedScreeningAllNullBy(wantedScreening jobs.Screening) bool {
	var totalCount int
	var nullCount int
	configs.MyDB.Raw("select count(*) from tb_oneseo where wanted_screening = ?", wantedScreening).Scan(&totalCount)
	configs.MyDB.Raw("select count(*) from tb_oneseo where wanted_screening = ? and applied_screening is null", wantedScreening).Scan(&nullCount)
	return totalCount == nullCount
}

func IsAppliedScreeningAllNotNullBy(wantedScreening jobs.Screening) bool {
	var totalCount int
	var notnullCount int
	configs.MyDB.Raw("select count(*) from tb_oneseo where wanted_screening = ?", wantedScreening).Scan(&totalCount)
	configs.MyDB.Raw("select count(*) from tb_oneseo where wanted_screening = ? and applied_screening is not null", wantedScreening).Scan(&notnullCount)
	return totalCount == notnullCount
}
