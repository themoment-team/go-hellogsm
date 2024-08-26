package repository

import (
	"themoment-team/go-hellogsm/configs"
)

func SelectOne() int {
	var result int
	tx := configs.MyDB.Raw("select 1").Scan(&result)
	if tx.Error != nil {
		panic("sql execute error")
	}
	return result
}
