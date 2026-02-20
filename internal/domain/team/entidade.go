package team

import (
	"strings"

	"github.com/hudsontheuz/saas_kanban/internal/domain/shared"
)

type Team struct {
	id       TeamID
	nome     string
	leaderID UserID
}

func NovaTeam(nome string, leaderID UserID) (*Team, error) {
	nome = strings.TrimSpace(nome)
	if nome == "" {
		return nil, ErrNomeObrigatorio
	}
	if leaderID == "" {
		return nil, ErrLeaderObrigatorio
	}

	return &Team{
		id:       TeamID(shared.NovoID()),
		nome:     nome,
		leaderID: leaderID,
	}, nil
}

func (t *Team) ID() TeamID             { return t.id }
func (t *Team) Nome() string           { return t.nome }
func (t *Team) LeaderID() UserID       { return t.leaderID }
func (t *Team) EhLeader(u UserID) bool { return t.leaderID == u }
