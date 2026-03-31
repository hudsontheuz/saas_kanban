package dto

type ListarTasksPorProjectRequest struct {
	ProjectID string
	Status    *string
}

type TaskListItem struct {
	TaskID            string  `json:"task_id"`
	ProjectID         string  `json:"project_id"`
	Titulo            string  `json:"titulo"`
	Descricao         string  `json:"descricao,omitempty"`
	ComentarioReview  string  `json:"comentario_review,omitempty"`
	ComentarioEntrega string  `json:"comentario_entrega"`
	Status            string  `json:"status"`
	AssigneeID        *string `json:"assignee_id,omitempty"`
	Paused            bool    `json:"paused"`
	Outcome           *string `json:"outcome,omitempty"`
}

type ListarTasksPorProjectResponse struct {
	Items []TaskListItem `json:"items"`
}
