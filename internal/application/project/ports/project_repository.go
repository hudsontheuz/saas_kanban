package ports

import (
	"github.com/hudsontheuz/saas_kanban/internal/domain/project"
	"github.com/hudsontheuz/saas_kanban/internal/domain/team"
)

type ProjectRepository interface {
	Salvar(p *project.Project) error
	BuscarPorID(id project.ProjectID) (*project.Project, error)
	BuscarAtivoPorTeamID(teamID team.TeamID) (*project.Project, error)
}
