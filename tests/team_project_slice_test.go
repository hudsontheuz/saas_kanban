package tests

import (
	"testing"

	projectdto "github.com/hudsontheuz/saas_kanban/internal/project/application/dto"
	projectusecase "github.com/hudsontheuz/saas_kanban/internal/project/application/usecase"
	projectmemory "github.com/hudsontheuz/saas_kanban/internal/project/infrastructure/persistence/memory"
	teamdto "github.com/hudsontheuz/saas_kanban/internal/team/application/dto"
	teamusecase "github.com/hudsontheuz/saas_kanban/internal/team/application/usecase"
	teammemory "github.com/hudsontheuz/saas_kanban/internal/team/infrastructure/persistence/memory"
)

func TestTeamProject_RegraUmProjectAtivoPorTeam(t *testing.T) {
	teamRepo := teammemory.NovoTeamRepoEmMemoria()
	projectRepo := projectmemory.NovoProjectRepoEmMemoria()

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
		Nome:                        "Projeto 1",
		PermitirSoltarDoingParaTodo: true,
	})
	if err != nil {
		t.Fatalf("erro criando project: %v", err)
	}

	// tentar criar outro ativo deve falhar
	_, err = ucCriarProject.Executar(projectdto.CriarProjectRequest{
		TeamID:   teamResp.TeamID,
		LeaderID: "leader-1",
		Nome:     "Projeto 2",
	})
	if err == nil {
		t.Fatalf("esperava erro: já existe project ativo")
	}
}
