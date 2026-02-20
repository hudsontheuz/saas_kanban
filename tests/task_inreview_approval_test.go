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

func TestMoverParaInReview_AssigneeOnly(t *testing.T) {
	projectRepo := memory.NovoProjectRepoEmMemoria()
	taskRepo := memory.NovoTaskRepoEmMemoria()

	// team + project
	leader := team.UserID("leader-1")
	tm, _ := team.NovaTeam("T", leader)
	teamRepo := memory.NovoTeamRepoEmMemoria()
	_ = teamRepo.Salvar(tm)

	p, _ := project.NovoProject(tm.ID(), project.ConfiguracoesProject{})
	_ = projectRepo.Salvar(p)

	// task + selfassign
	tk, _ := task.NovaTask(p.ID(), "Task")
	_ = taskRepo.Salvar(tk)

	ucAssign := usecase.NovoSelfAssignTaskUseCase(projectRepo, taskRepo)
	_ = ucAssign.Executar(dto.SelfAssignRequest{TaskID: string(tk.ID()), UserID: "user-1"})

	ucInReview := usecase.NovoMoverParaInReviewUseCase(projectRepo, taskRepo)

	// outro usuário tenta -> erro
	err := ucInReview.Executar(dto.MoverParaInReviewRequest{TaskID: string(tk.ID()), UserID: "user-2"})
	if err == nil {
		t.Fatalf("esperava erro: somente assignee pode mover para InReview")
	}

	// assignee -> ok
	err = ucInReview.Executar(dto.MoverParaInReviewRequest{TaskID: string(tk.ID()), UserID: "user-1"})
	if err != nil {
		t.Fatalf("esperava mover para InReview ok, veio: %v", err)
	}
}

func TestAprovar_LeaderOnly(t *testing.T) {
	projectRepo := memory.NovoProjectRepoEmMemoria()
	taskRepo := memory.NovoTaskRepoEmMemoria()

	leader := team.UserID("leader-1")
	tm, _ := team.NovaTeam("T", leader)
	teamRepo := memory.NovoTeamRepoEmMemoria()
	_ = teamRepo.Salvar(tm)

	p, _ := project.NovoProject(tm.ID(), project.ConfiguracoesProject{})
	_ = projectRepo.Salvar(p)

	tk, _ := task.NovaTask(p.ID(), "Task")
	_ = taskRepo.Salvar(tk)

	ucAssign := usecase.NovoSelfAssignTaskUseCase(projectRepo, taskRepo)
	ucInReview := usecase.NovoMoverParaInReviewUseCase(projectRepo, taskRepo)

	_ = ucAssign.Executar(dto.SelfAssignRequest{TaskID: string(tk.ID()), UserID: "user-1"})
	_ = ucInReview.Executar(dto.MoverParaInReviewRequest{TaskID: string(tk.ID()), UserID: "user-1"})

	ucAprovar := usecase.NovoAprovarTaskUseCase(projectRepo, teamRepo, taskRepo)

	// não-leader tenta -> erro
	err := ucAprovar.Executar(dto.AprovarTaskRequest{TaskID: string(tk.ID()), LeaderID: "user-1"})
	if err == nil {
		t.Fatalf("esperava erro: somente leader aprova")
	}

	// leader -> ok
	err = ucAprovar.Executar(dto.AprovarTaskRequest{TaskID: string(tk.ID()), LeaderID: string(leader)})
	if err != nil {
		t.Fatalf("esperava aprovar ok, veio: %v", err)
	}
}

func TestReprovar_InReviewVoltaParaDoing(t *testing.T) {
	projectRepo := memory.NovoProjectRepoEmMemoria()
	taskRepo := memory.NovoTaskRepoEmMemoria()

	leader := team.UserID("leader-1")
	tm, _ := team.NovaTeam("T", leader)
	teamRepo := memory.NovoTeamRepoEmMemoria()
	_ = teamRepo.Salvar(tm)

	p, _ := project.NovoProject(tm.ID(), project.ConfiguracoesProject{})
	_ = projectRepo.Salvar(p)

	tk, _ := task.NovaTask(p.ID(), "Task")
	_ = taskRepo.Salvar(tk)

	ucAssign := usecase.NovoSelfAssignTaskUseCase(projectRepo, taskRepo)
	ucInReview := usecase.NovoMoverParaInReviewUseCase(projectRepo, taskRepo)
	_ = ucAssign.Executar(dto.SelfAssignRequest{TaskID: string(tk.ID()), UserID: "user-1"})
	_ = ucInReview.Executar(dto.MoverParaInReviewRequest{TaskID: string(tk.ID()), UserID: "user-1"})

	ucReprovar := usecase.NovoReprovarTaskUseCase(projectRepo, teamRepo, taskRepo)
	err := ucReprovar.Executar(dto.ReprovarTaskRequest{TaskID: string(tk.ID()), LeaderID: string(leader)})
	if err != nil {
		t.Fatalf("esperava reprovar ok, veio: %v", err)
	}
}

func TestRejeitarEmToDo_ToDoParaDoneRejected(t *testing.T) {
	projectRepo := memory.NovoProjectRepoEmMemoria()
	taskRepo := memory.NovoTaskRepoEmMemoria()

	leader := team.UserID("leader-1")
	tm, _ := team.NovaTeam("T", leader)
	teamRepo := memory.NovoTeamRepoEmMemoria()
	_ = teamRepo.Salvar(tm)

	p, _ := project.NovoProject(tm.ID(), project.ConfiguracoesProject{})
	_ = projectRepo.Salvar(p)

	tk, _ := task.NovaTask(p.ID(), "Task")
	_ = taskRepo.Salvar(tk)

	ucRejeitar := usecase.NovoRejeitarTaskToDoUseCase(projectRepo, teamRepo, taskRepo)
	err := ucRejeitar.Executar(dto.RejeitarTaskToDoRequest{TaskID: string(tk.ID()), LeaderID: string(leader)})
	if err != nil {
		t.Fatalf("esperava rejeitar ok, veio: %v", err)
	}
}
