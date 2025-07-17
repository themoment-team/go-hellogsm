package shared

import "fmt"

var (
	CommonTieBreakerQuery = `
	td.general_subjects_score DESC, -- 일반교과성적이 우수한자
	td.score_3_1 DESC, -- 3-1,2-2,2-1,1-2 순으로 성적이 우수한자
	td.score_2_2 DESC,
	td.score_2_1 DESC,
	td.score_1_2 DESC,
	td.total_non_subjects_score DESC -- 비교과성적이 우수한자
`
	FinalTieBreakerQuery = fmt.Sprintf(`
	%s
	tr.competency_evaluation_score DESC, -- 역량검사 점수가 우수한자
	tr.interview_score DESC, -- 면접 점수가 우수한자
`, CommonTieBreakerQuery)
)
