package repository

import (
	"log"
	"themoment-team/go-hellogsm/configs"
)

func CountByGiveUpApplicant() int {
	var result int
	tx := configs.MyDB.Raw(`
							SELECT count(*) 
							FROM tb_entrance_test_result 
							WHERE second_test_pass_yn = 'YES' AND 
							      entrance_intention_yn = 'NO'
					`).Scan(&result)
	if tx.Error != nil {
		log.Println(tx.Error.Error())
	}
	return result
}

func CountByFinalPassApplicant() int {
	var result int
	tx := configs.MyDB.Raw(`
							SELECT count(*) 
							FROM tb_entrance_test_result 
							WHERE second_test_pass_yn = 'YES' AND
							      entrance_intention_yn = NULL
					`).Scan(&result)
	if tx.Error != nil {
		log.Println(tx.Error.Error())
	}
	return result
}

func QueryByRemainingDepartment() (int, int, int) {
	var sw, iot, ai int

	rows, err := configs.MyDB.Raw(`
        SELECT decided_major, COUNT(*)
        FROM tb_oneseo
        WHERE decided_major IS NOT NULL
        GROUP BY decided_major
    `).Rows()
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
