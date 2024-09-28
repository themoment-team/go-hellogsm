package jobs

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"themoment-team/go-hellogsm/configs"
	"themoment-team/go-hellogsm/service"
	"time"
)

type SimpleJob struct {
	name        string
	steps       []Step
	jobListener *JobListener
}

func NewSimpleJob(name string, steps []Step, jobListener JobListener) *SimpleJob {
	return &SimpleJob{
		name:        name,
		steps:       steps,
		jobListener: &jobListener,
	}
}

func (job *SimpleJob) Name() string {
	return job.name
}

func (job *SimpleJob) Start() {
	db, _ := configs.MyDB.DB()
	txManager := service.NewTransactionManager(db)
	tx, err := txManager.BeginTx()
	if err != nil {
		WrapFatalErr(err)
	}

	listener := *job.jobListener
	if listener == nil {
		listener = DefaultJobListener{job.name, time.Now()}
	}
	listener.BeforeJob()
	// AfterJob() 은 어떠한 상황에서도 무조건 실행된다.
	defer listener.AfterJob()

	stepContext := NewBatchContext()
	// step 은 기본적으로 배치된 순서에 따라 실행한다.
	for _, step := range job.steps {
		err = step.Processor(stepContext, tx.(*sql.Tx))
		if err != nil {
			rollback := processRollbackIfNeeded(err, txManager, tx)
			if rollback {
				// rollback 이 발생한 경우 더 이상 진행하지 않는다.
				return
			}
		}

		err = txManager.Commit(tx)
		if err != nil {
			WrapFatalErr(err)
		}
	}
}

type DefaultJobListener struct {
	jobName   string
	startTime time.Time
}

func (l DefaultJobListener) BeforeJob() {
	startMsg := fmt.Sprintf("job: [%s] is started at [%s]", l.jobName, l.startTime.Format("2006-01-02 15:04:05"))
	service.SendDiscordMsg(service.Template{
		Title:       "배치 작업 시작 알림",
		Content:     startMsg,
		NoticeLevel: service.Info,
	})
	log.Println(startMsg)
}

func (l DefaultJobListener) AfterJob() {
	endTime := time.Now()
	endMsg := fmt.Sprintf("job: [%s] is finished at [%s] (elapsed time: %s)", l.jobName, endTime.Format("2006-01-02 15:04:05"), endTime.Sub(l.startTime))
	service.SendDiscordMsg(service.Template{
		Title:       "배치 작업 종료 알림",
		Content:     endMsg,
		NoticeLevel: service.Info,
	})
	log.Println(endMsg)
}

func processRollbackIfNeeded(err error, txManager service.TransactionManager, tx interface{}) bool {
	var rollbackNeededError RollbackNeededError
	if errors.As(err, &rollbackNeededError) {
		log.Println("error occurred and rollback.", err)
		txRollbackErr := txManager.Rollback(tx)
		if txRollbackErr != nil {
			WrapFatalErr(txRollbackErr)
		}
		return true
	} else {
		log.Println("error occurred but not rollback.", err)
		return false
	}
}
