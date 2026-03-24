package main

import (
	"log"
	"net/http"
	"os"
	"time"

	deliveryhttp "github.com/hudsontheuz/saas_kanban/delivery/http"
	gormdb "github.com/hudsontheuz/saas_kanban/infrastructure/persistence/gorm"
	authusecase "github.com/hudsontheuz/saas_kanban/internal/auth/application/usecase"
	authhttp "github.com/hudsontheuz/saas_kanban/internal/auth/delivery/http"
	authhash "github.com/hudsontheuz/saas_kanban/internal/auth/infrastructure/hash"
	authjwt "github.com/hudsontheuz/saas_kanban/internal/auth/infrastructure/jwt"
	projectrepo "github.com/hudsontheuz/saas_kanban/internal/project/infrastructure/persistence/gorm/repo"
	"github.com/hudsontheuz/saas_kanban/internal/shared/envx"
	taskusecase "github.com/hudsontheuz/saas_kanban/internal/task/application/usecase"
	taskhttp "github.com/hudsontheuz/saas_kanban/internal/task/delivery/http"
	taskrepo "github.com/hudsontheuz/saas_kanban/internal/task/infrastructure/persistence/gorm/repo"
	teamusecase "github.com/hudsontheuz/saas_kanban/internal/team/application/usecase"
	teamhttp "github.com/hudsontheuz/saas_kanban/internal/team/delivery/http"
	teamrepo "github.com/hudsontheuz/saas_kanban/internal/team/infrastructure/persistence/gorm/repo"
	userrepo "github.com/hudsontheuz/saas_kanban/internal/user/infrastructure/persistence/gorm/repo"
)

func main() {
	if err := envx.Load(); err != nil {
		log.Fatalf("env: %v", err)
	}

	db, err := gormdb.Open()
	if err != nil {
		log.Fatalf("db: %v", err)
	}

	repoUsuario := userrepo.NewUserRepo(db)
	repoProjeto := projectrepo.NewProjectRepo(db)
	repoTarefa := taskrepo.NewTaskRepo(db)
	repoTeam := teamrepo.NewTeamRepo(db)

	segredoJWT := os.Getenv("JWT_SECRET")
	emissorJWT := os.Getenv("JWT_ISSUER")

	issuerJWT, err := authjwt.NovoIssuer(segredoJWT, emissorJWT, 24*time.Hour)
	if err != nil {
		log.Fatalf("jwt issuer: %v", err)
	}

	validadorJWT, err := authjwt.NovoValidador(segredoJWT, emissorJWT)
	if err != nil {
		log.Fatalf("jwt validator: %v", err)
	}

	hasher := authhash.NewBcryptHasher()
	ucRegister := authusecase.NovoRegisterUseCase(repoUsuario, hasher, issuerJWT)
	ucLogin := authusecase.NovoLoginUseCase(repoUsuario, hasher, issuerJWT)
	handlerAuth := authhttp.NewAuthHandler(ucRegister, ucLogin)

	casoUsoCriarTeam := teamusecase.NovoCriarTeamUseCase(repoTeam)
	handlerTeam := teamhttp.NewTeamHandler(casoUsoCriarTeam)

	casoUsoCriarTask := taskusecase.NovoCriarTaskUseCase(repoProjeto, repoTarefa)
	casoUsoSelfAssign := taskusecase.NovoSelfAssignTaskUseCase(repoProjeto, repoTarefa)
	casoUsoPausarTask := taskusecase.NovoPausarTaskUseCase(repoProjeto, repoTarefa)
	casoUsoRetomarTask := taskusecase.NovoRetomarTaskUseCase(repoProjeto, repoTarefa)
	casoUsoInReview := taskusecase.NovoMoverParaInReviewUseCase(repoProjeto, repoTarefa)
	casoUsoApprove := taskusecase.NovoAprovarTaskUseCase(repoProjeto, repoTeam, repoTarefa)
	casoUsoReject := taskusecase.NovoReprovarTaskUseCase(repoProjeto, repoTeam, repoTarefa)

	handlerTarefa := taskhttp.NewTaskHandlerWorkflow(
		casoUsoCriarTask,
		casoUsoSelfAssign,
		casoUsoPausarTask,
		casoUsoRetomarTask,
		casoUsoInReview,
		casoUsoApprove,
		casoUsoReject,
	)

	router := deliveryhttp.NewRouter(handlerAuth, handlerTeam, handlerTarefa, validadorJWT)

	porta := os.Getenv("PORT")
	if porta == "" {
		porta = "8080"
	}

	srv := &http.Server{
		Addr:              ":" + porta,
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
	}

	log.Printf("listening on :%s", porta)
	log.Fatal(srv.ListenAndServe())
}
