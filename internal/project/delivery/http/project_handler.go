package http

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"

	authctx "github.com/hudsontheuz/saas_kanban/internal/auth/delivery/http/authctx"
	projectdto "github.com/hudsontheuz/saas_kanban/internal/project/application/dto"
	projectusecase "github.com/hudsontheuz/saas_kanban/internal/project/application/usecase"
	project "github.com/hudsontheuz/saas_kanban/internal/project/domain"
	shared "github.com/hudsontheuz/saas_kanban/internal/shared/errors"
)

type ProjectHandler struct {
	create         *projectusecase.CriarProjectUseCase
	getActive      *projectusecase.BuscarProjectAtivoUseCase
	updateSettings *projectusecase.AtualizarSettingsProjectUseCase
	close          *projectusecase.FecharProjectUseCase
}

func NewProjectHandler(
	create *projectusecase.CriarProjectUseCase,
	getActive *projectusecase.BuscarProjectAtivoUseCase,
	updateSettings *projectusecase.AtualizarSettingsProjectUseCase,
	close *projectusecase.FecharProjectUseCase,
) *ProjectHandler {
	return &ProjectHandler{
		create:         create,
		getActive:      getActive,
		updateSettings: updateSettings,
		close:          close,
	}
}

type criarProjectBody struct {
	Nome string `json:"nome"`

	PermitirSoltarDoingParaTodo bool `json:"permitir_soltar_doing_para_todo"`
}

type atualizarSettingsBody struct {
	PermitirSoltarDoingParaTodo bool `json:"permitir_soltar_doing_para_todo"`
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

func (h *ProjectHandler) GetActiveByTeam(w http.ResponseWriter, r *http.Request) {
	teamID := chi.URLParam(r, "id")

	resp, err := h.getActive.Executar(projectdto.BuscarProjectAtivoRequest{
		TeamID: teamID,
	})
	if err != nil {
		writeProjectError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, resp)
}

func (h *ProjectHandler) UpdateSettings(w http.ResponseWriter, r *http.Request) {
	projectID := chi.URLParam(r, "id")

	idUsuario, ok := authctx.IDUsuarioDoContexto(r.Context())
	if !ok || idUsuario == "" {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}

	var body atualizarSettingsBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "json inválido"})
		return
	}

	err := h.updateSettings.Executar(projectdto.AtualizarSettingsProjectRequest{
		ProjectID:                   projectID,
		LeaderID:                    string(idUsuario),
		PermitirSoltarDoingParaTodo: body.PermitirSoltarDoingParaTodo,
	})
	if err != nil {
		writeProjectError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{"ok": true})
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

func writeProjectError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, shared.ErrNaoEncontrado):
		writeJSON(w, http.StatusNotFound, map[string]string{"error": err.Error()})

	case errors.Is(err, projectusecase.ErrSomenteLeaderPodeGerenciarProject):
		writeJSON(w, http.StatusForbidden, map[string]string{"error": err.Error()})

	case errors.Is(err, projectusecase.ErrJaExisteProjectAtivo),
		errors.Is(err, project.ErrProjetoFechado):
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
