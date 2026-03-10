package http

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	handlers "github.com/hudsontheuz/saas_kanban/delivery/http/handlers"
)

func NewRouter(taskHandler *handlers.TaskHandler) http.Handler {
	r := chi.NewRouter()
	r.Post("/tasks/{id}/self-assign", taskHandler.SelfAssign)
	return r
}
