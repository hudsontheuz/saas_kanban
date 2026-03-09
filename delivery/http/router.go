package http

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	authhttp "github.com/hudsontheuz/saas_kanban/internal/auth/delivery/http"
	authjwt "github.com/hudsontheuz/saas_kanban/internal/auth/infrastructure/jwt"
	taskhttp "github.com/hudsontheuz/saas_kanban/internal/task/delivery/http"
)

func NewRouter(authHandler *authhttp.AuthHandler, taskHandler *taskhttp.TaskHandler, validadorJWT *authjwt.Validador) http.Handler {
	r := chi.NewRouter()

	authhttp.Register(r, authHandler)
	taskhttp.Register(r, taskHandler, validadorJWT)

	return r
}
