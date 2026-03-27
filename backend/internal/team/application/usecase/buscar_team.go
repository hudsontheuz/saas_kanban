package usecase

import (
	teamdto "github.com/hudsontheuz/saas_kanban/internal/team/application/dto"
	teamports "github.com/hudsontheuz/saas_kanban/internal/team/application/ports"
	team "github.com/hudsontheuz/saas_kanban/internal/team/domain"
)

type BuscarTeamUseCase struct {
	teams teamports.TeamRepository
}

func NovoBuscarTeamUseCase(teams teamports.TeamRepository) *BuscarTeamUseCase {
	return &BuscarTeamUseCase{teams: teams}
}

func (uc *BuscarTeamUseCase) Executar(req teamdto.BuscarTeamRequest) (teamdto.BuscarTeamResponse, error) {
	t, err := uc.teams.BuscarPorID(team.TeamID(req.TeamID))
	if err != nil {
		return teamdto.BuscarTeamResponse{}, err
	}

	return teamdto.BuscarTeamResponse{
		TeamID:   string(t.ID()),
		Nome:     t.Nome(),
		LeaderID: string(t.LeaderID()),
	}, nil
}
