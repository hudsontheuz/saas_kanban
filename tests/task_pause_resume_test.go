package tests

import (
	"testing"

	project "github.com/hudsontheuz/saas_kanban/internal/project/domain"
	projectmemory "github.com/hudsontheuz/saas_kanban/internal/project/infrastructure/persistence/memory"
	dto "github.com/hudsontheuz/saas_kanban/internal/task/application/dto"
	usecase "github.com/hudsontheuz/saas_kanban/internal/task/application/usecase"
	task "github.com/hudsontheuz/saas_kanban/internal/task/domain"
	taskmemory "github.com/hudsontheuz/saas_kanban/internal/task/infrastructure/persistence/memory"
	team "github.com/hudsontheuz/saas_kanban/internal/team/domain"
)

func TestPause_NaoContaNoLimiteGlobal_SelfAssignEmOutraTaskFunciona(t *testing.T) {
	projectRepo := projectmemory.NovoProjectRepoEmMemoria()
	taskRepo := taskmemory.NovoTaskRepoEmMemoria()

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
	projectRepo := projectmemory.NovoProjectRepoEmMemoria()
	taskRepo := taskmemory.NovoTaskRepoEmMemoria()

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

	// tentar retomar t1 deve falhar porque t2 está DOING ativa
	if err := ucResume.Executar(dto.RetomarTaskRequest{TaskID: string(t1.ID()), UserID: user}); err == nil {
		t.Fatalf("esperava erro ao retomar: existe outra DOING ativa")
	}
}
