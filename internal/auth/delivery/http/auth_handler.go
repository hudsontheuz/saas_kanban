package http

import (
	"encoding/json"
	"errors"
	"net/http"

	authdto "github.com/hudsontheuz/saas_kanban/internal/auth/application/dto"
	authusecase "github.com/hudsontheuz/saas_kanban/internal/auth/application/usecase"
	auth "github.com/hudsontheuz/saas_kanban/internal/auth/domain"
)

type AuthHandler struct {
	register *authusecase.RegisterUseCase
	login    *authusecase.LoginUseCase
}

func NewAuthHandler(register *authusecase.RegisterUseCase, login *authusecase.LoginUseCase) *AuthHandler {
	return &AuthHandler{register: register, login: login}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req authdto.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "json inválido"})
		return
	}

	resp, err := h.register.Executar(req)
	if err != nil {
		h.handleAuthError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, resp)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req authdto.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "json inválido"})
		return
	}

	resp, err := h.login.Executar(req)
	if err != nil {
		h.handleAuthError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, resp)
}

func (h *AuthHandler) handleAuthError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, auth.ErrEmailJaCadastrado):
		writeJSON(w, http.StatusConflict, map[string]string{"error": err.Error()})
	case errors.Is(err, auth.ErrCredenciaisInvalidas):
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": err.Error()})
	default:
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
}

func writeJSON(w http.ResponseWriter, code int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(v)
}
