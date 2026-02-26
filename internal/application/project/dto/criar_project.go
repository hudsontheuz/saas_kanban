package dto

type CriarProjectRequest struct {
	TeamID   string
	LeaderID string

	PermitirSoltarDoingParaTodo bool
}

type CriarProjectResponse struct {
	ProjectID string
}
