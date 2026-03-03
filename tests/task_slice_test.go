package tests

import (
	"testing"
	"time"

	"github.com/hudsontheuz/saas_kanban/infrastructure/persistence/memory"
	"github.com/hudsontheuz/saas_kanban/internal/application/task/dto"
	usecase "github.com/hudsontheuz/saas_kanban/internal/application/task/usecase"
	"github.com/hudsontheuz/saas_kanban/internal/domain/project"
	"github.com/hudsontheuz/saas_kanban/internal/domain/team"
)

func TestSlice_CriarTask_OK(t *testing.T) {
	projectRepo := memory.NovoProjectRepoEmMemoria()
	taskRepo := memory.NovoTaskRepoEmMemoria()

	p, err := project.NovoProject(team.TeamID("1"), "Projeto Teste", project.ConfiguracoesProject{
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
		CriadorID: "1",
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

	p, err := project.NovoProject(team.TeamID("1"), "Projeto Teste", project.ConfiguracoesProject{})
	if err != nil {
		t.Fatalf("erro ao criar project: %v", err)
	}
	p.Fechar(time.Now().UTC())
	_ = projectRepo.Salvar(p)

	uc := usecase.NovoCriarTaskUseCase(projectRepo, taskRepo)
	_, err = uc.Executar(dto.CriarTaskRequest{
		ProjectID: string(p.ID()),
		Titulo:    "Task bloqueada",
	})
	if err == nil {
		t.Fatalf("esperava erro por projeto fechado")
	}
}
