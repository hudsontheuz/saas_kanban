package usecase

import (
	"github.com/hudsontheuz/saas_kanban/internal/team/application/dto"
	teamports "github.com/hudsontheuz/saas_kanban/internal/team/application/ports"
	"github.com/hudsontheuz/saas_kanban/internal/team/domain"
	"github.com/hudsontheuz/saas_kanban/internal/user/domain"
)

type CriarTeamUseCase struct {
	teams teamports.TeamRepository
}

func NovoCriarTeamUseCase(teams teamports.TeamRepository) *CriarTeamUseCase {
	return &CriarTeamUseCase{teams: teams}
}

func (uc *CriarTeamUseCase) Executar(req dto.CriarTeamRequest) (dto.CriarTeamResponse, error) {
	t, err := team.NovaTeam(req.Nome, user.UserID(req.LeaderID))
	if err != nil {
		return dto.CriarTeamResponse{}, err
	}

	if err := uc.teams.Salvar(t); err != nil {
		return dto.CriarTeamResponse{}, err
	}

	return dto.CriarTeamResponse{TeamID: string(t.ID())}, nil
}
