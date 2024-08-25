package internal

import (
	"os"
)

const activeProfileKey = "active_profile"

func GetActiveProfile() AppProfile {
	value := os.Getenv(activeProfileKey)
	for _, profile := range getAllProfiles() {
		if value == profile.Value {
			return profile
		}
	}
	panic("program을 실행하랴면 profile 설정이 필요합니다.")
}

func SetActiveProfile(activeProfile string) {
	os.Setenv(activeProfileKey, activeProfile)
}

func getAllProfiles() []AppProfile {
	return []AppProfile{Prod, Stage, Local}
}

var Prod = AppProfile{
	Value: "prod",
	Desc:  "상용환경",
}

var Stage = AppProfile{
	Value: "stage",
	Desc:  "개발환경",
}

var Local = AppProfile{
	Value: "local",
	Desc:  "로컬환경",
}

type AppProfile struct {
	Value string
	Desc  string
}
