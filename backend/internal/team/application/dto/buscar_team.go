package dto

type BuscarTeamRequest struct {
	TeamID string
}

type BuscarTeamResponse struct {
	TeamID   string `json:"team_id"`
	Nome     string `json:"nome"`
	LeaderID string `json:"leader_id"`
}
