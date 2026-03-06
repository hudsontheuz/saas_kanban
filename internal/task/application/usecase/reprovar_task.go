package usecase

import (
	projectports "github.com/hudsontheuz/saas_kanban/internal/project/application/ports"
	"github.com/hudsontheuz/saas_kanban/internal/project/domain"
	"github.com/hudsontheuz/saas_kanban/internal/task/application/dto"
	taskports "github.com/hudsontheuz/saas_kanban/internal/task/application/ports"
	"github.com/hudsontheuz/saas_kanban/internal/task/domain"
	teamports "github.com/hudsontheuz/saas_kanban/internal/team/application/ports"
	"github.com/hudsontheuz/saas_kanban/internal/team/domain"
)

type ReprovarTaskUseCase struct {
	projects projectports.ProjectRepository
	teams    teamports.TeamRepository
	tasks    taskports.TaskRepository
}

func NovoReprovarTaskUseCase(projects projectports.ProjectRepository, teams teamports.TeamRepository, tasks taskports.TaskRepository) *ReprovarTaskUseCase {
	return &ReprovarTaskUseCase{projects: projects, teams: teams, tasks: tasks}
}

func (uc *ReprovarTaskUseCase) Executar(req dto.ReprovarTaskRequest) error {
	tk, err := uc.tasks.BuscarPorID(task.TaskID(req.TaskID))
	if err != nil {
		return err
	}

	p, err := uc.projects.BuscarPorID(tk.ProjectID())
	if err != nil {
		return err
	}
	if p.EstaFechado() {
		return project.ErrProjetoFechado
	}

	tm, err := uc.teams.BuscarPorID(p.TeamID())
	if err != nil {
		return err
	}

	leaderID := team.UserID(req.LeaderID)
	if !tm.EhLeader(leaderID) {
		return ErrSomenteLeaderPodeDecidir
	}

	if err := tk.ReprovarParaAjustes(); err != nil {
		return err
	}
	if err := tk.ValidarInvariantes(); err != nil {
		return err
	}

	return uc.tasks.Salvar(tk)
}
