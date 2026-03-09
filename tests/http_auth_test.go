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
	projectmemory "github.com/hudsontheuz/saas_kanban/internal/project/infrastructure/persistence/memory"
	taskusecase "github.com/hudsontheuz/saas_kanban/internal/task/application/usecase"
	taskhttp "github.com/hudsontheuz/saas_kanban/internal/task/delivery/http"
	taskmemory "github.com/hudsontheuz/saas_kanban/internal/task/infrastructure/persistence/memory"
	usermemory "github.com/hudsontheuz/saas_kanban/internal/user/infrastructure/persistence/memory"
)

func TestHTTP_Auth_RegisterELogin(t *testing.T) {
	repoUsuario := usermemory.NovoUserRepoEmMemoria()
	repoProjeto := projectmemory.NovoProjectRepoEmMemoria()
	repoTarefa := taskmemory.NovoTaskRepoEmMemoria()

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

	handlerTask := taskhttp.NewTaskHandler(taskusecase.NovoSelfAssignTaskUseCase(repoProjeto, repoTarefa))
	router := deliveryhttp.NewRouter(handlerAuth, handlerTask, validador)
	server := httptest.NewServer(router)
	defer server.Close()

	registerBody, _ := json.Marshal(authdto.RegisterRequest{Nome: "Matheus", Email: "hudson@teste.com", Senha: "123456"})
	resp, err := http.Post(server.URL+"/auth/register", "application/json", bytes.NewReader(registerBody))
	if err != nil {
		t.Fatalf("erro no register: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("esperava 201 no register, veio %d", resp.StatusCode)
	}

	var registerResp authdto.AuthResponse
	if err := json.NewDecoder(resp.Body).Decode(&registerResp); err != nil {
		t.Fatalf("erro lendo register response: %v", err)
	}
	if registerResp.Token == "" || registerResp.UserID == "" {
		t.Fatalf("register deveria retornar token e user_id: %#v", registerResp)
	}

	resp, err = http.Post(server.URL+"/auth/register", "application/json", bytes.NewReader(registerBody))
	if err != nil {
		t.Fatalf("erro no register duplicado: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusConflict {
		t.Fatalf("esperava 409 no register duplicado, veio %d", resp.StatusCode)
	}

	loginBody, _ := json.Marshal(authdto.LoginRequest{Email: "hudson@teste.com", Senha: "123456"})
	resp, err = http.Post(server.URL+"/auth/login", "application/json", bytes.NewReader(loginBody))
	if err != nil {
		t.Fatalf("erro no login: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("esperava 200 no login, veio %d", resp.StatusCode)
	}

	var loginResp authdto.AuthResponse
	if err := json.NewDecoder(resp.Body).Decode(&loginResp); err != nil {
		t.Fatalf("erro lendo login response: %v", err)
	}
	if loginResp.Token == "" || loginResp.UserID == "" {
		t.Fatalf("login deveria retornar token e user_id: %#v", loginResp)
	}

	loginInvalidoBody, _ := json.Marshal(authdto.LoginRequest{Email: "hudson@teste.com", Senha: "errada"})
	resp, err = http.Post(server.URL+"/auth/login", "application/json", bytes.NewReader(loginInvalidoBody))
	if err != nil {
		t.Fatalf("erro no login inválido: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("esperava 401 no login inválido, veio %d", resp.StatusCode)
	}
}
