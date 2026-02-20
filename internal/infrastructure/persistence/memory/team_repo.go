package memory

import (
	"sync"

	"github.com/hudsontheuz/saas_kanban/internal/domain/team"
)

type TeamRepoEmMemoria struct {
	mu    sync.RWMutex
	dados map[team.TeamID]*team.Team
}

func NovoTeamRepoEmMemoria() *TeamRepoEmMemoria {
	return &TeamRepoEmMemoria{dados: map[team.TeamID]*team.Team{}}
}

func (r *TeamRepoEmMemoria) Salvar(t *team.Team) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.dados[t.ID()] = t
	return nil
}

func (r *TeamRepoEmMemoria) BuscarPorID(id team.TeamID) (*team.Team, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	t, ok := r.dados[id]
	if !ok {
		return nil, ErrNaoEncontrado
	}
	return t, nil
}
