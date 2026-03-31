package usecase

import (
	teamdto "github.com/hudsontheuz/saas_kanban/internal/team/application/dto"
	teamports "github.com/hudsontheuz/saas_kanban/internal/team/application/ports"
	user "github.com/hudsontheuz/saas_kanban/internal/user/domain"
)

type ListarMinhasTeamsUseCase struct {
	teams teamports.TeamRepository
}

func NovoListarMinhasTeamsUseCase(teams teamports.TeamRepository) *ListarMinhasTeamsUseCase {
	return &ListarMinhasTeamsUseCase{teams: teams}
}

func (uc *ListarMinhasTeamsUseCase) Executar(req teamdto.ListarMinhasTeamsRequest) (teamdto.ListarMinhasTeamsResponse, error) {
	lista, err := uc.teams.ListarPorUsuarioID(user.UserID(req.UserID))
	if err != nil {
		return teamdto.ListarMinhasTeamsResponse{}, err
	}

	items := make([]teamdto.TeamListItem, 0, len(lista))
	for _, tm := range lista {
		items = append(items, teamdto.TeamListItem{
			TeamID:   string(tm.ID()),
			Nome:     tm.Nome(),
			LeaderID: string(tm.LeaderID()),
		})
	}

	return teamdto.ListarMinhasTeamsResponse{Items: items}, nil
}
