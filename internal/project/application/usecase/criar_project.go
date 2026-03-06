package usecase

import (
	"github.com/hudsontheuz/saas_kanban/internal/project/application/dto"
	projectports "github.com/hudsontheuz/saas_kanban/internal/project/application/ports"
	"github.com/hudsontheuz/saas_kanban/internal/project/domain"
	teamports "github.com/hudsontheuz/saas_kanban/internal/team/application/ports"
	"github.com/hudsontheuz/saas_kanban/internal/team/domain"
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

	p, err := project.NovoProject(teamID, req.Nome, project.ConfiguracoesProject{
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
