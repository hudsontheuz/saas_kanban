package taskdto

type CriarTaskRequest struct {
	ProjectID string
	Titulo    string
	CriadorID string
}

type CriarTaskResponse struct {
	TaskID string
}
