package team

import user "github.com/hudsontheuz/saas_kanban/internal/user/domain"

func HidratarTeam(id TeamID, nome string, leaderID user.UserID) *Team {
	return &Team{
		id:       id,
		nome:     nome,
		leaderID: leaderID,
	}
}
