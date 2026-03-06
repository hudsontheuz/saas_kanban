package http

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	authjwt "github.com/hudsontheuz/saas_kanban/internal/auth/infrastructure/jwt"
	taskhttp "github.com/hudsontheuz/saas_kanban/internal/task/delivery/http"
)

// NewRouter monta o router principal, chamando Register() de cada módulo.
func NewRouter(taskHandler *taskhttp.TaskHandler, validadorJWT *authjwt.Validador) http.Handler {
	r := chi.NewRouter()

	// módulos
	taskhttp.Register(r, taskHandler, validadorJWT)

	return r
}
