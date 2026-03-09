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
	teammemory "github.com/hudsontheuz/saas_kanban/internal/team/infrastructure/persistence/memory"
	user "github.com/hudsontheuz/saas_kanban/internal/user/domain"
)

func TestMoverParaInReview_AssigneeOnly(t *testing.T) {
	projectRepo := projectmemory.NovoProjectRepoEmMemoria()
	taskRepo := taskmemory.NovoTaskRepoEmMemoria()

	leader := user.UserID("1")
	tm, err := team.NovaTeam("T", leader)
	if err != nil {
		t.Fatalf("erro ao criar team: %v", err)
	}
	teamRepo := teammemory.NovoTeamRepoEmMemoria()
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

	if tk.Status() != task.InReview {
		t.Fatalf("esperava task em InReview, veio %s", tk.Status())
	}
}

func TestReprovar_RetornaParaToDo(t *testing.T) {
	projectRepo := projectmemory.NovoProjectRepoEmMemoria()
	taskRepo := taskmemory.NovoTaskRepoEmMemoria()

	leader := user.UserID("1")
	tm, err := team.NovaTeam("T", leader)
	if err != nil {
		t.Fatalf("erro ao criar team: %v", err)
	}
	teamRepo := teammemory.NovoTeamRepoEmMemoria()
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
	ucReprovar := usecase.NovoReprovarTaskUseCase(projectRepo, teamRepo, taskRepo)

	_ = ucAssign.Executar(dto.SelfAssignRequest{TaskID: string(tk.ID()), UserID: "2"})
	_ = ucInReview.Executar(dto.MoverParaInReviewRequest{TaskID: string(tk.ID()), UserID: "2"})

	err = ucReprovar.Executar(dto.ReprovarTaskRequest{TaskID: string(tk.ID()), LeaderID: string(leader)})
	if err != nil {
		t.Fatalf("esperava reprovar ok, veio: %v", err)
	}

	if tk.Status() != task.ToDo {
		t.Fatalf("esperava task voltando para ToDo, veio %s", tk.Status())
	}
}

func TestAprovar_LeaderOnly(t *testing.T) {
	projectRepo := projectmemory.NovoProjectRepoEmMemoria()
	taskRepo := taskmemory.NovoTaskRepoEmMemoria()

	leader := user.UserID("1")
	tm, err := team.NovaTeam("T", leader)
	if err != nil {
		t.Fatalf("erro ao criar team: %v", err)
	}
	teamRepo := teammemory.NovoTeamRepoEmMemoria()
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
	ucApprove := usecase.NovoAprovarTaskUseCase(projectRepo, teamRepo, taskRepo)

	_ = ucAssign.Executar(dto.SelfAssignRequest{TaskID: string(tk.ID()), UserID: "2"})
	_ = ucInReview.Executar(dto.MoverParaInReviewRequest{TaskID: string(tk.ID()), UserID: "2"})

	err = ucApprove.Executar(dto.AprovarTaskRequest{TaskID: string(tk.ID()), LeaderID: "2"})
	if err == nil {
		t.Fatalf("esperava erro: somente leader aprova")
	}

	err = ucApprove.Executar(dto.AprovarTaskRequest{TaskID: string(tk.ID()), LeaderID: string(leader)})
	if err != nil {
		t.Fatalf("esperava aprovar ok, veio: %v", err)
	}
}
