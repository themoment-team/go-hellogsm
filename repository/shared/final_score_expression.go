package shared

import "fmt"

// 2차 전형 최종 성적 반영 비율이다. 매년 입학전형요강에 따라 바뀔 수 있으므로,
// 비율이 바뀌면 이 상수만 고치면 second_evaluation_job_repository.go와
// major_assignment_job_repository.go의 정렬 기준이 함께 바뀐다.
const (
	DocumentWeight   = 0.6 // 서류 60%
	CompetencyWeight = 0.2 // 역량검사 20%
	InterviewWeight  = 0.2 // 면접 20%
)

// document_evaluation_score는 교과 성적 3개 요소(3-1, 2-2, 2-1, 1-2 중 상위 3개 학기 등)의
// 합산 점수이므로 /3으로 100점 만점 기준으로 정규화한 뒤 반영 비율을 곱한다.
var FinalScoreExpr = fmt.Sprintf(
	`(((tr.document_evaluation_score / 3) * %v) + (tr.competency_evaluation_score * %v) + (tr.interview_score * %v))`,
	DocumentWeight, CompetencyWeight, InterviewWeight,
)
