package http

import (
	"github.com/go-chi/chi/v5"

	authmiddleware "github.com/hudsontheuz/saas_kanban/internal/auth/delivery/http/middleware"
	authjwt "github.com/hudsontheuz/saas_kanban/internal/auth/infrastructure/jwt"
)

func Register(r chi.Router, handler *ProjectHandler, validadorJWT *authjwt.Validador) {
	if handler == nil {
		return
	}

	auth := r.With(authmiddleware.AutenticacaoJWT(validadorJWT))

	if handler.create != nil {
		auth.Post("/teams/{id}/projects", handler.Create)
	}

	if handler.getActive != nil {
		auth.Get("/teams/{id}/projects/active", handler.GetActiveByTeam)
	}

	if handler.updateSettings != nil {
		auth.Patch("/projects/{id}/settings", handler.UpdateSettings)
	}

	if handler.close != nil {
		auth.Post("/projects/{id}/close", handler.Close)
	}
}
