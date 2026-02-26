package tests

import (
	"testing"

	projectdto "github.com/hudsontheuz/saas_kanban/internal/application/project/dto"
	projectusecase "github.com/hudsontheuz/saas_kanban/internal/application/project/usecase"
	teamdto "github.com/hudsontheuz/saas_kanban/internal/application/team/dto"
	teamusecase "github.com/hudsontheuz/saas_kanban/internal/application/team/usecase"
	"github.com/hudsontheuz/saas_kanban/internal/infrastructure/persistence/memory"
)

func TestTeamProject_RegraUmProjectAtivoPorTeam(t *testing.T) {
	teamRepo := memory.NovoTeamRepoEmMemoria()
	projectRepo := memory.NovoProjectRepoEmMemoria()

	ucCriarTeam := teamusecase.NovoCriarTeamUseCase(teamRepo)
	teamResp, err := ucCriarTeam.Executar(teamdto.CriarTeamRequest{
		Nome:     "Minha Team",
		LeaderID: "leader-1",
	})
	if err != nil {
		t.Fatalf("erro criando team: %v", err)
	}

	ucCriarProject := projectusecase.NovoCriarProjectUseCase(teamRepo, projectRepo)

	_, err = ucCriarProject.Executar(projectdto.CriarProjectRequest{
		TeamID:                      teamResp.TeamID,
		LeaderID:                    "leader-1",
		PermitirSoltarDoingParaTodo: true,
	})
	if err != nil {
		t.Fatalf("erro criando project: %v", err)
	}

	// tentar criar outro ativo deve falhar
	_, err = ucCriarProject.Executar(projectdto.CriarProjectRequest{
		TeamID:   teamResp.TeamID,
		LeaderID: "leader-1",
	})
	if err == nil {
		t.Fatalf("esperava erro: já existe project ativo")
	}
}
