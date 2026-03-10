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
	taskusecase "github.com/hudsontheuz/saas_kanban/internal/task/application/usecase"
	taskhttp "github.com/hudsontheuz/saas_kanban/internal/task/delivery/http"
	taskrepo "github.com/hudsontheuz/saas_kanban/internal/task/infrastructure/persistence/gorm/repo"
	userrepo "github.com/hudsontheuz/saas_kanban/internal/user/infrastructure/persistence/gorm/repo"
)

func main() {
	db, err := gormdb.Open()
	if err != nil {
		log.Fatalf("db: %v", err)
	}

	repoUsuario := userrepo.NewUserRepo(db)
	repoProjeto := projectrepo.NewProjectRepo(db)
	repoTarefa := taskrepo.NewTaskRepo(db)

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

	casoUsoSelfAssign := taskusecase.NovoSelfAssignTaskUseCase(repoProjeto, repoTarefa)
	handlerTarefa := taskhttp.NewTaskHandler(casoUsoSelfAssign)

	router := deliveryhttp.NewRouter(handlerAuth, handlerTarefa, validadorJWT)

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
