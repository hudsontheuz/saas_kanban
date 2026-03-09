package task

import (
	"time"

	"github.com/hudsontheuz/saas_kanban/internal/project/domain"
	"github.com/hudsontheuz/saas_kanban/internal/user/domain"
)

func HidratarTask(
	id TaskID,
	projectID project.ProjectID,
	titulo string,
	status StatusTask,
	assignee *user.UserID,
	pausada bool,
	outcome *OutcomeTask,
	deletedAt *time.Time,
	deletedBy *user.UserID,
) *Task {
	return &Task{
		id:        id,
		projectID: projectID,
		titulo:    titulo,
		status:    status,
		assignee:  assignee,
		isPaused:  pausada,
		outcome:   outcome,
		deletedAt: deletedAt,
		deletedBy: deletedBy,
	}
}
