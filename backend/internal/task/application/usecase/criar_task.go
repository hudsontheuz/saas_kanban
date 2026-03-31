package usecase

import (
	projectports "github.com/hudsontheuz/saas_kanban/internal/project/application/ports"
	"github.com/hudsontheuz/saas_kanban/internal/project/domain"
	"github.com/hudsontheuz/saas_kanban/internal/task/application/dto"
	taskports "github.com/hudsontheuz/saas_kanban/internal/task/application/ports"
	"github.com/hudsontheuz/saas_kanban/internal/task/domain"
)

type CriarTaskUseCase struct {
	projects projectports.ProjectRepository
	tasks    taskports.TaskRepository
}

func NovoCriarTaskUseCase(projects projectports.ProjectRepository, tasks taskports.TaskRepository) *CriarTaskUseCase {
	return &CriarTaskUseCase{projects: projects, tasks: tasks}
}

func (uc *CriarTaskUseCase) Executar(req dto.CriarTaskRequest) (dto.CriarTaskResponse, error) {
	p, err := uc.projects.BuscarPorID(project.ProjectID(req.ProjectID))
	if err != nil {
		return dto.CriarTaskResponse{}, err
	}
	if p.EstaFechado() {
		return dto.CriarTaskResponse{}, project.ErrProjetoFechado
	}

	tk, err := task.NovaTask(p.ID(), req.Titulo, req.Descricao)
	if err != nil {
		return dto.CriarTaskResponse{}, err
	}

	if err := uc.tasks.Salvar(tk); err != nil {
		return dto.CriarTaskResponse{}, err
	}

	return dto.CriarTaskResponse{TaskID: string(tk.ID())}, nil
}
