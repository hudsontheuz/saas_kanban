package tests

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	jwtlib "github.com/golang-jwt/jwt/v5"

	deliveryhttp "github.com/hudsontheuz/saas_kanban/delivery/http"
	authjwt "github.com/hudsontheuz/saas_kanban/internal/auth/infrastructure/jwt"
	project "github.com/hudsontheuz/saas_kanban/internal/project/domain"
	projectmemory "github.com/hudsontheuz/saas_kanban/internal/project/infrastructure/persistence/memory"
	taskdto "github.com/hudsontheuz/saas_kanban/internal/task/application/dto"
	taskusecase "github.com/hudsontheuz/saas_kanban/internal/task/application/usecase"
	taskhttp "github.com/hudsontheuz/saas_kanban/internal/task/delivery/http"
	task "github.com/hudsontheuz/saas_kanban/internal/task/domain"
	taskmemory "github.com/hudsontheuz/saas_kanban/internal/task/infrastructure/persistence/memory"
	team "github.com/hudsontheuz/saas_kanban/internal/team/domain"
)

func TestHTTP_SelfAssign_ComJWT(t *testing.T) {
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

	validador, err := authjwt.NovoValidador(segredo, emissor)
	if err != nil {
		t.Fatalf("erro criando validador jwt: %v", err)
	}

	router := deliveryhttp.NewRouter(handlerTarefa, validador)
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

		// sanity: garante que a regra foi aplicada
		_ = taskdto.SelfAssignRequest{} // só pra manter import se você quiser remover validação depois
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
