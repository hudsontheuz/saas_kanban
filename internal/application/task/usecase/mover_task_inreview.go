package usecase

import (
	"github.com/hudsontheuz/saas_kanban/internal/application/task/dto"
	projectports "github.com/hudsontheuz/saas_kanban/internal/application/project/ports"
	taskports "github.com/hudsontheuz/saas_kanban/internal/application/task/ports"
	"github.com/hudsontheuz/saas_kanban/internal/domain/project"
	"github.com/hudsontheuz/saas_kanban/internal/domain/task"
	"github.com/hudsontheuz/saas_kanban/internal/domain/team"
)

type MoverParaInReviewUseCase struct {
	projects projectports.ProjectRepository
	tasks    taskports.TaskRepository
}

func NovoMoverParaInReviewUseCase(projects projectports.ProjectRepository, tasks taskports.TaskRepository) *MoverParaInReviewUseCase {
	return &MoverParaInReviewUseCase{projects: projects, tasks: tasks}
}

func (uc *MoverParaInReviewUseCase) Executar(req taskdto.MoverParaInReviewRequest) error {
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

	if err := tk.PodeMoverParaInReview(userID); err != nil {
		return err
	}

	if err := tk.MoverParaInReview(); err != nil {
		return err
	}
	if err := tk.ValidarInvariantes(); err != nil {
		return err
	}

	return uc.tasks.Salvar(tk)
}