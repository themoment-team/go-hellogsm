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
	db, err := configs.MyDB.DB()
	if err != nil {
		fatalErr(err)
	}
	defer func() {
		if closeErr := db.Close(); closeErr != nil {
			log.Println("Error closing DB connection:", closeErr)
		}
	}()

	// DB transaction begin
	tx, err := db.Begin()
	if err != nil {
		fatalErr(err)
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
		err := step.Processor(stepContext)
		if err != nil {
			handle(tx, err)
		}
	}

	// DB transaction commit
	err = tx.Commit()
	if err != nil {
		fatalErr(err)
	}
}

func handle(tx *sql.Tx, err error) {
	var rollbackErr RollbackNeededError
	if errors.As(err, &rollbackErr) {
		log.Printf("An RollbackNeededError occurred. 상세: [%s]", err.Error())
		err := tx.Rollback()
		if err != nil {
			fatalErr(err)
		}
		log.Println("transaction rollback completed.")
		panic("RollbackNeededError 발생으로 작업을 더 이상 진행하지 않음.")
	} else {
		log.Println("An unknown error occurred:", err)
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

func fatalErr(err error) {
	log.Println("An fatal error occurred. Exiting...")
	panic(err)
}
