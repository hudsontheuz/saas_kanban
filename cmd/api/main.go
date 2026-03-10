package main

import (
	"log"
	"net/http"
	"os"
	"time"

	deliveryhttp "github.com/hudsontheuz/saas_kanban/delivery/http"
	gormdb "github.com/hudsontheuz/saas_kanban/infrastructure/persistence/gorm"
	authjwt "github.com/hudsontheuz/saas_kanban/internal/auth/infrastructure/jwt"
	projectrepo "github.com/hudsontheuz/saas_kanban/internal/project/infrastructure/persistence/gorm/repo"
	taskusecase "github.com/hudsontheuz/saas_kanban/internal/task/application/usecase"
	taskhttp "github.com/hudsontheuz/saas_kanban/internal/task/delivery/http"
	taskrepo "github.com/hudsontheuz/saas_kanban/internal/task/infrastructure/persistence/gorm/repo"
)

func main() {
	db, err := gormdb.Open()
	if err != nil {
		log.Fatalf("db: %v", err)
	}

	repoProjeto := projectrepo.NewProjectRepo(db)
	repoTarefa := taskrepo.NewTaskRepo(db)

	casoUsoSelfAssign := taskusecase.NovoSelfAssignTaskUseCase(repoProjeto, repoTarefa)
	handlerTarefa := taskhttp.NewTaskHandler(casoUsoSelfAssign)

	segredoJWT := os.Getenv("JWT_SECRET")
	emissorJWT := os.Getenv("JWT_ISSUER") // opcional; default no validador = "saas_kanban"

	validadorJWT, err := authjwt.NovoValidador(segredoJWT, emissorJWT)
	if err != nil {
		log.Fatalf("jwt: %v", err)
	}

	router := deliveryhttp.NewRouter(handlerTarefa, validadorJWT)

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
