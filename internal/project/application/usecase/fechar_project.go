package usecase

import (
	"time"

	"github.com/hudsontheuz/saas_kanban/internal/project/application/dto"
	projectports "github.com/hudsontheuz/saas_kanban/internal/project/application/ports"
	"github.com/hudsontheuz/saas_kanban/internal/project/domain"
	teamports "github.com/hudsontheuz/saas_kanban/internal/team/application/ports"
	"github.com/hudsontheuz/saas_kanban/internal/user/domain"
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

	leaderID := user.UserID(req.LeaderID)
	if !tm.EhLeader(leaderID) {
		return ErrSomenteLeaderPodeGerenciarProject
	}

	p.Fechar(time.Now().UTC())
	return uc.projects.Salvar(p)
}
