package tests

import (
	"testing"

	"github.com/hudsontheuz/saas_kanban/infrastructure/persistence/memory"
	"github.com/hudsontheuz/saas_kanban/internal/application/task/dto"
	usecase "github.com/hudsontheuz/saas_kanban/internal/application/task/usecase"
	"github.com/hudsontheuz/saas_kanban/internal/domain/project"
	"github.com/hudsontheuz/saas_kanban/internal/domain/task"
	"github.com/hudsontheuz/saas_kanban/internal/domain/team"
)

func TestPause_NaoContaNoLimiteGlobal_SelfAssignEmOutraTaskFunciona(t *testing.T) {
	projectRepo := memory.NovoProjectRepoEmMemoria()
	taskRepo := memory.NovoTaskRepoEmMemoria()

	p, err := project.NovoProject(team.TeamID("1"), "Projeto Teste", project.ConfiguracoesProject{})
	if err != nil {
		t.Fatalf("erro ao criar project: %v", err)
	}
	_ = projectRepo.Salvar(p)

	t1, err := task.NovaTask(p.ID(), "T1")
	if err != nil {
		t.Fatalf("erro ao criar task t1: %v", err)
	}
	t2, err := task.NovaTask(p.ID(), "T2")
	if err != nil {
		t.Fatalf("erro ao criar task t2: %v", err)
	}
	_ = taskRepo.Salvar(t1)
	_ = taskRepo.Salvar(t2)

	user := "2"

	ucAssign := usecase.NovoSelfAssignTaskUseCase(projectRepo, taskRepo)
	if err := ucAssign.Executar(dto.SelfAssignRequest{TaskID: string(t1.ID()), UserID: user}); err != nil {
		t.Fatalf("esperava self-assign em t1 ok, veio: %v", err)
	}

	ucPause := usecase.NovoPausarTaskUseCase(projectRepo, taskRepo)
	if err := ucPause.Executar(dto.PausarTaskRequest{TaskID: string(t1.ID()), UserID: user}); err != nil {
		t.Fatalf("esperava pausar ok, veio: %v", err)
	}

	if err := ucAssign.Executar(dto.SelfAssignRequest{TaskID: string(t2.ID()), UserID: user}); err != nil {
		t.Fatalf("esperava self-assign em t2 ok (t1 pausada), veio: %v", err)
	}
}

func TestRetomar_FalhaSeExisteOutraDoingAtiva(t *testing.T) {
	projectRepo := memory.NovoProjectRepoEmMemoria()
	taskRepo := memory.NovoTaskRepoEmMemoria()

	p, err := project.NovoProject(team.TeamID("1"), "Projeto Teste", project.ConfiguracoesProject{})
	if err != nil {
		t.Fatalf("erro ao criar project: %v", err)
	}
	_ = projectRepo.Salvar(p)

	t1, err := task.NovaTask(p.ID(), "T1")
	if err != nil {
		t.Fatalf("erro ao criar task t1: %v", err)
	}
	t2, err := task.NovaTask(p.ID(), "T2")
	if err != nil {
		t.Fatalf("erro ao criar task t2: %v", err)
	}
	_ = taskRepo.Salvar(t1)
	_ = taskRepo.Salvar(t2)

	user := "2"

	ucAssign := usecase.NovoSelfAssignTaskUseCase(projectRepo, taskRepo)
	ucPause := usecase.NovoPausarTaskUseCase(projectRepo, taskRepo)
	ucResume := usecase.NovoRetomarTaskUseCase(projectRepo, taskRepo)

	_ = ucAssign.Executar(dto.SelfAssignRequest{TaskID: string(t1.ID()), UserID: user})
	_ = ucPause.Executar(dto.PausarTaskRequest{TaskID: string(t1.ID()), UserID: user})

	_ = ucAssign.Executar(dto.SelfAssignRequest{TaskID: string(t2.ID()), UserID: user})

	err = ucResume.Executar(dto.RetomarTaskRequest{TaskID: string(t1.ID()), UserID: user})
	if err == nil {
		t.Fatalf("esperava erro ao retomar com outra DOING ativa")
	}
}
