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
	projectrepo "github.com/hudsontheuz/saas_kanban/internal/project/infrastructure/persistence/gorm/repo"
	taskusecase "github.com/hudsontheuz/saas_kanban/internal/task/application/usecase"
	taskhttp "github.com/hudsontheuz/saas_kanban/internal/task/delivery/http"
	taskrepo "github.com/hudsontheuz/saas_kanban/internal/task/infrastructure/persistence/gorm/repo"
	teamdto "github.com/hudsontheuz/saas_kanban/internal/team/application/dto"
	teamusecase "github.com/hudsontheuz/saas_kanban/internal/team/application/usecase"
	teamhttp "github.com/hudsontheuz/saas_kanban/internal/team/delivery/http"
	teamrepo "github.com/hudsontheuz/saas_kanban/internal/team/infrastructure/persistence/gorm/repo"
	userrepo "github.com/hudsontheuz/saas_kanban/internal/user/infrastructure/persistence/gorm/repo"
)

func TestHTTP_Team_Create(t *testing.T) {
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

	handlerTask := taskhttp.NewTaskHandler(taskusecase.NovoSelfAssignTaskUseCase(repoProjeto, repoTarefa))

	ucCriarTeam := teamusecase.NovoCriarTeamUseCase(repoTeam)
	handlerTeam := teamhttp.NewTeamHandler(ucCriarTeam)

	router := deliveryhttp.NewRouter(handlerAuth, handlerTeam, handlerTask, validador)
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

	body, _ := json.Marshal(map[string]any{
		"nome": "Equipe Alpha",
	})

	req, _ := http.NewRequest(http.MethodPost, server.URL+"/teams", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+loginResp.Token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("erro criando team: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("esperava 201, veio %d", resp.StatusCode)
	}

	var teamResp teamdto.CriarTeamResponse
	if err := json.NewDecoder(resp.Body).Decode(&teamResp); err != nil {
		t.Fatalf("erro lendo response: %v", err)
	}

	if teamResp.TeamID == "" {
		t.Fatalf("esperava team_id preenchido")
	}
}
