package internal

import "flag"

const (
	activeProfileFlag = "profile"
	jobsFlag          = "jobs"
)

var (
	activeProfile string
	jobs          string
)

func ProgramEssentialMeta() {
	parseFlagArgs()
}

func parseFlagArgs() {
	flag.StringVar(&activeProfile, activeProfileFlag, "local", "실행 환경을 선택해주세요.")
	flag.StringVar(&jobs, jobsFlag, "defaultJob", "실행할 Job을 선택해주세요.")
	flag.Parse()

	SetActiveProfile(activeProfile)
	SetJobs(jobs)
}
