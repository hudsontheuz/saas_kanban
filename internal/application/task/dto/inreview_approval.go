package dto

type MoverParaInReviewRequest struct {
	TaskID string
	UserID string // assignee
}

type AprovarTaskRequest struct {
	TaskID   string
	LeaderID string
}

type ReprovarTaskRequest struct {
	TaskID   string
	LeaderID string
}

type RejeitarTaskToDoRequest struct {
	TaskID   string
	LeaderID string
}