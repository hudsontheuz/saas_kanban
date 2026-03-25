package ports

import (
	project "github.com/hudsontheuz/saas_kanban/internal/project/domain"
	task "github.com/hudsontheuz/saas_kanban/internal/task/domain"
	user "github.com/hudsontheuz/saas_kanban/internal/user/domain"
)

type TaskRepository interface {
	Salvar(t *task.Task) error
	BuscarPorID(id task.TaskID) (*task.Task, error)
	ListarPorProjectID(projectID project.ProjectID, status *task.StatusTask) ([]*task.Task, error)
	ExisteDoingAtivaParaUser(userID user.UserID) (bool, error)
}
