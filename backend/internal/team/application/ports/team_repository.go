package ports

import "github.com/hudsontheuz/saas_kanban/internal/team/domain"

type TeamRepository interface {
	Salvar(t *team.Team) error
	BuscarPorID(id team.TeamID) (*team.Team, error)
}
