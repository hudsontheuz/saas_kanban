package usecase

import (
	"github.com/hudsontheuz/saas_kanban/internal/application/project/dto"
	projectports "github.com/hudsontheuz/saas_kanban/internal/application/project/ports"
	teamports "github.com/hudsontheuz/saas_kanban/internal/application/team/ports"
	"github.com/hudsontheuz/saas_kanban/internal/domain/project"
	"github.com/hudsontheuz/saas_kanban/internal/domain/team"
)

type CriarProjectUseCase struct {
	teams    teamports.TeamRepository
	projects projectports.ProjectRepository
}

func NovoCriarProjectUseCase(teams teamports.TeamRepository, projects projectports.ProjectRepository) *CriarProjectUseCase {
	return &CriarProjectUseCase{teams: teams, projects: projects}
}

func (uc *CriarProjectUseCase) Executar(req dto.CriarProjectRequest) (dto.CriarProjectResponse, error) {
	teamID := team.TeamID(req.TeamID)
	leaderID := team.UserID(req.LeaderID)

	tm, err := uc.teams.BuscarPorID(teamID)
	if err != nil {
		return dto.CriarProjectResponse{}, err
	}
	if !tm.EhLeader(leaderID) {
		return dto.CriarProjectResponse{}, ErrSomenteLeaderPodeGerenciarProject
	}

	if _, err := uc.projects.BuscarAtivoPorTeamID(teamID); err == nil {
		return dto.CriarProjectResponse{}, ErrJaExisteProjectAtivo
	}

	p, err := project.NovoProject(teamID, project.ConfiguracoesProject{
		PermitirSoltarDoingParaTodo: req.PermitirSoltarDoingParaTodo,
	})
	if err != nil {
		return dto.CriarProjectResponse{}, err
	}

	if err := uc.projects.Salvar(p); err != nil {
		return dto.CriarProjectResponse{}, err
	}

	return dto.CriarProjectResponse{ProjectID: string(p.ID())}, nil
}
