package http

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"

	authctx "github.com/hudsontheuz/saas_kanban/internal/auth/delivery/http/authctx"
	teamdto "github.com/hudsontheuz/saas_kanban/internal/team/application/dto"
	teamusecase "github.com/hudsontheuz/saas_kanban/internal/team/application/usecase"
)

type TeamHandler struct {
	create      *teamusecase.CriarTeamUseCase
	getByID     *teamusecase.BuscarTeamUseCase
	listMyTeams *teamusecase.ListarMinhasTeamsUseCase
}

func NewTeamHandler(
	create *teamusecase.CriarTeamUseCase,
	getByID *teamusecase.BuscarTeamUseCase,
	listMyTeams *teamusecase.ListarMinhasTeamsUseCase,
) *TeamHandler {
	return &TeamHandler{create: create, getByID: getByID, listMyTeams: listMyTeams}
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

func (h *TeamHandler) GetMyTeams(w http.ResponseWriter, r *http.Request) {
	if h.listMyTeams == nil {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "not found"})
		return
	}

	idUsuario, ok := authctx.IDUsuarioDoContexto(r.Context())
	if !ok || idUsuario == "" {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}

	resp, err := h.listMyTeams.Executar(teamdto.ListarMinhasTeamsRequest{UserID: string(idUsuario)})
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, resp)
}

func (h *TeamHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	if h.getByID == nil {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "not found"})
		return
	}

	teamID := chi.URLParam(r, "id")
	if teamID == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "team_id obrigatório"})
		return
	}

	resp, err := h.getByID.Executar(teamdto.BuscarTeamRequest{
		TeamID: teamID,
	})
	if err != nil {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, resp)
}

func writeJSON(w http.ResponseWriter, code int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(v)
}
