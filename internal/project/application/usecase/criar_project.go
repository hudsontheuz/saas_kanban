package usecase

import (
	"errors"

	"github.com/hudsontheuz/saas_kanban/internal/project/application/dto"
	projectports "github.com/hudsontheuz/saas_kanban/internal/project/application/ports"
	project "github.com/hudsontheuz/saas_kanban/internal/project/domain"
	shared "github.com/hudsontheuz/saas_kanban/internal/shared/errors"
	teamports "github.com/hudsontheuz/saas_kanban/internal/team/application/ports"
	team "github.com/hudsontheuz/saas_kanban/internal/team/domain"
	user "github.com/hudsontheuz/saas_kanban/internal/user/domain"
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
	leaderID := user.UserID(req.LeaderID)

	tm, err := uc.teams.BuscarPorID(teamID)
	if err != nil {
		return dto.CriarProjectResponse{}, err
	}
	if !tm.EhLeader(leaderID) {
		return dto.CriarProjectResponse{}, ErrSomenteLeaderPodeGerenciarProject
	}

	if _, err := uc.projects.BuscarAtivoPorTeamID(teamID); err == nil {
		return dto.CriarProjectResponse{}, ErrJaExisteProjectAtivo
	} else if !errors.Is(err, shared.ErrNaoEncontrado) {
		return dto.CriarProjectResponse{}, err
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
