package usecase

import (
	projectports "github.com/hudsontheuz/saas_kanban/internal/project/application/ports"
	"github.com/hudsontheuz/saas_kanban/internal/project/domain"
	"github.com/hudsontheuz/saas_kanban/internal/task/application/dto"
	taskports "github.com/hudsontheuz/saas_kanban/internal/task/application/ports"
	"github.com/hudsontheuz/saas_kanban/internal/task/domain"
	"github.com/hudsontheuz/saas_kanban/internal/team/domain"
)

type RetomarTaskUseCase struct {
	projects projectports.ProjectRepository
	tasks    taskports.TaskRepository
}

func NovoRetomarTaskUseCase(projects projectports.ProjectRepository, tasks taskports.TaskRepository) *RetomarTaskUseCase {
	return &RetomarTaskUseCase{projects: projects, tasks: tasks}
}

func (uc *RetomarTaskUseCase) Executar(req dto.RetomarTaskRequest) error {
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

	// Regra global: só pode retomar se não existir outra DOING ativa
	existe, err := uc.tasks.ExisteDoingAtivaParaUser(userID)
	if err != nil {
		return err
	}
	// cuidado: a própria task atual está em DOING, mas se ela está pausada, ela não conta.
	// Nosso repo em memória já ignora pausadas, então aqui "existe" só será true se tiver OUTRA DOING ativa.
	if existe {
		return ErrLimiteGlobalDoing
	}

	if err := tk.Retomar(); err != nil {
		return err
	}
	if err := tk.ValidarInvariantes(); err != nil {
		return err
	}

	return uc.tasks.Salvar(tk)
}
