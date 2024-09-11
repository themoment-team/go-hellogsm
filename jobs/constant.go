package jobs

type Screening string
type Major string

const (
	// 전형
	GeneralScreening        Screening = "GENERAL"
	SpecialScreening        Screening = "SPECIAL"
	ExtraVeteransScreening  Screening = "EXTRA_VETERANS"
	ExtraAdmissionScreening Screening = "EXTRA_ADMISSION"

	// [1차 평가] 전형 별 합격자 수
	GeneralSuccessfulApplicantOf1E        int = 84
	SpecialSuccessfulApplicantOf1E        int = 11
	ExtraVeteransSuccessfulApplicantOf1E  int = 2
	ExtraAdmissionSuccessfulApplicantOf1E int = 1

	// [2차 평가] 전형 별 합격자 수
	GeneralSpecialSuccessfulApplicantOf2E int = 72
	ExtraVeteransSuccessfulApplicantOf2E  int = 2
	ExtraAdmissionSuccessfulApplicantOf2E int = 1

	// 학과 별 정원
	SWDepartment    = 36
	IOTDepartment   = 18
	AIDepartment    = 18
	ExtraDepartment = 2

	// 학과
	SW  Major = "SW"
	IOT Major = "IOT"
	AI  Major = "AI"

	// 학과 배정시 정원외특별전형의 구분을 위한 값
	NORMAL = "NORMAL"
	EXTRA  = "EXTRA"

	// 그냥 전체. 발생하지 않을 수
	JustAll int = 99999
)
