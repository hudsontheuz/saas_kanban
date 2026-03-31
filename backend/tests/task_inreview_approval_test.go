package tests

import (
	"testing"

	project "github.com/hudsontheuz/saas_kanban/internal/project/domain"
	projectrepo "github.com/hudsontheuz/saas_kanban/internal/project/infrastructure/persistence/gorm/repo"
	dto "github.com/hudsontheuz/saas_kanban/internal/task/application/dto"
	usecase "github.com/hudsontheuz/saas_kanban/internal/task/application/usecase"
	task "github.com/hudsontheuz/saas_kanban/internal/task/domain"
	taskrepo "github.com/hudsontheuz/saas_kanban/internal/task/infrastructure/persistence/gorm/repo"
	teamrepo "github.com/hudsontheuz/saas_kanban/internal/team/infrastructure/persistence/gorm/repo"
)

const comentarioEntregaTeste = "Implementei a funcionalidade e validei o fluxo principal."

func TestMoverParaInReview_AssigneeOnly(t *testing.T) {
	db := openTestDB(t)

	projectRepo := projectrepo.NewProjectRepo(db)
	taskRepo := taskrepo.NewTaskRepo(db)

	leaderID := seedUser(t, db, "Leader Teste")
	assigneeID := seedUser(t, db, "Assignee Teste")
	outroUserID := seedUser(t, db, "Outro User")
	teamID := seedTeam(t, db, "Team Teste", leaderID)
	projectID := seedProject(t, db, teamID, "Projeto Teste", project.ConfiguracoesProject{})
	taskID := seedTask(t, db, projectID, "Task")

	ucAssign := usecase.NovoSelfAssignTaskUseCase(projectRepo, taskRepo)
	if err := ucAssign.Executar(dto.SelfAssignRequest{
		TaskID: string(taskID),
		UserID: string(assigneeID),
	}); err != nil {
		t.Fatalf("erro no self-assign: %v", err)
	}

	ucInReview := usecase.NovoMoverParaInReviewUseCase(projectRepo, taskRepo)

	err := ucInReview.Executar(dto.MoverParaInReviewRequest{
		TaskID:            string(taskID),
		UserID:            string(outroUserID),
		ComentarioEntrega: comentarioEntregaTeste,
	})
	if err == nil {
		t.Fatalf("esperava erro: somente assignee pode mover para InReview")
	}

	err = ucInReview.Executar(dto.MoverParaInReviewRequest{
		TaskID:            string(taskID),
		UserID:            string(assigneeID),
		ComentarioEntrega: comentarioEntregaTeste,
	})
	if err != nil {
		t.Fatalf("esperava mover para InReview ok, veio: %v", err)
	}

	tk, err := taskRepo.BuscarPorID(taskID)
	if err != nil {
		t.Fatalf("erro buscando task após mover para in review: %v", err)
	}
	if tk.Status() != task.InReview {
		t.Fatalf("esperava task em InReview, veio %s", tk.Status())
	}
	if tk.ComentarioEntrega() != comentarioEntregaTeste {
		t.Fatalf("esperava comentário de entrega salvo, veio %q", tk.ComentarioEntrega())
	}
}

func TestReprovar_RetornaParaToDo(t *testing.T) {
	db := openTestDB(t)

	projectRepo := projectrepo.NewProjectRepo(db)
	teamRepo := teamrepo.NewTeamRepo(db)
	taskRepo := taskrepo.NewTaskRepo(db)

	leaderID := seedUser(t, db, "Leader Teste")
	assigneeID := seedUser(t, db, "Assignee Teste")
	teamID := seedTeam(t, db, "Team Teste", leaderID)
	projectID := seedProject(t, db, teamID, "Projeto Teste", project.ConfiguracoesProject{})
	taskID := seedTask(t, db, projectID, "Task")

	ucAssign := usecase.NovoSelfAssignTaskUseCase(projectRepo, taskRepo)
	ucInReview := usecase.NovoMoverParaInReviewUseCase(projectRepo, taskRepo)
	ucReprovar := usecase.NovoReprovarTaskUseCase(projectRepo, teamRepo, taskRepo)

	if err := ucAssign.Executar(dto.SelfAssignRequest{
		TaskID: string(taskID),
		UserID: string(assigneeID),
	}); err != nil {
		t.Fatalf("erro no self-assign: %v", err)
	}

	if err := ucInReview.Executar(dto.MoverParaInReviewRequest{
		TaskID:            string(taskID),
		UserID:            string(assigneeID),
		ComentarioEntrega: comentarioEntregaTeste,
	}); err != nil {
		t.Fatalf("erro ao mover para in review: %v", err)
	}

	motivo := "Ajustar a validação antes de reenviar"
	err := ucReprovar.Executar(dto.ReprovarTaskRequest{
		TaskID:   string(taskID),
		LeaderID: string(leaderID),
		Motivo:   motivo,
	})
	if err != nil {
		t.Fatalf("esperava reprovar ok, veio: %v", err)
	}

	tk, err := taskRepo.BuscarPorID(taskID)
	if err != nil {
		t.Fatalf("erro buscando task após reprovar: %v", err)
	}
	if tk.Status() != task.ToDo {
		t.Fatalf("esperava task voltando para ToDo, veio %s", tk.Status())
	}
	if tk.ComentarioReview() != motivo {
		t.Fatalf("esperava comentário de review salvo, veio %q", tk.ComentarioReview())
	}
	if tk.ComentarioEntrega() != comentarioEntregaTeste {
		t.Fatalf("esperava manter comentário de entrega, veio %q", tk.ComentarioEntrega())
	}
}

func TestAprovar_LeaderOnly(t *testing.T) {
	db := openTestDB(t)

	projectRepo := projectrepo.NewProjectRepo(db)
	teamRepo := teamrepo.NewTeamRepo(db)
	taskRepo := taskrepo.NewTaskRepo(db)

	leaderID := seedUser(t, db, "Leader Teste")
	assigneeID := seedUser(t, db, "Assignee Teste")
	teamID := seedTeam(t, db, "Team Teste", leaderID)
	projectID := seedProject(t, db, teamID, "Projeto Teste", project.ConfiguracoesProject{})
	taskID := seedTask(t, db, projectID, "Task")

	ucAssign := usecase.NovoSelfAssignTaskUseCase(projectRepo, taskRepo)
	ucInReview := usecase.NovoMoverParaInReviewUseCase(projectRepo, taskRepo)
	ucApprove := usecase.NovoAprovarTaskUseCase(projectRepo, teamRepo, taskRepo)

	if err := ucAssign.Executar(dto.SelfAssignRequest{
		TaskID: string(taskID),
		UserID: string(assigneeID),
	}); err != nil {
		t.Fatalf("erro no self-assign: %v", err)
	}

	if err := ucInReview.Executar(dto.MoverParaInReviewRequest{
		TaskID:            string(taskID),
		UserID:            string(assigneeID),
		ComentarioEntrega: comentarioEntregaTeste,
	}); err != nil {
		t.Fatalf("erro ao mover para in review: %v", err)
	}

	err := ucApprove.Executar(dto.AprovarTaskRequest{
		TaskID:   string(taskID),
		LeaderID: string(assigneeID),
	})
	if err == nil {
		t.Fatalf("esperava erro: somente leader aprova")
	}

	err = ucApprove.Executar(dto.AprovarTaskRequest{
		TaskID:   string(taskID),
		LeaderID: string(leaderID),
	})
	if err != nil {
		t.Fatalf("esperava aprovar ok, veio: %v", err)
	}

	tk, err := taskRepo.BuscarPorID(taskID)
	if err != nil {
		t.Fatalf("erro buscando task após aprovar: %v", err)
	}
	if tk.Status() != task.Done {
		t.Fatalf("esperava task em Done, veio %s", tk.Status())
	}
	if tk.Outcome() == nil || *tk.Outcome() != task.OutcomeApproved {
		t.Fatalf("esperava outcome APPROVED")
	}
	if tk.ComentarioEntrega() != comentarioEntregaTeste {
		t.Fatalf("esperava manter comentário de entrega após aprovação, veio %q", tk.ComentarioEntrega())
	}
}
