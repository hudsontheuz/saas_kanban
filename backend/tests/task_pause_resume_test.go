package tests

import (
	"testing"

	project "github.com/hudsontheuz/saas_kanban/internal/project/domain"
	projectrepo "github.com/hudsontheuz/saas_kanban/internal/project/infrastructure/persistence/gorm/repo"
	dto "github.com/hudsontheuz/saas_kanban/internal/task/application/dto"
	usecase "github.com/hudsontheuz/saas_kanban/internal/task/application/usecase"
	task "github.com/hudsontheuz/saas_kanban/internal/task/domain"
	taskrepo "github.com/hudsontheuz/saas_kanban/internal/task/infrastructure/persistence/gorm/repo"
)

func TestPause_NaoContaNoLimiteGlobal_SelfAssignEmOutraTaskFunciona(t *testing.T) {
	db := openTestDB(t)

	projectRepo := projectrepo.NewProjectRepo(db)
	taskRepo := taskrepo.NewTaskRepo(db)

	leaderID := seedUser(t, db, "Leader Teste")
	userID := seedUser(t, db, "Executor Teste")
	teamID := seedTeam(t, db, "Team Teste", leaderID)
	projectID := seedProject(t, db, teamID, "Projeto Teste", project.ConfiguracoesProject{})
	t1ID := seedTask(t, db, projectID, "T1")
	t2ID := seedTask(t, db, projectID, "T2")

	ucAssign := usecase.NovoSelfAssignTaskUseCase(projectRepo, taskRepo)
	if err := ucAssign.Executar(dto.SelfAssignRequest{
		TaskID: string(t1ID),
		UserID: string(userID),
	}); err != nil {
		t.Fatalf("esperava self-assign em t1 ok, veio: %v", err)
	}

	ucPause := usecase.NovoPausarTaskUseCase(projectRepo, taskRepo)
	if err := ucPause.Executar(dto.PausarTaskRequest{
		TaskID: string(t1ID),
		UserID: string(userID),
	}); err != nil {
		t.Fatalf("esperava pausar ok, veio: %v", err)
	}

	if err := ucAssign.Executar(dto.SelfAssignRequest{
		TaskID: string(t2ID),
		UserID: string(userID),
	}); err != nil {
		t.Fatalf("esperava self-assign em t2 ok (t1 pausada), veio: %v", err)
	}

	t1, err := taskRepo.BuscarPorID(t1ID)
	if err != nil {
		t.Fatalf("erro buscando t1 após pausa: %v", err)
	}
	if !t1.IsPaused() {
		t.Fatalf("esperava t1 pausada")
	}

	t2, err := taskRepo.BuscarPorID(t2ID)
	if err != nil {
		t.Fatalf("erro buscando t2 após self-assign: %v", err)
	}
	if t2.Status() != task.Doing {
		t.Fatalf("esperava t2 em DOING, veio %s", t2.Status())
	}
	if t2.Assignee() == nil || *t2.Assignee() != userID {
		t.Fatalf("esperava t2 atribuída ao usuário %s", userID)
	}
}

func TestRetomar_FalhaSeExisteOutraDoingAtiva(t *testing.T) {
	db := openTestDB(t)

	projectRepo := projectrepo.NewProjectRepo(db)
	taskRepo := taskrepo.NewTaskRepo(db)

	leaderID := seedUser(t, db, "Leader Teste")
	userID := seedUser(t, db, "Executor Teste")
	teamID := seedTeam(t, db, "Team Teste", leaderID)
	projectID := seedProject(t, db, teamID, "Projeto Teste", project.ConfiguracoesProject{})
	t1ID := seedTask(t, db, projectID, "T1")
	t2ID := seedTask(t, db, projectID, "T2")

	ucAssign := usecase.NovoSelfAssignTaskUseCase(projectRepo, taskRepo)
	ucPause := usecase.NovoPausarTaskUseCase(projectRepo, taskRepo)
	ucResume := usecase.NovoRetomarTaskUseCase(projectRepo, taskRepo)

	if err := ucAssign.Executar(dto.SelfAssignRequest{
		TaskID: string(t1ID),
		UserID: string(userID),
	}); err != nil {
		t.Fatalf("erro no self-assign de t1: %v", err)
	}

	if err := ucPause.Executar(dto.PausarTaskRequest{
		TaskID: string(t1ID),
		UserID: string(userID),
	}); err != nil {
		t.Fatalf("erro ao pausar t1: %v", err)
	}

	if err := ucAssign.Executar(dto.SelfAssignRequest{
		TaskID: string(t2ID),
		UserID: string(userID),
	}); err != nil {
		t.Fatalf("erro no self-assign de t2: %v", err)
	}

	if err := ucResume.Executar(dto.RetomarTaskRequest{
		TaskID: string(t1ID),
		UserID: string(userID),
	}); err == nil {
		t.Fatalf("esperava erro ao retomar: existe outra DOING ativa")
	}

	t1, err := taskRepo.BuscarPorID(t1ID)
	if err != nil {
		t.Fatalf("erro buscando t1 após retomar falhar: %v", err)
	}
	if !t1.IsPaused() {
		t.Fatalf("esperava t1 continuar pausada")
	}

	t2, err := taskRepo.BuscarPorID(t2ID)
	if err != nil {
		t.Fatalf("erro buscando t2 após retomar falhar: %v", err)
	}
	if t2.Status() != task.Doing || t2.IsPaused() {
		t.Fatalf("esperava t2 continuar em DOING ativa")
	}
}
