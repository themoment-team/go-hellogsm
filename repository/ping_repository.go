package repository

import (
	"gorm.io/gorm"
)

func SelectOne(db *gorm.DB) int {
	var result int
	tx := db.Raw("select 1").Scan(&result)
	if tx.Error != nil {
		panic("sql execute error")
	}
	return result
}
