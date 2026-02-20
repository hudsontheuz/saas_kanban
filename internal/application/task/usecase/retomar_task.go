package usecase

import (
	projectports "github.com/hudsontheuz/saas_kanban/internal/application/project/ports"
	"github.com/hudsontheuz/saas_kanban/internal/application/task/dto"
	taskports "github.com/hudsontheuz/saas_kanban/internal/application/task/ports"
	"github.com/hudsontheuz/saas_kanban/internal/domain/project"
	"github.com/hudsontheuz/saas_kanban/internal/domain/task"
	"github.com/hudsontheuz/saas_kanban/internal/domain/team"
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
