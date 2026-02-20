package usecase

import (
	"github.com/hudsontheuz/saas_kanban/internal/application/task/dto"
	projectports "github.com/hudsontheuz/saas_kanban/internal/application/project/ports"
	taskports "github.com/hudsontheuz/saas_kanban/internal/application/task/ports"
	teamports "github.com/hudsontheuz/saas_kanban/internal/application/team/ports"
	"github.com/hudsontheuz/saas_kanban/internal/domain/project"
	"github.com/hudsontheuz/saas_kanban/internal/domain/task"
	"github.com/hudsontheuz/saas_kanban/internal/domain/team"
)

type RejeitarTaskToDoUseCase struct {
	projects projectports.ProjectRepository
	teams    teamports.TeamRepository
	tasks    taskports.TaskRepository
}

func NovoRejeitarTaskToDoUseCase(projects projectports.ProjectRepository, teams teamports.TeamRepository, tasks taskports.TaskRepository) *RejeitarTaskToDoUseCase {
	return &RejeitarTaskToDoUseCase{projects: projects, teams: teams, tasks: tasks}
}

func (uc *RejeitarTaskToDoUseCase) Executar(req dto.RejeitarTaskToDoRequest) error {
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

	if err := tk.RejeitarEmToDo(); err != nil {
		return err
	}
	if err := tk.ValidarInvariantes(); err != nil {
		return err
	}

	return uc.tasks.Salvar(tk)
}