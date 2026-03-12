package dto

type AtualizarSettingsProjectRequest struct {
	ProjectID string
	LeaderID  string

	PermitirSoltarDoingParaTodo bool
}
