package tests

import (
	"testing"

	projectdto "github.com/hudsontheuz/saas_kanban/internal/project/application/dto"
	projectusecase "github.com/hudsontheuz/saas_kanban/internal/project/application/usecase"
	projectrepo "github.com/hudsontheuz/saas_kanban/internal/project/infrastructure/persistence/gorm/repo"
	teamdto "github.com/hudsontheuz/saas_kanban/internal/team/application/dto"
	teamusecase "github.com/hudsontheuz/saas_kanban/internal/team/application/usecase"
	teamrepo "github.com/hudsontheuz/saas_kanban/internal/team/infrastructure/persistence/gorm/repo"
)

func TestTeamProject_RegraUmProjectAtivoPorTeam(t *testing.T) {
	db := openTestDB(t)

	leaderID := seedUser(t, db, "Leader Teste")

	teamRepo := teamrepo.NewTeamRepo(db)
	projectRepo := projectrepo.NewProjectRepo(db)

	ucCriarTeam := teamusecase.NovoCriarTeamUseCase(teamRepo)
	teamResp, err := ucCriarTeam.Executar(teamdto.CriarTeamRequest{
		Nome:     "Minha Team",
		LeaderID: string(leaderID),
	})
	if err != nil {
		t.Fatalf("erro criando team: %v", err)
	}

	ucCriarProject := projectusecase.NovoCriarProjectUseCase(teamRepo, projectRepo)

	_, err = ucCriarProject.Executar(projectdto.CriarProjectRequest{
		TeamID:                      teamResp.TeamID,
		LeaderID:                    string(leaderID),
		Nome:                        "Projeto 1",
		PermitirSoltarDoingParaTodo: true,
	})
	if err != nil {
		t.Fatalf("erro criando project: %v", err)
	}

	_, err = ucCriarProject.Executar(projectdto.CriarProjectRequest{
		TeamID:   teamResp.TeamID,
		LeaderID: string(leaderID),
		Nome:     "Projeto 2",
	})
	if err == nil {
		t.Fatalf("esperava erro: já existe project ativo")
	}
}
