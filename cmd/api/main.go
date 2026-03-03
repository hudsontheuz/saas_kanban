package main

import (
	"log"
	"net/http"
	"os"
	"time"

	delivery "github.com/hudsontheuz/saas_kanban/delivery/http"
	handlers "github.com/hudsontheuz/saas_kanban/delivery/http/handlers"
	gormdb "github.com/hudsontheuz/saas_kanban/infrastructure/persistence/gorm"
	"github.com/hudsontheuz/saas_kanban/infrastructure/persistence/gorm/repo"
	taskusecase "github.com/hudsontheuz/saas_kanban/internal/application/task/usecase"
)

func main() {
	db, err := gormdb.Open()
	if err != nil {
		log.Fatalf("db: %v", err)
	}

	projectRepo := repo.NewProjectRepo(db)
	taskRepo := repo.NewTaskRepo(db)

	selfAssignUC := taskusecase.NovoSelfAssignTaskUseCase(projectRepo, taskRepo)

	taskHandler := handlers.NewTaskHandler(selfAssignUC)
	router := delivery.NewRouter(taskHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	srv := &http.Server{
		Addr:              ":" + port,
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
	}

	log.Printf("listening on :%s", port)
	log.Fatal(srv.ListenAndServe())
}
