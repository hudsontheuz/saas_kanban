package usecase

import (
	"github.com/hudsontheuz/saas_kanban/internal/application/task/dto"
	projectports "github.com/hudsontheuz/saas_kanban/internal/application/project/ports"
	taskports "github.com/hudsontheuz/saas_kanban/internal/application/task/ports"
	"github.com/hudsontheuz/saas_kanban/internal/domain/project"
	"github.com/hudsontheuz/saas_kanban/internal/domain/task"
	"github.com/hudsontheuz/saas_kanban/internal/domain/team"
)

type PausarTaskUseCase struct {
	projects projectports.ProjectRepository
	tasks    taskports.TaskRepository
}

func NovoPausarTaskUseCase(projects projectports.ProjectRepository, tasks taskports.TaskRepository) *PausarTaskUseCase {
	return &PausarTaskUseCase{projects: projects, tasks: tasks}
}

func (uc *PausarTaskUseCase) Executar(req taskdto.PausarTaskRequest) error {
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

	if err := tk.PodePausarOuRetomar(userID); err != nil {
		return err
	}

	if err := tk.Pausar(); err != nil {
		return err
	}
	if err := tk.ValidarInvariantes(); err != nil {
		return err
	}

	return uc.tasks.Salvar(tk)
}
