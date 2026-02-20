package dto

type PausarTaskRequest struct {
	TaskID string
	UserID string
}

type RetomarTaskRequest struct {
	TaskID string
	UserID string
}
