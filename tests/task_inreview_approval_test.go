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

func TestMoverParaInReview_AssigneeOnly(t *testing.T) {
	projectRepo := memory.NovoProjectRepoEmMemoria()
	taskRepo := memory.NovoTaskRepoEmMemoria()

	leader := team.UserID("1")
	tm, err := team.NovaTeam("T", leader)
	if err != nil {
		t.Fatalf("erro ao criar team: %v", err)
	}
	teamRepo := memory.NovoTeamRepoEmMemoria()
	_ = teamRepo.Salvar(tm)

	p, err := project.NovoProject(tm.ID(), "Projeto Teste", project.ConfiguracoesProject{})
	if err != nil {
		t.Fatalf("erro ao criar project: %v", err)
	}
	_ = projectRepo.Salvar(p)

	tk, err := task.NovaTask(p.ID(), "Task")
	if err != nil {
		t.Fatalf("erro ao criar task: %v", err)
	}
	_ = taskRepo.Salvar(tk)

	ucAssign := usecase.NovoSelfAssignTaskUseCase(projectRepo, taskRepo)
	_ = ucAssign.Executar(dto.SelfAssignRequest{TaskID: string(tk.ID()), UserID: "2"})

	ucInReview := usecase.NovoMoverParaInReviewUseCase(projectRepo, taskRepo)

	err = ucInReview.Executar(dto.MoverParaInReviewRequest{TaskID: string(tk.ID()), UserID: "3"})
	if err == nil {
		t.Fatalf("esperava erro: somente assignee pode mover para InReview")
	}

	err = ucInReview.Executar(dto.MoverParaInReviewRequest{TaskID: string(tk.ID()), UserID: "2"})
	if err != nil {
		t.Fatalf("esperava mover para InReview ok, veio: %v", err)
	}
}

func TestAprovar_LeaderOnly(t *testing.T) {
	projectRepo := memory.NovoProjectRepoEmMemoria()
	taskRepo := memory.NovoTaskRepoEmMemoria()

	leader := team.UserID("1")
	tm, err := team.NovaTeam("T", leader)
	if err != nil {
		t.Fatalf("erro ao criar team: %v", err)
	}
	teamRepo := memory.NovoTeamRepoEmMemoria()
	_ = teamRepo.Salvar(tm)

	p, err := project.NovoProject(tm.ID(), "Projeto Teste", project.ConfiguracoesProject{})
	if err != nil {
		t.Fatalf("erro ao criar project: %v", err)
	}
	_ = projectRepo.Salvar(p)

	tk, err := task.NovaTask(p.ID(), "Task")
	if err != nil {
		t.Fatalf("erro ao criar task: %v", err)
	}
	_ = taskRepo.Salvar(tk)

	ucAssign := usecase.NovoSelfAssignTaskUseCase(projectRepo, taskRepo)
	ucInReview := usecase.NovoMoverParaInReviewUseCase(projectRepo, taskRepo)

	_ = ucAssign.Executar(dto.SelfAssignRequest{TaskID: string(tk.ID()), UserID: "2"})
	_ = ucInReview.Executar(dto.MoverParaInReviewRequest{TaskID: string(tk.ID()), UserID: "2"})

	ucAprovar := usecase.NovoAprovarTaskUseCase(projectRepo, teamRepo, taskRepo)

	err = ucAprovar.Executar(dto.AprovarTaskRequest{TaskID: string(tk.ID()), LeaderID: "2"})
	if err == nil {
		t.Fatalf("esperava erro: somente leader aprova")
	}

	err = ucAprovar.Executar(dto.AprovarTaskRequest{TaskID: string(tk.ID()), LeaderID: string(leader)})
	if err != nil {
		t.Fatalf("esperava aprovar ok, veio: %v", err)
	}
}

func TestReprovar_InReviewVoltaParaDoing(t *testing.T) {
	projectRepo := memory.NovoProjectRepoEmMemoria()
	taskRepo := memory.NovoTaskRepoEmMemoria()

	leader := team.UserID("1")
	tm, err := team.NovaTeam("T", leader)
	if err != nil {
		t.Fatalf("erro ao criar team: %v", err)
	}
	teamRepo := memory.NovoTeamRepoEmMemoria()
	_ = teamRepo.Salvar(tm)

	p, err := project.NovoProject(tm.ID(), "Projeto Teste", project.ConfiguracoesProject{})
	if err != nil {
		t.Fatalf("erro ao criar project: %v", err)
	}
	_ = projectRepo.Salvar(p)

	tk, err := task.NovaTask(p.ID(), "Task")
	if err != nil {
		t.Fatalf("erro ao criar task: %v", err)
	}
	_ = taskRepo.Salvar(tk)

	ucAssign := usecase.NovoSelfAssignTaskUseCase(projectRepo, taskRepo)
	ucInReview := usecase.NovoMoverParaInReviewUseCase(projectRepo, taskRepo)
	_ = ucAssign.Executar(dto.SelfAssignRequest{TaskID: string(tk.ID()), UserID: "2"})
	_ = ucInReview.Executar(dto.MoverParaInReviewRequest{TaskID: string(tk.ID()), UserID: "2"})

	ucReprovar := usecase.NovoReprovarTaskUseCase(projectRepo, teamRepo, taskRepo)
	err = ucReprovar.Executar(dto.ReprovarTaskRequest{TaskID: string(tk.ID()), LeaderID: string(leader)})
	if err != nil {
		t.Fatalf("esperava reprovar ok, veio: %v", err)
	}
}

func TestRejeitarEmToDo_ToDoParaDoneRejected(t *testing.T) {
	projectRepo := memory.NovoProjectRepoEmMemoria()
	taskRepo := memory.NovoTaskRepoEmMemoria()

	leader := team.UserID("1")
	tm, err := team.NovaTeam("T", leader)
	if err != nil {
		t.Fatalf("erro ao criar team: %v", err)
	}
	teamRepo := memory.NovoTeamRepoEmMemoria()
	_ = teamRepo.Salvar(tm)

	p, err := project.NovoProject(tm.ID(), "Projeto Teste", project.ConfiguracoesProject{})
	if err != nil {
		t.Fatalf("erro ao criar project: %v", err)
	}
	_ = projectRepo.Salvar(p)

	tk, err := task.NovaTask(p.ID(), "Task")
	if err != nil {
		t.Fatalf("erro ao criar task: %v", err)
	}
	_ = taskRepo.Salvar(tk)

	ucRejeitar := usecase.NovoRejeitarTaskToDoUseCase(projectRepo, teamRepo, taskRepo)
	err = ucRejeitar.Executar(dto.RejeitarTaskToDoRequest{TaskID: string(tk.ID()), LeaderID: string(leader)})
	if err != nil {
		t.Fatalf("esperava rejeitar ok, veio: %v", err)
	}
}
