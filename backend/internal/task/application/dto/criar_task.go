package dto

type CriarTaskRequest struct {
	ProjectID string
	Titulo    string
	Descricao string
	CriadorID string
}

type CriarTaskResponse struct {
	TaskID string `json:"task_id"`
}
