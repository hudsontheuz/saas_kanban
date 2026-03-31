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

	membros, err := uc.teams.ListarMembros(t.ID())
	if err != nil {
		return teamdto.BuscarTeamResponse{}, err
	}

	responseMembros := make([]teamdto.TeamMemberResponse, 0, len(membros))
	for _, membro := range membros {
		responseMembros = append(responseMembros, teamdto.TeamMemberResponse{
			UserID: string(membro.UserID),
			Nome:   membro.Nome,
			Email:  membro.Email,
			Role:   membro.Role,
		})
	}

	return teamdto.BuscarTeamResponse{
		TeamID:   string(t.ID()),
		Nome:     t.Nome(),
		LeaderID: string(t.LeaderID()),
		Membros:  responseMembros,
	}, nil
}
