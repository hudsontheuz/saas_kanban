package http

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"

	taskdto "github.com/hudsontheuz/saas_kanban/internal/application/task/dto"
	taskusecase "github.com/hudsontheuz/saas_kanban/internal/application/task/usecase"
	shared "github.com/hudsontheuz/saas_kanban/internal/domain/shared"
)

type TaskHandler struct {
	selfAssign *taskusecase.SelfAssignTaskUseCase
}

func NewTaskHandler(selfAssign *taskusecase.SelfAssignTaskUseCase) *TaskHandler {
	return &TaskHandler{selfAssign: selfAssign}
}

func (h *TaskHandler) SelfAssign(w http.ResponseWriter, r *http.Request) {
	taskID := chi.URLParam(r, "id")
	userID := r.Header.Get("X-User-Id")
	if userID == "" {
		writeJSON(w, 400, map[string]string{"error": "missing X-User-Id header"})
		return
	}

	err := h.selfAssign.Executar(taskdto.SelfAssignRequest{TaskID: taskID, UserID: userID})
	if err != nil {
		if err == shared.ErrNaoEncontrado {
			writeJSON(w, 404, map[string]string{"error": err.Error()})
			return
		}
		if err == taskusecase.ErrLimiteGlobalDoing {
			writeJSON(w, 409, map[string]string{"error": err.Error()})
			return
		}
		writeJSON(w, 400, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, 200, map[string]any{"ok": true})
}

func writeJSON(w http.ResponseWriter, code int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(v)
}
