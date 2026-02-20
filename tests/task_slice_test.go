package tests

import (
	"testing"
	"time"

	"github.com/hudsontheuz/saas_kanban/internal/application/task/dto"
	usecase "github.com/hudsontheuz/saas_kanban/internal/application/task/usecase"
	"github.com/hudsontheuz/saas_kanban/internal/domain/project"
	"github.com/hudsontheuz/saas_kanban/internal/domain/team"
	"github.com/hudsontheuz/saas_kanban/internal/infrastructure/persistence/memory"
)

func TestSlice_CriarTask_OK(t *testing.T) {
	projectRepo := memory.NovoProjectRepoEmMemoria()
	taskRepo := memory.NovoTaskRepoEmMemoria()

	p, _ := project.NovoProject(team.TeamID("team-1"), project.ConfiguracoesProject{
		PermitirSoltarDoingParaTodo: true,
	})
	_ = projectRepo.Salvar(p)

	uc := usecase.NovoCriarTaskUseCase(projectRepo, taskRepo)
	resp, err := uc.Executar(taskdto.CriarTaskRequest{
		ProjectID: string(p.ID()),
		Titulo:    "Primeira task",
		CriadorID: "user-1",
	})
	if err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}
	if resp.TaskID == "" {
		t.Fatalf("esperava TaskID preenchido")
	}
}

func TestSlice_ProjectFechado_BloqueiaCriarTask(t *testing.T) {
	projectRepo := memory.NovoProjectRepoEmMemoria()
	taskRepo := memory.NovoTaskRepoEmMemoria()

	p, _ := project.NovoProject(team.TeamID("team-1"), project.ConfiguracoesProject{})
	p.Fechar(time.Now().UTC())
	_ = projectRepo.Salvar(p)

	uc := usecase.NovoCriarTaskUseCase(projectRepo, taskRepo)
	_, err := uc.Executar(taskdto.CriarTaskRequest{
		ProjectID: string(p.ID()),
		Titulo:    "Task bloqueada",
	})
	if err == nil {
		t.Fatalf("esperava erro por projeto fechado")
	}
}
