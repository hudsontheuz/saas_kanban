package usecase

import (
	projectports "github.com/hudsontheuz/saas_kanban/internal/project/application/ports"
	"github.com/hudsontheuz/saas_kanban/internal/project/domain"
	"github.com/hudsontheuz/saas_kanban/internal/task/application/dto"
	taskports "github.com/hudsontheuz/saas_kanban/internal/task/application/ports"
	"github.com/hudsontheuz/saas_kanban/internal/task/domain"
	"github.com/hudsontheuz/saas_kanban/internal/user/domain"
)

type MoverParaInReviewUseCase struct {
	projects projectports.ProjectRepository
	tasks    taskports.TaskRepository
}

func NovoMoverParaInReviewUseCase(projects projectports.ProjectRepository, tasks taskports.TaskRepository) *MoverParaInReviewUseCase {
	return &MoverParaInReviewUseCase{projects: projects, tasks: tasks}
}

func (uc *MoverParaInReviewUseCase) Executar(req dto.MoverParaInReviewRequest) error {
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

	userID := user.UserID(req.UserID)

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
