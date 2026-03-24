package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	deliveryhttp "github.com/hudsontheuz/saas_kanban/delivery/http"
	authusecase "github.com/hudsontheuz/saas_kanban/internal/auth/application/usecase"
	authhttp "github.com/hudsontheuz/saas_kanban/internal/auth/delivery/http"
	authhash "github.com/hudsontheuz/saas_kanban/internal/auth/infrastructure/hash"
	authjwt "github.com/hudsontheuz/saas_kanban/internal/auth/infrastructure/jwt"
	projectdto "github.com/hudsontheuz/saas_kanban/internal/project/application/dto"
	projectusecase "github.com/hudsontheuz/saas_kanban/internal/project/application/usecase"
	projecthttp "github.com/hudsontheuz/saas_kanban/internal/project/delivery/http"
	projectdomain "github.com/hudsontheuz/saas_kanban/internal/project/domain"
	projectrepo "github.com/hudsontheuz/saas_kanban/internal/project/infrastructure/persistence/gorm/repo"
	taskusecase "github.com/hudsontheuz/saas_kanban/internal/task/application/usecase"
	taskhttp "github.com/hudsontheuz/saas_kanban/internal/task/delivery/http"
	taskrepo "github.com/hudsontheuz/saas_kanban/internal/task/infrastructure/persistence/gorm/repo"
	teamusecase "github.com/hudsontheuz/saas_kanban/internal/team/application/usecase"
	teamhttp "github.com/hudsontheuz/saas_kanban/internal/team/delivery/http"
	teamrepo "github.com/hudsontheuz/saas_kanban/internal/team/infrastructure/persistence/gorm/repo"
	userrepo "github.com/hudsontheuz/saas_kanban/internal/user/infrastructure/persistence/gorm/repo"
)

func TestHTTP_Project_CreateEBuscarAtivo(t *testing.T) {
	db := openTestDB(t)

	repoUsuario := userrepo.NewUserRepo(db)
	repoProjeto := projectrepo.NewProjectRepo(db)
	repoTarefa := taskrepo.NewTaskRepo(db)
	repoTeam := teamrepo.NewTeamRepo(db)

	leaderID := seedUser(t, db, "Leader Project HTTP")
	teamID := seedTeam(t, db, "Team HTTP Project", leaderID)

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
	handlerTeam := teamhttp.NewTeamHandler(ucCriarTeam)

	ucCriarProject := projectusecase.NovoCriarProjectUseCase(repoTeam, repoProjeto)
	ucFecharProject := projectusecase.NovoFecharProjectUseCase(repoTeam, repoProjeto)
	ucBuscarProjectAtivo := projectusecase.NovoBuscarProjectAtivoUseCase(repoProjeto)
	handlerProject := projecthttp.NewProjectHandler(ucCriarProject, ucFecharProject, ucBuscarProjectAtivo)

	handlerTask := taskhttp.NewTaskHandlerWorkflow(
		nil,
		taskusecase.NovoListarTasksPorProjectUseCase(repoProjeto, repoTarefa),
		taskusecase.NovoSelfAssignTaskUseCase(repoProjeto, repoTarefa),
		nil,
		nil,
		nil,
		nil,
		nil,
	)

	router := deliveryhttp.NewRouter(handlerAuth, handlerTeam, handlerProject, handlerTask, validador)
	server := httptest.NewServer(router)
	defer server.Close()

	leaderToken := gerarJWT(t, segredo, emissor, string(leaderID), time.Now().Add(1*time.Hour))

	createBody, _ := json.Marshal(map[string]any{
		"nome":                            "Projeto HTTP",
		"permitir_soltar_doing_para_todo": true,
	})

	req, _ := http.NewRequest(http.MethodPost, server.URL+"/teams/"+string(teamID)+"/projects", bytes.NewReader(createBody))
	req.Header.Set("Authorization", "Bearer "+leaderToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("erro no create project: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("esperava 201 no create project, veio %d", resp.StatusCode)
	}

	var createResp projectdto.CriarProjectResponse
	if err := json.NewDecoder(resp.Body).Decode(&createResp); err != nil {
		t.Fatalf("erro lendo create project response: %v", err)
	}
	if createResp.ProjectID == "" {
		t.Fatalf("esperava project_id preenchido")
	}

	req, _ = http.NewRequest(http.MethodGet, server.URL+"/teams/"+string(teamID)+"/projects/active", nil)
	req.Header.Set("Authorization", "Bearer "+leaderToken)

	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("erro no get active project: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("esperava 200 no get active project, veio %d", resp.StatusCode)
	}

	var activeResp projectdto.BuscarProjectAtivoResponse
	if err := json.NewDecoder(resp.Body).Decode(&activeResp); err != nil {
		t.Fatalf("erro lendo active project response: %v", err)
	}

	if activeResp.ProjectID != createResp.ProjectID {
		t.Fatalf("esperava project_id %s, veio %s", createResp.ProjectID, activeResp.ProjectID)
	}
	if activeResp.Nome != "Projeto HTTP" {
		t.Fatalf("esperava nome do projeto na consulta ativa")
	}
	if !activeResp.PermitirSoltarDoingParaTodo {
		t.Fatalf("esperava permitir_soltar_doing_para_todo=true")
	}
}

func TestHTTP_Project_Close(t *testing.T) {
	db := openTestDB(t)

	repoUsuario := userrepo.NewUserRepo(db)
	repoProjeto := projectrepo.NewProjectRepo(db)
	repoTarefa := taskrepo.NewTaskRepo(db)
	repoTeam := teamrepo.NewTeamRepo(db)

	leaderID := seedUser(t, db, "Leader Close HTTP")
	teamID := seedTeam(t, db, "Team Close HTTP", leaderID)
	projectID := seedProject(t, db, teamID, "Projeto para fechar", projectdomain.ConfiguracoesProject{})

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
	handlerTeam := teamhttp.NewTeamHandler(ucCriarTeam)

	ucCriarProject := projectusecase.NovoCriarProjectUseCase(repoTeam, repoProjeto)
	ucFecharProject := projectusecase.NovoFecharProjectUseCase(repoTeam, repoProjeto)
	ucBuscarProjectAtivo := projectusecase.NovoBuscarProjectAtivoUseCase(repoProjeto)
	handlerProject := projecthttp.NewProjectHandler(ucCriarProject, ucFecharProject, ucBuscarProjectAtivo)

	handlerTask := taskhttp.NewTaskHandlerWorkflow(
		nil,
		taskusecase.NovoListarTasksPorProjectUseCase(repoProjeto, repoTarefa),
		taskusecase.NovoSelfAssignTaskUseCase(repoProjeto, repoTarefa),
		nil,
		nil,
		nil,
		nil,
		nil,
	)

	router := deliveryhttp.NewRouter(handlerAuth, handlerTeam, handlerProject, handlerTask, validador)
	server := httptest.NewServer(router)
	defer server.Close()

	leaderToken := gerarJWT(t, segredo, emissor, string(leaderID), time.Now().Add(1*time.Hour))

	req, _ := http.NewRequest(http.MethodPost, server.URL+"/projects/"+string(projectID)+"/close", nil)
	req.Header.Set("Authorization", "Bearer "+leaderToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("erro no close project: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("esperava 200 no close project, veio %d", resp.StatusCode)
	}

	proj, err := repoProjeto.BuscarPorID(projectdomain.ProjectID(projectID))
	if err != nil {
		t.Fatalf("erro buscando projeto após close: %v", err)
	}
	if !proj.EstaFechado() {
		t.Fatalf("esperava projeto fechado após endpoint /close")
	}
	if proj.FechadoEm() == nil {
		t.Fatalf("esperava fechado_em preenchido após close")
	}
}
