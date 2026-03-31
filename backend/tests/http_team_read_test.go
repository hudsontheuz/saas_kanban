package tests

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	deliveryhttp "github.com/hudsontheuz/saas_kanban/delivery/http"
	authdto "github.com/hudsontheuz/saas_kanban/internal/auth/application/dto"
	authusecase "github.com/hudsontheuz/saas_kanban/internal/auth/application/usecase"
	authhttp "github.com/hudsontheuz/saas_kanban/internal/auth/delivery/http"
	authhash "github.com/hudsontheuz/saas_kanban/internal/auth/infrastructure/hash"
	authjwt "github.com/hudsontheuz/saas_kanban/internal/auth/infrastructure/jwt"
	projectusecase "github.com/hudsontheuz/saas_kanban/internal/project/application/usecase"
	projecthttp "github.com/hudsontheuz/saas_kanban/internal/project/delivery/http"
	projectrepo "github.com/hudsontheuz/saas_kanban/internal/project/infrastructure/persistence/gorm/repo"
	taskusecase "github.com/hudsontheuz/saas_kanban/internal/task/application/usecase"
	taskhttp "github.com/hudsontheuz/saas_kanban/internal/task/delivery/http"
	taskrepo "github.com/hudsontheuz/saas_kanban/internal/task/infrastructure/persistence/gorm/repo"
	teamusecase "github.com/hudsontheuz/saas_kanban/internal/team/application/usecase"
	teamhttp "github.com/hudsontheuz/saas_kanban/internal/team/delivery/http"
	teamrepo "github.com/hudsontheuz/saas_kanban/internal/team/infrastructure/persistence/gorm/repo"
	userrepo "github.com/hudsontheuz/saas_kanban/internal/user/infrastructure/persistence/gorm/repo"
)

func TestHTTP_Team_GetByID(t *testing.T) {
	db := openTestDB(t)

	repoUsuario := userrepo.NewUserRepo(db)
	repoProjeto := projectrepo.NewProjectRepo(db)
	repoTarefa := taskrepo.NewTaskRepo(db)
	repoTeam := teamrepo.NewTeamRepo(db)

	segredo := "segredo-teste"
	emissor := "saas_kanban"

	issuer, err := authjwt.NovoIssuer(segredo, emissor, time.Hour)
	if err != nil {
		t.Fatalf("erro criando issuer: %v", err)
	}

	validador, err := authjwt.NovoValidador(segredo, emissor)
	if err != nil {
		t.Fatalf("erro criando validador: %v", err)
	}

	hasher := authhash.NewBcryptHasher()
	ucRegister := authusecase.NovoRegisterUseCase(repoUsuario, hasher, issuer)
	ucLogin := authusecase.NovoLoginUseCase(repoUsuario, hasher, issuer)
	handlerAuth := authhttp.NewAuthHandler(ucRegister, ucLogin)

	ucCriarTeam := teamusecase.NovoCriarTeamUseCase(repoTeam)
	ucBuscarTeam := teamusecase.NovoBuscarTeamUseCase(repoTeam)
	handlerTeam := teamhttp.NewTeamHandler(ucCriarTeam, ucBuscarTeam, nil)

	ucCriarProject := projectusecase.NovoCriarProjectUseCase(repoTeam, repoProjeto)
	ucBuscarProjectAtivo := projectusecase.NovoBuscarProjectAtivoUseCase(repoProjeto)
	ucAtualizarSettingsProject := projectusecase.NovoAtualizarSettingsProjectUseCase(repoTeam, repoProjeto)
	ucFecharProject := projectusecase.NovoFecharProjectUseCase(repoTeam, repoProjeto)

	handlerProject := projecthttp.NewProjectHandler(
		ucCriarProject,
		ucBuscarProjectAtivo,
		ucAtualizarSettingsProject,
		ucFecharProject,
	)

	handlerTask := taskhttp.NewTaskHandlerWorkflowWithList(
		taskusecase.NovoCriarTaskUseCase(repoProjeto, repoTarefa),
		taskusecase.NovoListarTasksPorProjectUseCase(repoProjeto, repoTarefa),
		taskusecase.NovoSelfAssignTaskUseCase(repoProjeto, repoTarefa),
		taskusecase.NovoPausarTaskUseCase(repoProjeto, repoTarefa),
		taskusecase.NovoRetomarTaskUseCase(repoProjeto, repoTarefa),
		taskusecase.NovoMoverParaInReviewUseCase(repoProjeto, repoTarefa),
		taskusecase.NovoAprovarTaskUseCase(repoProjeto, repoTeam, repoTarefa),
		taskusecase.NovoReprovarTaskUseCase(repoProjeto, repoTeam, repoTarefa),
	)

	router := deliveryhttp.NewRouter(handlerAuth, handlerTeam, handlerProject, handlerTask, validador)
	server := httptest.NewServer(router)
	defer server.Close()

	_, err = ucRegister.Executar(authdto.RegisterRequest{
		Nome:  "Matheus",
		Email: "matheus@teste.com",
		Senha: "123456",
	})
	if err != nil {
		t.Fatalf("erro no register: %v", err)
	}

	loginResp, err := ucLogin.Executar(authdto.LoginRequest{
		Email: "matheus@teste.com",
		Senha: "123456",
	})
	if err != nil {
		t.Fatalf("erro no login: %v", err)
	}

	token := loginResp.Token
	teamID := criarTeamViaHTTP(t, server.URL, token, "Equipe Alpha")

	req, _ := http.NewRequest(http.MethodGet, server.URL+"/teams/"+teamID, nil)
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("erro buscando team por id: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("esperava 200, veio %d", resp.StatusCode)
	}

	var payload map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		t.Fatalf("erro lendo response: %v", err)
	}

	if payload["team_id"] != teamID {
		t.Fatalf("esperava team_id %s, veio %v", teamID, payload["team_id"])
	}

	if payload["nome"] != "Equipe Alpha" {
		t.Fatalf("esperava nome Equipe Alpha, veio %v", payload["nome"])
	}

	if payload["leader_id"] == "" || payload["leader_id"] == nil {
		t.Fatalf("esperava leader_id preenchido")
	}
}
