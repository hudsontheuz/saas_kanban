package task

type TaskID string
type StatusTask string
type OutcomeTask string

const (
	ToDo            StatusTask  = "TODO"
	Doing           StatusTask  = "DOING"
	InReview        StatusTask  = "INREVIEW"
	Done            StatusTask  = "DONE"
	OutcomeApproved OutcomeTask = "APPROVED"
	OutcomeRejected OutcomeTask = "REJECTED"
)
