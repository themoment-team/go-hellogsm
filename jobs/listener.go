package jobs

type JobListener interface {
	BeforeJob()
	AfterJob()
}
