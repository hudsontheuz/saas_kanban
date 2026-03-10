package memory

import (
	"strconv"
	"sync"

	"github.com/hudsontheuz/saas_kanban/internal/project/domain"
	"github.com/hudsontheuz/saas_kanban/internal/team/domain"
)

type ProjectRepoEmMemoria struct {
	mu     sync.RWMutex
	dados  map[project.ProjectID]*project.Project
	nextID int64
}

func NovoProjectRepoEmMemoria() *ProjectRepoEmMemoria {
	return &ProjectRepoEmMemoria{
		dados:  map[project.ProjectID]*project.Project{},
		nextID: 1,
	}
}

func (r *ProjectRepoEmMemoria) Salvar(p *project.Project) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if p == nil {
		return ErrNaoEncontrado
	}

	// Se vier sem ID, simula DB gerando ID
	if string(p.ID()) == "" {
		id := project.ProjectID(strconv.FormatInt(r.nextID, 10))
		r.nextID++
		_ = p.DefinirID(id)
	}

	r.dados[p.ID()] = p
	return nil
}

func (r *ProjectRepoEmMemoria) BuscarPorID(id project.ProjectID) (*project.Project, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	p, ok := r.dados[id]
	if !ok {
		return nil, ErrNaoEncontrado
	}
	return p, nil
}

func (r *ProjectRepoEmMemoria) BuscarAtivoPorTeamID(teamID team.TeamID) (*project.Project, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, p := range r.dados {
		if p.TeamID() == teamID && !p.EstaFechado() {
			return p, nil
		}
	}
	return nil, ErrNaoEncontrado
}
