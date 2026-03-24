package http

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"

	authctx "github.com/hudsontheuz/saas_kanban/internal/auth/delivery/http/authctx"
	projectdto "github.com/hudsontheuz/saas_kanban/internal/project/application/dto"
	projectusecase "github.com/hudsontheuz/saas_kanban/internal/project/application/usecase"
	shared "github.com/hudsontheuz/saas_kanban/internal/shared/errors"
)

type ProjectHandler struct {
	create    *projectusecase.CriarProjectUseCase
	close     *projectusecase.FecharProjectUseCase
	getActive *projectusecase.BuscarProjectAtivoUseCase
}

func NewProjectHandler(
	create *projectusecase.CriarProjectUseCase,
	close *projectusecase.FecharProjectUseCase,
	getActive *projectusecase.BuscarProjectAtivoUseCase,
) *ProjectHandler {
	return &ProjectHandler{
		create:    create,
		close:     close,
		getActive: getActive,
	}
}

type criarProjectBody struct {
	Nome                        string `json:"nome"`
	PermitirSoltarDoingParaTodo bool   `json:"permitir_soltar_doing_para_todo"`
}

func (h *ProjectHandler) Create(w http.ResponseWriter, r *http.Request) {
	teamID := chi.URLParam(r, "id")

	idUsuario, ok := authctx.IDUsuarioDoContexto(r.Context())
	if !ok || idUsuario == "" {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}

	var body criarProjectBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "json inválido"})
		return
	}

	resp, err := h.create.Executar(projectdto.CriarProjectRequest{
		TeamID:                      teamID,
		LeaderID:                    string(idUsuario),
		Nome:                        body.Nome,
		PermitirSoltarDoingParaTodo: body.PermitirSoltarDoingParaTodo,
	})
	if err != nil {
		writeProjectError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, resp)
}

func (h *ProjectHandler) Close(w http.ResponseWriter, r *http.Request) {
	projectID := chi.URLParam(r, "id")

	idUsuario, ok := authctx.IDUsuarioDoContexto(r.Context())
	if !ok || idUsuario == "" {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}

	err := h.close.Executar(projectdto.FecharProjectRequest{
		ProjectID: projectID,
		LeaderID:  string(idUsuario),
	})
	if err != nil {
		writeProjectError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{"ok": true})
}

func (h *ProjectHandler) GetActive(w http.ResponseWriter, r *http.Request) {
	teamID := chi.URLParam(r, "id")

	_, ok := authctx.IDUsuarioDoContexto(r.Context())
	if !ok {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}

	resp, err := h.getActive.Executar(projectdto.BuscarProjectAtivoRequest{TeamID: teamID})
	if err != nil {
		writeProjectError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, resp)
}

func writeProjectError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, shared.ErrNaoEncontrado):
		writeJSON(w, http.StatusNotFound, map[string]string{"error": err.Error()})

	case errors.Is(err, projectusecase.ErrJaExisteProjectAtivo):
		writeJSON(w, http.StatusConflict, map[string]string{"error": err.Error()})

	default:
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
}

func writeJSON(w http.ResponseWriter, code int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(v)
}
