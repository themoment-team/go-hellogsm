package test

import (
	"database/sql"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
	"testing"
	"themoment-team/go-hellogsm/configs"
	"themoment-team/go-hellogsm/internal"
	"themoment-team/go-hellogsm/jobs"
)

func beforeTest() {
	// test env setting
	setActiveProfile()
	setMysqlDB()

	// create mock table for test
	createTable()
}

func afterTest() {
	configs.MyDB.Exec("drop table tb_a")
}

func TestTxRollback(t *testing.T) {
	// given
	beforeTest()

	// when
	jobs.NewSimpleJob("testJob", getSteps(), getJobListener()).Start()

	// then
	actual := findAll()
	assert.Equal(t, nil, actual[0])

	afterTest()
}

func findAll() []interface{} {
	list := make([]interface{}, 1)
	configs.MyDB.Raw("select id from tb_a").Scan(&list)
	return list
}

type TestStep struct {
}

func (t TestStep) Processor(bc *jobs.BatchContext, tx *sql.Tx) error {
	err := doA(tx)
	if err != nil {
		return err
	}

	err = doXReturnRollbackErr()
	if err != nil {
		return err
	}

	return nil
}

func doA(tx *sql.Tx) error {
	_, err := tx.Exec("insert into tb_a values (1)")
	return err
}

// rollback 을 해야하는 에러를 반환한다.
func doXReturnRollbackErr() jobs.RollbackNeededError {
	return jobs.WrapRollbackNeededError("error occurred")
}

func getSteps() []jobs.Step {
	return []jobs.Step{&TestStep{}}
}

func getJobListener() TestJobListener {
	return TestJobListener{}
}

// 테스트에서는 Job 공통 Listener 기능을 사용하지 않기 위해 별도로 만든다.
type TestJobListener struct {
}

func (t TestJobListener) BeforeJob() {
	// do nothing
}

func (t TestJobListener) AfterJob() {
	// do nothing
}

func createTable() *gorm.DB {
	return configs.MyDB.Exec("create table if not exists tb_a (id int)")
}

func setMysqlDB() {
	// test DB 에 해당하는 설정을 기입한다.
	configs.CreateMysqlDB(configs.CreateMysqlDsn(
		internal.MysqlProperties{
			Host:     "localhost",
			Port:     "3306",
			Username: "root",
			Password: "1234",
			Database: "hellogsm",
		},
	))
}

func setActiveProfile() {
	internal.SetActiveProfile("local")
}
