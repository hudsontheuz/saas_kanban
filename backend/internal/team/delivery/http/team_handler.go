package http

import (
	"encoding/json"
	"net/http"

	authctx "github.com/hudsontheuz/saas_kanban/internal/auth/delivery/http/authctx"
	teamdto "github.com/hudsontheuz/saas_kanban/internal/team/application/dto"
	teamusecase "github.com/hudsontheuz/saas_kanban/internal/team/application/usecase"
)

type TeamHandler struct {
	create *teamusecase.CriarTeamUseCase
}

func NewTeamHandler(create *teamusecase.CriarTeamUseCase) *TeamHandler {
	return &TeamHandler{create: create}
}

type criarTeamBody struct {
	Nome string `json:"nome"`
}

func (h *TeamHandler) Create(w http.ResponseWriter, r *http.Request) {
	idUsuario, ok := authctx.IDUsuarioDoContexto(r.Context())
	if !ok || idUsuario == "" {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}

	var body criarTeamBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "json inválido"})
		return
	}

	resp, err := h.create.Executar(teamdto.CriarTeamRequest{
		Nome:     body.Nome,
		LeaderID: string(idUsuario),
	})
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusCreated, resp)
}

func writeJSON(w http.ResponseWriter, code int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(v)
}
