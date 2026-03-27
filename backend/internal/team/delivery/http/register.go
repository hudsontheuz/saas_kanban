package http

import (
	"github.com/go-chi/chi/v5"

	authmiddleware "github.com/hudsontheuz/saas_kanban/internal/auth/delivery/http/middleware"
	authjwt "github.com/hudsontheuz/saas_kanban/internal/auth/infrastructure/jwt"
)

func Register(r chi.Router, handler *TeamHandler, validadorJWT *authjwt.Validador) {
	if handler == nil {
		return
	}

	auth := r.With(authmiddleware.AutenticacaoJWT(validadorJWT))
	auth.Post("/teams", handler.Create)

	if handler.getByID != nil {
		auth.Get("/teams/{id}", handler.GetByID)
	}
}
