package http

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"

	authctx "github.com/hudsontheuz/saas_kanban/internal/auth/delivery/http/authctx"
	project "github.com/hudsontheuz/saas_kanban/internal/project/domain"
	shared "github.com/hudsontheuz/saas_kanban/internal/shared/errors"
	taskdto "github.com/hudsontheuz/saas_kanban/internal/task/application/dto"
	taskusecase "github.com/hudsontheuz/saas_kanban/internal/task/application/usecase"
	task "github.com/hudsontheuz/saas_kanban/internal/task/domain"
)

type TaskHandler struct {
	create        *taskusecase.CriarTaskUseCase
	listByProject *taskusecase.ListarTasksPorProjectUseCase
	selfAssign    *taskusecase.SelfAssignTaskUseCase
	pause         *taskusecase.PausarTaskUseCase
	resume        *taskusecase.RetomarTaskUseCase
	inReview      *taskusecase.MoverParaInReviewUseCase
	approve       *taskusecase.AprovarTaskUseCase
	reject        *taskusecase.ReprovarTaskUseCase
}

func NewTaskHandler(selfAssign *taskusecase.SelfAssignTaskUseCase) *TaskHandler {
	return &TaskHandler{
		selfAssign: selfAssign,
	}
}

func NewTaskHandlerWorkflow(
	create *taskusecase.CriarTaskUseCase,
	selfAssign *taskusecase.SelfAssignTaskUseCase,
	pause *taskusecase.PausarTaskUseCase,
	resume *taskusecase.RetomarTaskUseCase,
	inReview *taskusecase.MoverParaInReviewUseCase,
	approve *taskusecase.AprovarTaskUseCase,
	reject *taskusecase.ReprovarTaskUseCase,
) *TaskHandler {
	return &TaskHandler{
		create:     create,
		selfAssign: selfAssign,
		pause:      pause,
		resume:     resume,
		inReview:   inReview,
		approve:    approve,
		reject:     reject,
	}
}

func NewTaskHandlerWorkflowWithList(
	create *taskusecase.CriarTaskUseCase,
	listByProject *taskusecase.ListarTasksPorProjectUseCase,
	selfAssign *taskusecase.SelfAssignTaskUseCase,
	pause *taskusecase.PausarTaskUseCase,
	resume *taskusecase.RetomarTaskUseCase,
	inReview *taskusecase.MoverParaInReviewUseCase,
	approve *taskusecase.AprovarTaskUseCase,
	reject *taskusecase.ReprovarTaskUseCase,
) *TaskHandler {
	return &TaskHandler{
		create:        create,
		listByProject: listByProject,
		selfAssign:    selfAssign,
		pause:         pause,
		resume:        resume,
		inReview:      inReview,
		approve:       approve,
		reject:        reject,
	}
}

type criarTaskBody struct {
	Titulo    string `json:"titulo"`
	Descricao string `json:"descricao"`
}

type reprovarTaskBody struct {
	Motivo string `json:"motivo"`
}

type moverParaInReviewBody struct {
	ComentarioEntrega string `json:"comentario_entrega"`
}

func (h *TaskHandler) Create(w http.ResponseWriter, r *http.Request) {
	projectID := chi.URLParam(r, "id")

	idUsuario, ok := authctx.IDUsuarioDoContexto(r.Context())
	if !ok || idUsuario == "" {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}

	var body criarTaskBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "json inválido"})
		return
	}

	resp, err := h.create.Executar(taskdto.CriarTaskRequest{
		ProjectID: projectID,
		Titulo:    body.Titulo,
		Descricao: body.Descricao,
		CriadorID: string(idUsuario),
	})
	if err != nil {
		writeTaskError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, resp)
}

func (h *TaskHandler) ListByProject(w http.ResponseWriter, r *http.Request) {
	projectID := chi.URLParam(r, "id")

	_, ok := authctx.IDUsuarioDoContexto(r.Context())
	if !ok {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}

	var status *string
	if s := r.URL.Query().Get("status"); s != "" {
		status = &s
	}

	resp, err := h.listByProject.Executar(taskdto.ListarTasksPorProjectRequest{
		ProjectID: projectID,
		Status:    status,
	})
	if err != nil {
		writeTaskError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, resp)
}

func (h *TaskHandler) Approve(w http.ResponseWriter, r *http.Request) {
	taskID := chi.URLParam(r, "id")

	idUsuario, ok := authctx.IDUsuarioDoContexto(r.Context())
	if !ok || idUsuario == "" {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}

	err := h.approve.Executar(taskdto.AprovarTaskRequest{
		TaskID:   taskID,
		LeaderID: string(idUsuario),
	})
	if err != nil {
		writeTaskError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{"ok": true})
}

func (h *TaskHandler) SelfAssign(w http.ResponseWriter, r *http.Request) {
	taskID := chi.URLParam(r, "id")

	idUsuario, ok := authctx.IDUsuarioDoContexto(r.Context())
	if !ok || idUsuario == "" {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}

	err := h.selfAssign.Executar(taskdto.SelfAssignRequest{
		TaskID: taskID,
		UserID: string(idUsuario),
	})
	if err != nil {
		writeTaskError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{"ok": true})
}

func (h *TaskHandler) Pause(w http.ResponseWriter, r *http.Request) {
	taskID := chi.URLParam(r, "id")

	idUsuario, ok := authctx.IDUsuarioDoContexto(r.Context())
	if !ok || idUsuario == "" {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}

	err := h.pause.Executar(taskdto.PausarTaskRequest{
		TaskID: taskID,
		UserID: string(idUsuario),
	})
	if err != nil {
		writeTaskError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{"ok": true})
}

func (h *TaskHandler) Resume(w http.ResponseWriter, r *http.Request) {
	taskID := chi.URLParam(r, "id")

	idUsuario, ok := authctx.IDUsuarioDoContexto(r.Context())
	if !ok || idUsuario == "" {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}

	err := h.resume.Executar(taskdto.RetomarTaskRequest{
		TaskID: taskID,
		UserID: string(idUsuario),
	})
	if err != nil {
		writeTaskError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{"ok": true})
}

func (h *TaskHandler) InReview(w http.ResponseWriter, r *http.Request) {
	taskID := chi.URLParam(r, "id")

	idUsuario, ok := authctx.IDUsuarioDoContexto(r.Context())
	if !ok || idUsuario == "" {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}

	var body moverParaInReviewBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "json inválido"})
		return
	}

	err := h.inReview.Executar(taskdto.MoverParaInReviewRequest{
		TaskID:            taskID,
		UserID:            string(idUsuario),
		ComentarioEntrega: body.ComentarioEntrega,
	})
	if err != nil {
		writeTaskError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{"ok": true})
}

func (h *TaskHandler) Reject(w http.ResponseWriter, r *http.Request) {
	taskID := chi.URLParam(r, "id")

	idUsuario, ok := authctx.IDUsuarioDoContexto(r.Context())
	if !ok || idUsuario == "" {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}

	var body reprovarTaskBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "json inválido"})
		return
	}

	err := h.reject.Executar(taskdto.ReprovarTaskRequest{
		TaskID:   taskID,
		LeaderID: string(idUsuario),
		Motivo:   body.Motivo,
	})
	if err != nil {
		writeTaskError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{"ok": true})
}

func writeTaskError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, shared.ErrNaoEncontrado):
		writeJSON(w, http.StatusNotFound, map[string]string{"error": err.Error()})

	case errors.Is(err, taskusecase.ErrLimiteGlobalDoing),
		errors.Is(err, project.ErrProjetoFechado),
		errors.Is(err, task.ErrTransicaoInvalida):
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
