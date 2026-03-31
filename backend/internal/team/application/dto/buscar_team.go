package dto

type BuscarTeamRequest struct {
	TeamID string
}

type TeamMemberResponse struct {
	UserID string `json:"user_id"`
	Nome   string `json:"nome"`
	Email  string `json:"email"`
	Role   string `json:"role"`
}

type BuscarTeamResponse struct {
	TeamID   string               `json:"team_id"`
	Nome     string               `json:"nome"`
	LeaderID string               `json:"leader_id"`
	Membros  []TeamMemberResponse `json:"membros"`
}
