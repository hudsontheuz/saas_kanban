package dto

type ListarMinhasTeamsRequest struct {
	UserID string
}

type TeamListItem struct {
	TeamID   string `json:"team_id"`
	Nome     string `json:"nome"`
	LeaderID string `json:"leader_id"`
}

type ListarMinhasTeamsResponse struct {
	Items []TeamListItem `json:"items"`
}
