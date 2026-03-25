package dto

type BuscarProjectAtivoRequest struct {
	TeamID string
}

type BuscarProjectAtivoResponse struct {
	ProjectID                   string `json:"project_id"`
	TeamID                      string `json:"team_id"`
	Nome                        string `json:"nome"`
	Status                      string `json:"status"`
	PermitirSoltarDoingParaTodo bool   `json:"permitir_soltar_doing_para_todo"`
}
