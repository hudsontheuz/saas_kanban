package http

import (
	"github.com/go-chi/chi/v5"

	authmiddleware "github.com/hudsontheuz/saas_kanban/internal/auth/delivery/http/middleware"
	authjwt "github.com/hudsontheuz/saas_kanban/internal/auth/infrastructure/jwt"
)

func Register(r chi.Router, handler *TaskHandler, validadorJWT *authjwt.Validador) {
	if handler == nil {
		return
	}

	auth := r.With(authmiddleware.AutenticacaoJWT(validadorJWT))

	if handler.create != nil {
		auth.Post("/projects/{id}/tasks", handler.Create)
	}

	if handler.listByProject != nil {
		auth.Get("/projects/{id}/tasks", handler.ListByProject)
	}

	if handler.selfAssign != nil {
		auth.Post("/tasks/{id}/self-assign", handler.SelfAssign)
	}

	if handler.pause != nil {
		auth.Post("/tasks/{id}/pause", handler.Pause)
	}

	if handler.resume != nil {
		auth.Post("/tasks/{id}/resume", handler.Resume)
	}

	if handler.inReview != nil {
		auth.Post("/tasks/{id}/in-review", handler.InReview)
	}

	if handler.approve != nil {
		auth.Post("/tasks/{id}/approve", handler.Approve)
	}

	if handler.reject != nil {
		auth.Post("/tasks/{id}/reject", handler.Reject)
	}
}
