package memory

import (
	"strconv"
	"sync"

	"github.com/hudsontheuz/saas_kanban/internal/domain/task"
	"github.com/hudsontheuz/saas_kanban/internal/domain/team"
)

type TaskRepoEmMemoria struct {
	mu     sync.RWMutex
	dados  map[task.TaskID]*task.Task
	nextID int64
}

func NovoTaskRepoEmMemoria() *TaskRepoEmMemoria {
	return &TaskRepoEmMemoria{
		dados:  map[task.TaskID]*task.Task{},
		nextID: 1,
	}
}

func (r *TaskRepoEmMemoria) Salvar(tk *task.Task) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if tk == nil {
		return ErrNaoEncontrado
	}

	// Se vier sem ID, simula DB gerando ID
	if string(tk.ID()) == "" {
		id := task.TaskID(strconv.FormatInt(r.nextID, 10))
		r.nextID++
		_ = tk.DefinirID(id)
	}

	r.dados[tk.ID()] = tk
	return nil
}

func (r *TaskRepoEmMemoria) BuscarPorID(id task.TaskID) (*task.Task, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	tk, ok := r.dados[id]
	if !ok {
		return nil, ErrNaoEncontrado
	}
	return tk, nil
}

func (r *TaskRepoEmMemoria) ExisteDoingAtivaParaUser(userID team.UserID) (bool, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, tk := range r.dados {
		a := tk.Assignee()
		if a == nil {
			continue
		}
		if *a == userID && tk.Status() == task.Doing && !tk.IsPaused() && tk.DeletedAt() == nil {
			return true, nil
		}
	}
	return false, nil
}
