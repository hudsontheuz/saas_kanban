package http

import "github.com/go-chi/chi/v5"

func Register(r chi.Router, handler *AuthHandler) {
	r.Route("/auth", func(r chi.Router) {
		r.Post("/register", handler.Register)
		r.Post("/login", handler.Login)
	})
}
