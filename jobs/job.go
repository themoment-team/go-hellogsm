package jobs

import (
	"fmt"
	"log"
	"themoment-team/go-hellogsm/configs"
	"themoment-team/go-hellogsm/service"
	"time"

	"gorm.io/gorm"
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
	(&configs.MyDB).Transaction(func(tx *gorm.DB) error {
		listener := *job.jobListener
		if listener == nil {
			listener = DefaultJobListener{job.name, time.Now()}
		}
		listener.BeforeJob()

		stepContext := NewBatchContext()
		// step 은 기본적으로 배치된 순서에 따라 실행한다.
		for _, step := range job.steps {
			step.Processor(stepContext)
		}

		listener.AfterJob()
		return nil
	})
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
