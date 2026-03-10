package project

import (
	"time"

	"github.com/hudsontheuz/saas_kanban/internal/domain/team"
)

func HidratarProject(
	id ProjectID,
	teamID team.TeamID,
	nome string,
	status StatusProject,
	settings ConfiguracoesProject,
	fechadoEm *time.Time,
) *Project {
	return &Project{
		id:        id,
		teamID:    teamID,
		nome:      nome,
		status:    status,
		settings:  settings,
		fechadoEm: fechadoEm,
	}
}
