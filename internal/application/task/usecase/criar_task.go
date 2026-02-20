package usecase

import (
	"github.com/hudsontheuz/saas_kanban/internal/application/task/dto"
	projectports "github.com/hudsontheuz/saas_kanban/internal/application/project/ports"
	taskports "github.com/hudsontheuz/saas_kanban/internal/application/task/ports"
	"github.com/hudsontheuz/saas_kanban/internal/domain/project"
	"github.com/hudsontheuz/saas_kanban/internal/domain/task"
)

type CriarTaskUseCase struct {
	projects projectports.ProjectRepository
	tasks    taskports.TaskRepository
}

func NovoCriarTaskUseCase(projects projectports.ProjectRepository, tasks taskports.TaskRepository) *CriarTaskUseCase {
	return &CriarTaskUseCase{projects: projects, tasks: tasks}
}

func (uc *CriarTaskUseCase) Executar(req taskdto.CriarTaskRequest) (taskdto.CriarTaskResponse, error) {
	p, err := uc.projects.BuscarPorID(project.ProjectID(req.ProjectID))
	if err != nil {
		return taskdto.CriarTaskResponse{}, err
	}
	if p.EstaFechado() {
		return taskdto.CriarTaskResponse{}, project.ErrProjetoFechado
	}

	tk, err := task.NovaTask(p.ID(), req.Titulo)
	if err != nil {
		return taskdto.CriarTaskResponse{}, err
	}

	if err := uc.tasks.Salvar(tk); err != nil {
		return taskdto.CriarTaskResponse{}, err
	}

	return taskdto.CriarTaskResponse{TaskID: string(tk.ID())}, nil
}
