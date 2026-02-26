package dto

type CriarTeamRequest struct {
	Nome     string
	LeaderID string
}

type CriarTeamResponse struct {
	TeamID string
}
