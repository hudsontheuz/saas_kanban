package ports

import (
	"github.com/hudsontheuz/saas_kanban/internal/task/domain"
	"github.com/hudsontheuz/saas_kanban/internal/team/domain"
)

type TaskRepository interface {
	Salvar(t *task.Task) error
	BuscarPorID(id task.TaskID) (*task.Task, error)
	ExisteDoingAtivaParaUser(userID team.UserID) (bool, error)
}
