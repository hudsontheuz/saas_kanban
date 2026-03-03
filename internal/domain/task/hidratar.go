package task

import (
	"time"

	"github.com/hudsontheuz/saas_kanban/internal/domain/project"
	"github.com/hudsontheuz/saas_kanban/internal/domain/team"
)

func HidratarTask(
	id TaskID,
	projectID project.ProjectID,
	titulo string,
	status StatusTask,
	assignee *team.UserID,
	pausada bool,
	outcome *OutcomeTask,
	deletedAt *time.Time,
	deletedBy *team.UserID,
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
