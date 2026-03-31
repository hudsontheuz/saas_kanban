package ports

import (
	"github.com/hudsontheuz/saas_kanban/internal/team/domain"
	user "github.com/hudsontheuz/saas_kanban/internal/user/domain"
)

type TeamRepository interface {
	Salvar(t *team.Team) error
	BuscarPorID(id team.TeamID) (*team.Team, error)
	ListarPorUsuarioID(userID user.UserID) ([]*team.Team, error)
	ListarMembros(id team.TeamID) ([]team.Membro, error)
}
