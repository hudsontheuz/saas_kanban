package tests

import (
	"testing"
	"time"

	project "github.com/hudsontheuz/saas_kanban/internal/project/domain"
	projectrepo "github.com/hudsontheuz/saas_kanban/internal/project/infrastructure/persistence/gorm/repo"
	dto "github.com/hudsontheuz/saas_kanban/internal/task/application/dto"
	usecase "github.com/hudsontheuz/saas_kanban/internal/task/application/usecase"
	taskrepo "github.com/hudsontheuz/saas_kanban/internal/task/infrastructure/persistence/gorm/repo"
)

func TestSlice_CriarTask_OK(t *testing.T) {
	db := openTestDB(t)

	leaderID := seedUser(t, db, "Leader Projeto")
	teamID := seedTeam(t, db, "Team Teste", leaderID)

	projectRepo := projectrepo.NewProjectRepo(db)
	taskRepo := taskrepo.NewTaskRepo(db)

	p, err := project.NovoProject(teamID, "Projeto Teste", project.ConfiguracoesProject{
		PermitirSoltarDoingParaTodo: true,
	})
	if err != nil {
		t.Fatalf("erro ao criar project: %v", err)
	}
	if err := projectRepo.Salvar(p); err != nil {
		t.Fatalf("erro ao salvar project: %v", err)
	}

	uc := usecase.NovoCriarTaskUseCase(projectRepo, taskRepo)
	resp, err := uc.Executar(dto.CriarTaskRequest{
		ProjectID: string(p.ID()),
		Titulo:    "Primeira task",
		Descricao: "Descrição da primeira task",
		CriadorID: string(leaderID),
	})
	if err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}
	if resp.TaskID == "" {
		t.Fatalf("esperava TaskID preenchido")
	}
}

func TestSlice_ProjectFechado_BloqueiaCriarTask(t *testing.T) {
	db := openTestDB(t)

	leaderID := seedUser(t, db, "Leader Projeto")
	teamID := seedTeam(t, db, "Team Teste", leaderID)

	projectRepo := projectrepo.NewProjectRepo(db)
	taskRepo := taskrepo.NewTaskRepo(db)

	p, err := project.NovoProject(teamID, "Projeto Teste", project.ConfiguracoesProject{})
	if err != nil {
		t.Fatalf("erro ao criar project: %v", err)
	}
	if err := p.Fechar(time.Now().UTC()); err != nil {
		t.Fatalf("erro ao fechar project: %v", err)
	}

	if err := projectRepo.Salvar(p); err != nil {
		t.Fatalf("erro ao salvar project: %v", err)
	}

	uc := usecase.NovoCriarTaskUseCase(projectRepo, taskRepo)
	_, err = uc.Executar(dto.CriarTaskRequest{
		ProjectID: string(p.ID()),
		Titulo:    "Task bloqueada",
		Descricao: "Descrição da task bloqueada",
		CriadorID: string(leaderID),
	})
	if err == nil {
		t.Fatalf("esperava erro por projeto fechado")
	}
}
