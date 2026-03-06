package ports

import (
	"github.com/hudsontheuz/saas_kanban/internal/project/domain"
	"github.com/hudsontheuz/saas_kanban/internal/team/domain"
)

type ProjectRepository interface {
	Salvar(p *project.Project) error
	BuscarPorID(id project.ProjectID) (*project.Project, error)
	BuscarAtivoPorTeamID(teamID team.TeamID) (*project.Project, error)
}
