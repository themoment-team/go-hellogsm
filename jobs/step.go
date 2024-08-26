package jobs

type Step interface {
	Reader()
	Writer()
}
