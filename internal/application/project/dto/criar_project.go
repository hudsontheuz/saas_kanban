package dto

type CriarProjectRequest struct {
	TeamID   string
	LeaderID string

	Nome string

	PermitirSoltarDoingParaTodo bool
}

type CriarProjectResponse struct {
	ProjectID string
}
