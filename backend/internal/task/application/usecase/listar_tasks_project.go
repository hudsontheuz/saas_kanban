package usecase

import (
	"strings"

	projectports "github.com/hudsontheuz/saas_kanban/internal/project/application/ports"
	project "github.com/hudsontheuz/saas_kanban/internal/project/domain"
	taskdto "github.com/hudsontheuz/saas_kanban/internal/task/application/dto"
	taskports "github.com/hudsontheuz/saas_kanban/internal/task/application/ports"
	task "github.com/hudsontheuz/saas_kanban/internal/task/domain"
)

type ListarTasksPorProjectUseCase struct {
	projects projectports.ProjectRepository
	tasks    taskports.TaskRepository
}

func NovoListarTasksPorProjectUseCase(
	projects projectports.ProjectRepository,
	tasks taskports.TaskRepository,
) *ListarTasksPorProjectUseCase {
	return &ListarTasksPorProjectUseCase{
		projects: projects,
		tasks:    tasks,
	}
}

func (uc *ListarTasksPorProjectUseCase) Executar(
	req taskdto.ListarTasksPorProjectRequest,
) (taskdto.ListarTasksPorProjectResponse, error) {
	_, err := uc.projects.BuscarPorID(project.ProjectID(req.ProjectID))
	if err != nil {
		return taskdto.ListarTasksPorProjectResponse{}, err
	}

	var filtro *task.StatusTask
	if req.Status != nil && strings.TrimSpace(*req.Status) != "" {
		status := task.StatusTask(strings.ToUpper(strings.TrimSpace(*req.Status)))

		switch status {
		case task.ToDo, task.Doing, task.InReview, task.Done:
			filtro = &status
		default:
			return taskdto.ListarTasksPorProjectResponse{}, task.ErrTransicaoInvalida
		}
	}

	lista, err := uc.tasks.ListarPorProjectID(project.ProjectID(req.ProjectID), filtro)
	if err != nil {
		return taskdto.ListarTasksPorProjectResponse{}, err
	}

	items := make([]taskdto.TaskListItem, 0, len(lista))
	for _, tk := range lista {
		item := taskdto.TaskListItem{
			TaskID:            string(tk.ID()),
			ProjectID:         string(tk.ProjectID()),
			Titulo:            tk.Titulo(),
			Descricao:         tk.Descricao(),
			ComentarioEntrega: tk.ComentarioEntrega(),
			ComentarioReview:  tk.ComentarioReview(),
			Status:            string(tk.Status()),
			Paused:            tk.IsPaused(),
		}

		if tk.Assignee() != nil {
			assignee := string(*tk.Assignee())
			item.AssigneeID = &assignee
		}

		if tk.Outcome() != nil {
			outcome := string(*tk.Outcome())
			item.Outcome = &outcome
		}

		items = append(items, item)
	}

	return taskdto.ListarTasksPorProjectResponse{Items: items}, nil
}
