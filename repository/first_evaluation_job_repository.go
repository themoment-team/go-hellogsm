package repository

import (
	"fmt"

	"gorm.io/gorm"
)

func CountOneseoByWantedScreening(wantedScreening string, tx *gorm.DB) (int, error) {
	var result int
	err := tx.Raw("select count(*) from tb_oneseo where wanted_screening = ?", wantedScreening).Scan(&result).Error
	if err != nil {
		return 0, err
	}
	return result, nil
}

func CountOneseoByAppliedScreening(appliedScreening string, tx *gorm.DB) (int, error) {
	var result int
	err := tx.Raw("select count(*) from tb_oneseo where applied_screening = ?", appliedScreening).Scan(&result).Error
	if err != nil {
		return 0, err
	}
	return result, nil
}

func SaveAppliedScreening(evaluateScreening []string, appliedScreening string, top int, tx *gorm.DB) error {
	query := fmt.Sprintf(`
update tb_oneseo tbo
    join (select tbo_inner.oneseo_id
          from tb_oneseo tbo_inner
                   join tb_entrance_test_result tbe
                        on tbo_inner.oneseo_id = tbe.oneseo_id
          where tbo_inner.wanted_screening in ?
            and tbo_inner.applied_screening is null
          order by tbe.document_evaluation_score
          LIMIT ?) as limited_tbo
    on tbo.oneseo_id = limited_tbo.oneseo_id
set tbo.applied_screening = ?
where tbo.oneseo_id is not null
`)

	if err := tx.Exec(query, evaluateScreening, top, appliedScreening).Error; err != nil {
		return err
	}

	return nil
}

func IsAppliedScreeningAllNull(tx *gorm.DB) (bool, error) {
	var result int
	err := tx.Raw("select count(*) from tb_oneseo where applied_screening is not null").Scan(&result).Error
	if err != nil {
		return false, err
	}
	return result < 1, nil
}

func IsAppliedScreeningAllNullBy(wantedScreening string, tx *gorm.DB) (bool, error) {
	var totalCount int
	var nullCount int

	err := tx.Raw("select count(*) from tb_oneseo where wanted_screening = ?", wantedScreening).Scan(&totalCount).Error
	if err != nil {
		return false, err
	}

	err = tx.Raw("select count(*) from tb_oneseo where wanted_screening = ? and applied_screening is null", wantedScreening).Scan(&nullCount).Error
	if err != nil {
		return false, err
	}

	return totalCount == nullCount, nil
}

func SaveFirstTestPassYn(tx *gorm.DB) error {
	query := `
update tb_entrance_test_result tbe
    join tb_oneseo tbo on tbe.oneseo_id = tbo.oneseo_id
set tbe.first_test_pass_yn = IF(tbo.applied_screening is not null, 'YES', 'NO')
where tbo.oneseo_id is not null;
`

	if err := tx.Exec(query).Error; err != nil {
		return err
	}

	return nil
}
