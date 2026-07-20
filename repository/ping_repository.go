package repository

import (
	"themoment-team/go-hellogsm/configs"
)

func SelectOne() (*int, error) {
	var result int
	tx := configs.MyDB.Raw("select 1").Scan(&result)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return &result, nil
}
