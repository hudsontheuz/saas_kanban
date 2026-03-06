package http

import (
	"github.com/go-chi/chi/v5"

	authmiddleware "github.com/hudsontheuz/saas_kanban/internal/auth/delivery/http/middleware"
	authjwt "github.com/hudsontheuz/saas_kanban/internal/auth/infrastructure/jwt"
)

// Register monta as rotas HTTP do módulo Task.
func Register(r chi.Router, handler *TaskHandler, validadorJWT *authjwt.Validador) {
	// MVP incremental: só esta rota por enquanto
	r.With(authmiddleware.AutenticacaoJWT(validadorJWT)).
		Post("/tasks/{id}/self-assign", handler.SelfAssign)
}
