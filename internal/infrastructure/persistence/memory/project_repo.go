package memory

import (
	"sync"

	"github.com/hudsontheuz/saas_kanban/internal/domain/project"
)

type ProjectRepoEmMemoria struct {
	mu    sync.RWMutex
	dados map[project.ProjectID]*project.Project
}

func NovoProjectRepoEmMemoria() *ProjectRepoEmMemoria {
	return &ProjectRepoEmMemoria{dados: map[project.ProjectID]*project.Project{}}
}

func (r *ProjectRepoEmMemoria) Salvar(p *project.Project) error {
	r.mu.Lock()
	defer r.mu.Unlock()
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
