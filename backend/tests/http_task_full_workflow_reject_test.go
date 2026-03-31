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
	project "github.com/hudsontheuz/saas_kanban/internal/project/domain"
	projectrepo "github.com/hudsontheuz/saas_kanban/internal/project/infrastructure/persistence/gorm/repo"
	taskdto "github.com/hudsontheuz/saas_kanban/internal/task/application/dto"
	taskusecase "github.com/hudsontheuz/saas_kanban/internal/task/application/usecase"
	taskhttp "github.com/hudsontheuz/saas_kanban/internal/task/delivery/http"
	task "github.com/hudsontheuz/saas_kanban/internal/task/domain"
	taskrepo "github.com/hudsontheuz/saas_kanban/internal/task/infrastructure/persistence/gorm/repo"
	teamrepo "github.com/hudsontheuz/saas_kanban/internal/team/infrastructure/persistence/gorm/repo"
	userrepo "github.com/hudsontheuz/saas_kanban/internal/user/infrastructure/persistence/gorm/repo"
)

func TestHTTP_TaskFullWorkflow_Reject(t *testing.T) {
	db := openTestDB(t)

	repoUsuario := userrepo.NewUserRepo(db)
	repoProjeto := projectrepo.NewProjectRepo(db)
	repoTarefa := taskrepo.NewTaskRepo(db)
	repoTeam := teamrepo.NewTeamRepo(db)

	leaderID := seedUser(t, db, "Leader Reject")
	userID := seedUser(t, db, "Executor Reject")
	teamID := seedTeam(t, db, "Team Reject", leaderID)
	projectID := seedProject(t, db, teamID, "Projeto Reject", project.ConfiguracoesProject{})

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

	ucCriarTask := taskusecase.NovoCriarTaskUseCase(repoProjeto, repoTarefa)
	ucSelfAssign := taskusecase.NovoSelfAssignTaskUseCase(repoProjeto, repoTarefa)
	ucPausarTask := taskusecase.NovoPausarTaskUseCase(repoProjeto, repoTarefa)
	ucRetomarTask := taskusecase.NovoRetomarTaskUseCase(repoProjeto, repoTarefa)
	ucInReview := taskusecase.NovoMoverParaInReviewUseCase(repoProjeto, repoTarefa)
	ucApprove := taskusecase.NovoAprovarTaskUseCase(repoProjeto, repoTeam, repoTarefa)
	ucReject := taskusecase.NovoReprovarTaskUseCase(repoProjeto, repoTeam, repoTarefa)

	handlerTask := taskhttp.NewTaskHandlerWorkflow(
		ucCriarTask,
		ucSelfAssign,
		ucPausarTask,
		ucRetomarTask,
		ucInReview,
		ucApprove,
		ucReject,
	)

	router := deliveryhttp.NewRouter(handlerAuth, nil, nil, handlerTask, validador)
	server := httptest.NewServer(router)
	defer server.Close()

	userToken := gerarJWT(t, segredo, emissor, string(userID), time.Now().Add(1*time.Hour))
	leaderToken := gerarJWT(t, segredo, emissor, string(leaderID), time.Now().Add(1*time.Hour))

	createBody, _ := json.Marshal(map[string]string{"titulo": "Task completa reject", "descricao": "Implementar fluxo completo de reprovação"})

	req, _ := http.NewRequest(http.MethodPost, server.URL+"/projects/"+string(projectID)+"/tasks", bytes.NewReader(createBody))
	req.Header.Set("Authorization", "Bearer "+userToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("erro no create task: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("esperava 201 no create task, veio %d", resp.StatusCode)
	}

	var createResp taskdto.CriarTaskResponse
	if err := json.NewDecoder(resp.Body).Decode(&createResp); err != nil {
		t.Fatalf("erro lendo create response: %v", err)
	}

	req, _ = http.NewRequest(http.MethodPost, server.URL+"/tasks/"+createResp.TaskID+"/self-assign", nil)
	req.Header.Set("Authorization", "Bearer "+userToken)
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("erro no self-assign: %v", err)
	}
	defer resp.Body.Close()

	inReviewBody, _ := json.Marshal(map[string]string{"comentario_entrega": "Implementei o fluxo inicial, mas ainda falta revisar um caso de borda."})
	req, _ = http.NewRequest(http.MethodPost, server.URL+"/tasks/"+createResp.TaskID+"/in-review", bytes.NewReader(inReviewBody))
	req.Header.Set("Authorization", "Bearer "+userToken)
	req.Header.Set("Content-Type", "application/json")
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("erro no in-review: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("esperava 200 no in-review, veio %d", resp.StatusCode)
	}

	rejectBody, _ := json.Marshal(map[string]string{"motivo": "Faltou validar o fluxo quando volta para To Do"})
	req, _ = http.NewRequest(http.MethodPost, server.URL+"/tasks/"+createResp.TaskID+"/reject", bytes.NewReader(rejectBody))
	req.Header.Set("Authorization", "Bearer "+leaderToken)
	req.Header.Set("Content-Type", "application/json")
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("erro no reject: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("esperava 200 no reject, veio %d", resp.StatusCode)
	}

	tk, err := repoTarefa.BuscarPorID(task.TaskID(createResp.TaskID))
	if err != nil {
		t.Fatalf("erro buscando task após reject: %v", err)
	}

	if tk.Status() != task.ToDo {
		t.Fatalf("esperava task voltando para TODO, veio %s", tk.Status())
	}
	if tk.ComentarioReview() == "" {
		t.Fatalf("esperava comentário de review salvo")
	}
	if tk.ComentarioEntrega() == "" {
		t.Fatalf("esperava comentário de entrega salvo")
	}
}
