package tests

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	jwtlib "github.com/golang-jwt/jwt/v5"

	deliveryhttp "github.com/hudsontheuz/saas_kanban/delivery/http"
	authdto "github.com/hudsontheuz/saas_kanban/internal/auth/application/dto"
	authusecase "github.com/hudsontheuz/saas_kanban/internal/auth/application/usecase"
	authhttp "github.com/hudsontheuz/saas_kanban/internal/auth/delivery/http"
	authhash "github.com/hudsontheuz/saas_kanban/internal/auth/infrastructure/hash"
	authjwt "github.com/hudsontheuz/saas_kanban/internal/auth/infrastructure/jwt"
	project "github.com/hudsontheuz/saas_kanban/internal/project/domain"
	projectmemory "github.com/hudsontheuz/saas_kanban/internal/project/infrastructure/persistence/memory"
	taskdto "github.com/hudsontheuz/saas_kanban/internal/task/application/dto"
	taskusecase "github.com/hudsontheuz/saas_kanban/internal/task/application/usecase"
	taskhttp "github.com/hudsontheuz/saas_kanban/internal/task/delivery/http"
	task "github.com/hudsontheuz/saas_kanban/internal/task/domain"
	taskmemory "github.com/hudsontheuz/saas_kanban/internal/task/infrastructure/persistence/memory"
	team "github.com/hudsontheuz/saas_kanban/internal/team/domain"
	usermemory "github.com/hudsontheuz/saas_kanban/internal/user/infrastructure/persistence/memory"
)

func TestHTTP_SelfAssign_ComJWT(t *testing.T) {
	repoUsuario := usermemory.NovoUserRepoEmMemoria()
	repoProjeto := projectmemory.NovoProjectRepoEmMemoria()
	repoTarefa := taskmemory.NovoTaskRepoEmMemoria()

	p, err := project.NovoProject(team.TeamID("1"), "Projeto Teste", project.ConfiguracoesProject{})
	if err != nil {
		t.Fatalf("erro ao criar project: %v", err)
	}
	_ = repoProjeto.Salvar(p)

	tk, err := task.NovaTask(p.ID(), "Task Teste")
	if err != nil {
		t.Fatalf("erro ao criar task: %v", err)
	}
	_ = repoTarefa.Salvar(tk)

	casoUso := taskusecase.NovoSelfAssignTaskUseCase(repoProjeto, repoTarefa)
	handlerTarefa := taskhttp.NewTaskHandler(casoUso)

	segredo := "segredo-teste"
	emissor := "saas_kanban"

	issuer, err := authjwt.NovoIssuer(segredo, emissor, time.Hour)
	if err != nil {
		t.Fatalf("erro criando issuer jwt: %v", err)
	}
	validador, err := authjwt.NovoValidador(segredo, emissor)
	if err != nil {
		t.Fatalf("erro criando validador jwt: %v", err)
	}

	hasher := authhash.NewBcryptHasher()
	ucRegister := authusecase.NovoRegisterUseCase(repoUsuario, hasher, issuer)
	ucLogin := authusecase.NovoLoginUseCase(repoUsuario, hasher, issuer)
	handlerAuth := authhttp.NewAuthHandler(ucRegister, ucLogin)

	router := deliveryhttp.NewRouter(handlerAuth, handlerTarefa, validador)
	server := httptest.NewServer(router)
	defer server.Close()

	url := server.URL + "/tasks/" + string(tk.ID()) + "/self-assign"

	t.Run("sem Authorization retorna 401", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodPost, url, nil)
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("request falhou: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusUnauthorized {
			t.Fatalf("esperava 401, veio %d", resp.StatusCode)
		}
	})

	t.Run("token inválido retorna 401", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodPost, url, nil)
		req.Header.Set("Authorization", "Bearer token_invalido")

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("request falhou: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusUnauthorized {
			t.Fatalf("esperava 401, veio %d", resp.StatusCode)
		}
	})

	t.Run("token expirado retorna 401", func(t *testing.T) {
		tokenExpirado := gerarJWT(t, segredo, emissor, "2", time.Now().Add(-1*time.Minute))

		req, _ := http.NewRequest(http.MethodPost, url, nil)
		req.Header.Set("Authorization", "Bearer "+tokenExpirado)

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("request falhou: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusUnauthorized {
			t.Fatalf("esperava 401, veio %d", resp.StatusCode)
		}
	})

	t.Run("token válido retorna 200 e ok true", func(t *testing.T) {
		tokenValido := gerarJWT(t, segredo, emissor, "2", time.Now().Add(1*time.Hour))

		req, _ := http.NewRequest(http.MethodPost, url, nil)
		req.Header.Set("Authorization", "Bearer "+tokenValido)

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("request falhou: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Fatalf("esperava 200, veio %d", resp.StatusCode)
		}

		var body map[string]any
		if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
			t.Fatalf("erro lendo json: %v", err)
		}

		ok, _ := body["ok"].(bool)
		if !ok {
			t.Fatalf("esperava ok=true, veio: %#v", body)
		}

		_ = taskdto.SelfAssignRequest{}
	})

	t.Run("token emitido pelo login também funciona", func(t *testing.T) {
		_, err := ucRegister.Executar(authdto.RegisterRequest{Nome: "Matheus", Email: "matheus@teste.com", Senha: "123456"})
		if err != nil {
			t.Fatalf("erro no register: %v", err)
		}

		loginResp, err := ucLogin.Executar(authdto.LoginRequest{Email: "matheus@teste.com", Senha: "123456"})
		if err != nil {
			t.Fatalf("erro no login: %v", err)
		}

		tk2, err := task.NovaTask(p.ID(), "Task com token real")
		if err != nil {
			t.Fatalf("erro ao criar task: %v", err)
		}
		_ = repoTarefa.Salvar(tk2)

		req, _ := http.NewRequest(http.MethodPost, server.URL+"/tasks/"+string(tk2.ID())+"/self-assign", nil)
		req.Header.Set("Authorization", "Bearer "+loginResp.Token)

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("request falhou: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Fatalf("esperava 200, veio %d", resp.StatusCode)
		}
	})
}

func gerarJWT(t *testing.T, segredo, emissor, sub string, exp time.Time) string {
	t.Helper()

	claims := jwtlib.RegisteredClaims{
		Issuer:    emissor,
		Subject:   sub,
		ExpiresAt: jwtlib.NewNumericDate(exp),
		IssuedAt:  jwtlib.NewNumericDate(time.Now()),
	}

	token := jwtlib.NewWithClaims(jwtlib.SigningMethodHS256, claims)

	assinado, err := token.SignedString([]byte(segredo))
	if err != nil {
		t.Fatalf("erro assinando token: %v", err)
	}
	return assinado
}
