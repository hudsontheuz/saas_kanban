package usecase

import (
	"errors"

	projectports "github.com/hudsontheuz/saas_kanban/internal/project/application/ports"
	"github.com/hudsontheuz/saas_kanban/internal/project/domain"
	"github.com/hudsontheuz/saas_kanban/internal/task/application/dto"
	taskports "github.com/hudsontheuz/saas_kanban/internal/task/application/ports"
	"github.com/hudsontheuz/saas_kanban/internal/task/domain"
	"github.com/hudsontheuz/saas_kanban/internal/team/domain"
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
