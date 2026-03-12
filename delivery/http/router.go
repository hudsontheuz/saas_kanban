package http

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	authhttp "github.com/hudsontheuz/saas_kanban/internal/auth/delivery/http"
	authjwt "github.com/hudsontheuz/saas_kanban/internal/auth/infrastructure/jwt"
	taskhttp "github.com/hudsontheuz/saas_kanban/internal/task/delivery/http"
	teamhttp "github.com/hudsontheuz/saas_kanban/internal/team/delivery/http"
)

func NewRouter(
	authHandler *authhttp.AuthHandler,
	teamHandler *teamhttp.TeamHandler,
	taskHandler *taskhttp.TaskHandler,
	validadorJWT *authjwt.Validador,
) http.Handler {
	r := chi.NewRouter()

	authhttp.Register(r, authHandler)
	teamhttp.Register(r, teamHandler, validadorJWT)
	taskhttp.Register(r, taskHandler, validadorJWT)

	return r
}
