package usecase

import (
	"time"

	"github.com/hudsontheuz/saas_kanban/internal/application/project/dto"
	projectports "github.com/hudsontheuz/saas_kanban/internal/application/project/ports"
	teamports "github.com/hudsontheuz/saas_kanban/internal/application/team/ports"
	"github.com/hudsontheuz/saas_kanban/internal/domain/project"
	"github.com/hudsontheuz/saas_kanban/internal/domain/team"
)

type FecharProjectUseCase struct {
	teams    teamports.TeamRepository
	projects projectports.ProjectRepository
}

func NovoFecharProjectUseCase(teams teamports.TeamRepository, projects projectports.ProjectRepository) *FecharProjectUseCase {
	return &FecharProjectUseCase{teams: teams, projects: projects}
}

func (uc *FecharProjectUseCase) Executar(req dto.FecharProjectRequest) error {
	p, err := uc.projects.BuscarPorID(project.ProjectID(req.ProjectID))
	if err != nil {
		return err
	}

	tm, err := uc.teams.BuscarPorID(p.TeamID())
	if err != nil {
		return err
	}

	leaderID := team.UserID(req.LeaderID)
	if !tm.EhLeader(leaderID) {
		return ErrSomenteLeaderPodeGerenciarProject
	}

	p.Fechar(time.Now().UTC())
	return uc.projects.Salvar(p)
}
