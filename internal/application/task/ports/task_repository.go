package ports

import (
	"github.com/hudsontheuz/saas_kanban/internal/domain/task"
	"github.com/hudsontheuz/saas_kanban/internal/domain/team"
)

type TaskRepository interface {
	Salvar(t *task.Task) error
	BuscarPorID(id task.TaskID) (*task.Task, error)
	ExisteDoingAtivaParaUser(userID team.UserID) (bool, error)
}
