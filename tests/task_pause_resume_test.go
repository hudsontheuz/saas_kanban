package tests

import (
	"testing"

	"github.com/hudsontheuz/saas_kanban/internal/application/task/dto"
	usecase "github.com/hudsontheuz/saas_kanban/internal/application/task/usecase"
	"github.com/hudsontheuz/saas_kanban/internal/domain/project"
	"github.com/hudsontheuz/saas_kanban/internal/domain/task"
	"github.com/hudsontheuz/saas_kanban/internal/domain/team"
	"github.com/hudsontheuz/saas_kanban/internal/infrastructure/persistence/memory"
)

func TestPause_NaoContaNoLimiteGlobal_SelfAssignEmOutraTaskFunciona(t *testing.T) {
	projectRepo := memory.NovoProjectRepoEmMemoria()
	taskRepo := memory.NovoTaskRepoEmMemoria()

	p, _ := project.NovoProject(team.TeamID("team-1"), project.ConfiguracoesProject{})
	_ = projectRepo.Salvar(p)

	// cria duas tasks
	t1, _ := task.NovaTask(p.ID(), "T1")
	t2, _ := task.NovaTask(p.ID(), "T2")
	_ = taskRepo.Salvar(t1)
	_ = taskRepo.Salvar(t2)

	user := "user-1"

	// self-assign na primeira (vira DOING ativa)
	ucAssign := usecase.NovoSelfAssignTaskUseCase(projectRepo, taskRepo)
	if err := ucAssign.Executar(taskdto.SelfAssignRequest{TaskID: string(t1.ID()), UserID: user}); err != nil {
		t.Fatalf("esperava self-assign em t1 ok, veio: %v", err)
	}

	// pausar t1 (continua DOING, mas pausada)
	ucPause := usecase.NovoPausarTaskUseCase(projectRepo, taskRepo)
	if err := ucPause.Executar(taskdto.PausarTaskRequest{TaskID: string(t1.ID()), UserID: user}); err != nil {
		t.Fatalf("esperava pausar ok, veio: %v", err)
	}

	// agora self-assign na segunda deve funcionar (porque t1 pausada não conta)
	if err := ucAssign.Executar(taskdto.SelfAssignRequest{TaskID: string(t2.ID()), UserID: user}); err != nil {
		t.Fatalf("esperava self-assign em t2 ok (t1 pausada), veio: %v", err)
	}
}

func TestRetomar_FalhaSeExisteOutraDoingAtiva(t *testing.T) {
	projectRepo := memory.NovoProjectRepoEmMemoria()
	taskRepo := memory.NovoTaskRepoEmMemoria()

	p, _ := project.NovoProject(team.TeamID("team-1"), project.ConfiguracoesProject{})
	_ = projectRepo.Salvar(p)

	t1, _ := task.NovaTask(p.ID(), "T1")
	t2, _ := task.NovaTask(p.ID(), "T2")
	_ = taskRepo.Salvar(t1)
	_ = taskRepo.Salvar(t2)

	user := "user-1"

	ucAssign := usecase.NovoSelfAssignTaskUseCase(projectRepo, taskRepo)
	ucPause := usecase.NovoPausarTaskUseCase(projectRepo, taskRepo)
	ucResume := usecase.NovoRetomarTaskUseCase(projectRepo, taskRepo)

	// self-assign t1 e pausa
	_ = ucAssign.Executar(taskdto.SelfAssignRequest{TaskID: string(t1.ID()), UserID: user})
	_ = ucPause.Executar(taskdto.PausarTaskRequest{TaskID: string(t1.ID()), UserID: user})

	// self-assign t2 (fica DOING ativa)
	_ = ucAssign.Executar(taskdto.SelfAssignRequest{TaskID: string(t2.ID()), UserID: user})

	// tentar retomar t1 deve falhar porque existe outra DOING ativa (t2)
	err := ucResume.Executar(taskdto.RetomarTaskRequest{TaskID: string(t1.ID()), UserID: user})
	if err == nil {
		t.Fatalf("esperava erro ao retomar com outra DOING ativa")
	}
}
