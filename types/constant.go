package types

type Screening string
type Major string

type Applicant struct {
	MemberID           int       `json:"member_id"`
	AppliedScreening   Screening `json:"applied_screening"`
	FirstDesiredMajor  Major     `json:"first_desired_major"`
	SecondDesiredMajor Major     `json:"second_desired_major"`
	ThirdDesiredMajor  Major     `json:"third_desired_major"`
}

const (
	// 전형
	GeneralScreening        Screening = "GENERAL"
	SpecialScreening        Screening = "SPECIAL"
	ExtraVeteransScreening  Screening = "EXTRA_VETERANS"
	ExtraAdmissionScreening Screening = "EXTRA_ADMISSION"

	// [1차 평가] 전형 별 합격자 수
	GeneralSuccessfulApplicantOf1E        int = 76
	SpecialSuccessfulApplicantOf1E        int = 8
	ExtraVeteransSuccessfulApplicantOf1E  int = 1
	ExtraAdmissionSuccessfulApplicantOf1E int = 1

	// [2차 평가] 전형 별 합격자 수
	GeneralSpecialSuccessfulApplicantOf2E int = 64
	GeneralSuccessfulApplicantOf2E        int = 58
	SpecialSuccessfulApplicantOf2E        int = 6
	ExtraVeteransSuccessfulApplicantOf2E  int = 1
	ExtraAdmissionSuccessfulApplicantOf2E int = 1

	// 학과 별 정원
	SWMajor    = 32
	IOTMajor   = 16
	AIMajor    = 16
	ExtraMajor = 2

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

// 고사실 정원 (고사장 상황 변동 시 수정 필요)
// 06번 고사실은 정원 외 특별전형(최대 2명)까지 포함한 정원으로 다른 고사실보다 2명 더 크게 잡는다.
var ExaminationRoomCapacities = []int{15, 15, 15, 15, 12, 14}

// TotalExaminationRoomCapacity 전체 고사실 정원 합계를 반환한다.
func TotalExaminationRoomCapacity() int {
	total := 0
	for _, c := range ExaminationRoomCapacities {
		total += c
	}
	return total
}
