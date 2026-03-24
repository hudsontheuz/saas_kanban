package tests

import (
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
	projectdomain "github.com/hudsontheuz/saas_kanban/internal/project/domain"
	projectrepo "github.com/hudsontheuz/saas_kanban/internal/project/infrastructure/persistence/gorm/repo"
	taskdto "github.com/hudsontheuz/saas_kanban/internal/task/application/dto"
	taskusecase "github.com/hudsontheuz/saas_kanban/internal/task/application/usecase"
	taskhttp "github.com/hudsontheuz/saas_kanban/internal/task/delivery/http"
	taskdomain "github.com/hudsontheuz/saas_kanban/internal/task/domain"
	taskrepo "github.com/hudsontheuz/saas_kanban/internal/task/infrastructure/persistence/gorm/repo"
	userrepo "github.com/hudsontheuz/saas_kanban/internal/user/infrastructure/persistence/gorm/repo"
)

func TestHTTP_Task_ListByProject(t *testing.T) {
	db := openTestDB(t)

	repoUsuario := userrepo.NewUserRepo(db)
	repoProjeto := projectrepo.NewProjectRepo(db)
	repoTarefa := taskrepo.NewTaskRepo(db)

	leaderID := seedUser(t, db, "Leader Lista Tasks")
	userID := seedUser(t, db, "User Lista Tasks")
	teamID := seedTeam(t, db, "Team Lista Tasks", leaderID)
	projectID := seedProject(t, db, teamID, "Projeto Lista Tasks", projectdomain.ConfiguracoesProject{})

	taskTodoID := seedTask(t, db, projectID, "Task TODO")
	taskDoingID := seedTask(t, db, projectID, "Task DOING")

	tkDoing, err := repoTarefa.BuscarPorID(taskdomain.TaskID(taskDoingID))
	if err != nil {
		t.Fatalf("erro buscando task doing: %v", err)
	}
	if err := tkDoing.SelfAssign(userID); err != nil {
		t.Fatalf("erro movendo task para doing: %v", err)
	}
	if err := repoTarefa.Salvar(tkDoing); err != nil {
		t.Fatalf("erro salvando task doing: %v", err)
	}

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

	router := deliveryhttp.NewRouter(handlerAuth, nil, nil, handlerTask, validador)
	server := httptest.NewServer(router)
	defer server.Close()

	token := gerarJWT(t, segredo, emissor, string(userID), time.Now().Add(1*time.Hour))

	req, _ := http.NewRequest(http.MethodGet, server.URL+"/projects/"+string(projectID)+"/tasks", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("erro listando tasks: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("esperava 200 na listagem de tasks, veio %d", resp.StatusCode)
	}

	var listResp taskdto.ListarTasksPorProjectResponse
	if err := json.NewDecoder(resp.Body).Decode(&listResp); err != nil {
		t.Fatalf("erro lendo response da listagem: %v", err)
	}

	if len(listResp.Items) != 2 {
		t.Fatalf("esperava 2 tasks na listagem, veio %d", len(listResp.Items))
	}

	if listResp.Items[0].TaskID != string(taskTodoID) {
		t.Fatalf("esperava primeira task %s, veio %s", taskTodoID, listResp.Items[0].TaskID)
	}

	req, _ = http.NewRequest(http.MethodGet, server.URL+"/projects/"+string(projectID)+"/tasks?status=DOING", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("erro listando tasks filtradas: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("esperava 200 na listagem filtrada, veio %d", resp.StatusCode)
	}

	if err := json.NewDecoder(resp.Body).Decode(&listResp); err != nil {
		t.Fatalf("erro lendo response filtrada: %v", err)
	}

	if len(listResp.Items) != 1 {
		t.Fatalf("esperava 1 task com filtro DOING, veio %d", len(listResp.Items))
	}
	if listResp.Items[0].TaskID != string(taskDoingID) {
		t.Fatalf("esperava task doing %s, veio %s", taskDoingID, listResp.Items[0].TaskID)
	}
}
