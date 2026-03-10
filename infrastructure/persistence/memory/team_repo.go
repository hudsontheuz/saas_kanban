package memory

import (
	"strconv"
	"sync"

	"github.com/hudsontheuz/saas_kanban/internal/domain/team"
)

type TeamRepoEmMemoria struct {
	mu     sync.RWMutex
	dados  map[team.TeamID]*team.Team
	nextID int64
}

func NovoTeamRepoEmMemoria() *TeamRepoEmMemoria {
	return &TeamRepoEmMemoria{
		dados:  map[team.TeamID]*team.Team{},
		nextID: 1,
	}
}

func (r *TeamRepoEmMemoria) Salvar(t *team.Team) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if t == nil {
		return ErrNaoEncontrado
	}

	// Se vier sem ID, simula DB gerando ID
	if string(t.ID()) == "" {
		id := team.TeamID(strconv.FormatInt(r.nextID, 10))
		r.nextID++
		_ = t.DefinirID(id)
	}

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
