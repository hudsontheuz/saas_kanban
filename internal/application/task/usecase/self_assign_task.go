package usecase

import (
	"errors"

	"github.com/hudsontheuz/saas_kanban/internal/application/task/dto"
	projectports "github.com/hudsontheuz/saas_kanban/internal/application/project/ports"
	taskports "github.com/hudsontheuz/saas_kanban/internal/application/task/ports"
	"github.com/hudsontheuz/saas_kanban/internal/domain/project"
	"github.com/hudsontheuz/saas_kanban/internal/domain/task"
	"github.com/hudsontheuz/saas_kanban/internal/domain/team"
)

var ErrLimiteGlobalDoing = errors.New("limite global: usuário já possui task DOING ativa")

type SelfAssignTaskUseCase struct {
	projects projectports.ProjectRepository
	tasks    taskports.TaskRepository
}

func NovoSelfAssignTaskUseCase(projects projectports.ProjectRepository, tasks taskports.TaskRepository) *SelfAssignTaskUseCase {
	return &SelfAssignTaskUseCase{projects: projects, tasks: tasks}
}

func (uc *SelfAssignTaskUseCase) Executar(req dto.SelfAssignRequest) error {
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

	userID := team.UserID(req.UserID)

	existe, err := uc.tasks.ExisteDoingAtivaParaUser(userID)
	if err != nil {
		return err
	}
	if existe {
		return ErrLimiteGlobalDoing
	}

	if err := tk.SelfAssign(userID); err != nil {
		return err
	}
	if err := tk.ValidarInvariantes(); err != nil {
		return err
	}

	return uc.tasks.Salvar(tk)
}
