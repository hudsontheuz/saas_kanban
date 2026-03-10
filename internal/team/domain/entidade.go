package team

import (
	"strings"

	shared "github.com/hudsontheuz/saas_kanban/internal/shared/errors"
	"github.com/hudsontheuz/saas_kanban/internal/user/domain"
)

type Team struct {
	id       TeamID
	nome     string
	leaderID user.UserID
}

func NovaTeam(nome string, leaderID user.UserID) (*Team, error) {
	nome = strings.TrimSpace(nome)
	if nome == "" {
		return nil, ErrNomeObrigatorio
	}
	if leaderID == "" {
		return nil, ErrLeaderObrigatorio
	}

	return &Team{
		id:       "",
		nome:     nome,
		leaderID: leaderID,
	}, nil
}

func (t *Team) ID() TeamID                  { return t.id }
func (t *Team) Nome() string                { return t.nome }
func (t *Team) LeaderID() user.UserID       { return t.leaderID }
func (t *Team) EhLeader(u user.UserID) bool { return t.leaderID == u }

func (t *Team) DefinirID(id TeamID) error {
	if id == "" {
		return shared.ErrIDInvalido
	}
	t.id = id
	return nil
}
