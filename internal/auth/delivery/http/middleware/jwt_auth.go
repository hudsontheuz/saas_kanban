package middleware

import (
	"net/http"
	"strings"

	"github.com/hudsontheuz/saas_kanban/internal/auth/delivery/http/authctx"
	"github.com/hudsontheuz/saas_kanban/internal/auth/infrastructure/jwt"
)

func AutenticacaoJWT(validador *jwt.Validador) func(http.Handler) http.Handler {
	return func(proximo http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cabecalhoAuth := r.Header.Get("Authorization")
			if cabecalhoAuth == "" {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}

			partes := strings.SplitN(cabecalhoAuth, " ", 2)
			if len(partes) != 2 || !strings.EqualFold(partes[0], "Bearer") || strings.TrimSpace(partes[1]) == "" {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}

			idUsuario, err := validador.ValidarEObterIDUsuario(strings.TrimSpace(partes[1]))
			if err != nil {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}

			ctx := authctx.ComIDUsuario(r.Context(), idUsuario)
			proximo.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
