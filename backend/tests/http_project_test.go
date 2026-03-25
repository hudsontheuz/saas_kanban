package tests

import (
	"bytes"
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

func TestHTTP_Project_Workflow(t *testing.T) {
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
	handlerTeam := teamhttp.NewTeamHandler(ucCriarTeam)

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

	handlerTask := taskhttp.NewTaskHandlerWorkflow(
		taskusecase.NovoCriarTaskUseCase(repoProjeto, repoTarefa),
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

	projectID := criarProjectViaHTTP(t, server.URL, token, teamID, "Projeto Alpha", true)

	t.Run("não permite segundo projeto ativo", func(t *testing.T) {
		body, _ := json.Marshal(map[string]any{
			"nome":                            "Projeto Beta",
			"permitir_soltar_doing_para_todo": true,
		})

		req, _ := http.NewRequest(http.MethodPost, server.URL+"/teams/"+teamID+"/projects", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("erro criando segundo projeto: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusConflict {
			t.Fatalf("esperava 409, veio %d", resp.StatusCode)
		}
	})

	t.Run("busca projeto ativo", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, server.URL+"/teams/"+teamID+"/projects/active", nil)
		req.Header.Set("Authorization", "Bearer "+token)

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("erro buscando projeto ativo: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Fatalf("esperava 200, veio %d", resp.StatusCode)
		}

		var payload map[string]any
		if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
			t.Fatalf("erro decodificando projeto ativo: %v", err)
		}

		if payload["project_id"] == "" || payload["project_id"] == nil {
			t.Fatalf("esperava project_id preenchido")
		}
	})

	t.Run("atualiza settings", func(t *testing.T) {
		body, _ := json.Marshal(map[string]any{
			"permitir_soltar_doing_para_todo": false,
		})

		req, _ := http.NewRequest(http.MethodPatch, server.URL+"/projects/"+projectID+"/settings", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("erro atualizando settings: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Fatalf("esperava 200, veio %d", resp.StatusCode)
		}
	})

	t.Run("fecha projeto", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodPost, server.URL+"/projects/"+projectID+"/close", nil)
		req.Header.Set("Authorization", "Bearer "+token)

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("erro fechando projeto: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Fatalf("esperava 200, veio %d", resp.StatusCode)
		}
	})

	t.Run("fechar de novo retorna conflito", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodPost, server.URL+"/projects/"+projectID+"/close", nil)
		req.Header.Set("Authorization", "Bearer "+token)

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("erro fechando projeto pela segunda vez: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusConflict {
			t.Fatalf("esperava 409, veio %d", resp.StatusCode)
		}
	})

	t.Run("após fechar permite criar novo projeto", func(t *testing.T) {
		novoProjectID := criarProjectViaHTTP(t, server.URL, token, teamID, "Projeto Beta", false)
		if novoProjectID == "" {
			t.Fatalf("esperava novo project_id preenchido")
		}
	})
}

func criarTeamViaHTTP(t *testing.T, baseURL, token, nome string) string {
	t.Helper()

	body, _ := json.Marshal(map[string]any{
		"nome": nome,
	})

	req, _ := http.NewRequest(http.MethodPost, baseURL+"/teams", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("erro criando team: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("esperava 201 ao criar team, veio %d", resp.StatusCode)
	}

	var payload map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		t.Fatalf("erro lendo response de team: %v", err)
	}

	teamID, _ := payload["team_id"].(string)
	if teamID == "" {
		t.Fatalf("esperava team_id preenchido")
	}

	return teamID
}

func criarProjectViaHTTP(t *testing.T, baseURL, token, teamID, nome string, permitirSoltar bool) string {
	t.Helper()

	body, _ := json.Marshal(map[string]any{
		"nome":                            nome,
		"permitir_soltar_doing_para_todo": permitirSoltar,
	})

	req, _ := http.NewRequest(http.MethodPost, baseURL+"/teams/"+teamID+"/projects", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("erro criando project: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("esperava 201 ao criar project, veio %d", resp.StatusCode)
	}

	var payload map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		t.Fatalf("erro lendo response de project: %v", err)
	}

	projectID, _ := payload["project_id"].(string)
	if projectID == "" {
		t.Fatalf("esperava project_id preenchido")
	}

	return projectID
}
