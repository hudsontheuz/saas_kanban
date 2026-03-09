package memory

import (
	"strconv"
	"strings"
	"sync"

	"github.com/hudsontheuz/saas_kanban/internal/user/domain"
)

type UserRepoEmMemoria struct {
	mu       sync.RWMutex
	dados    map[user.UserID]*user.Usuario
	porEmail map[string]user.UserID
	nextID   int64
}

func NovoUserRepoEmMemoria() *UserRepoEmMemoria {
	return &UserRepoEmMemoria{
		dados:    map[user.UserID]*user.Usuario{},
		porEmail: map[string]user.UserID{},
		nextID:   1,
	}
}

func (r *UserRepoEmMemoria) Salvar(u *user.Usuario) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if u == nil {
		return ErrNaoEncontrado
	}

	if string(u.ID()) == "" {
		id := user.UserID(strconv.FormatInt(r.nextID, 10))
		r.nextID++
		_ = u.DefinirID(id)
	}

	email := strings.TrimSpace(strings.ToLower(u.Email()))
	r.dados[u.ID()] = u
	r.porEmail[email] = u.ID()
	return nil
}

func (r *UserRepoEmMemoria) BuscarPorEmail(email string) (*user.Usuario, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	id, ok := r.porEmail[strings.TrimSpace(strings.ToLower(email))]
	if !ok {
		return nil, ErrNaoEncontrado
	}
	u, ok := r.dados[id]
	if !ok {
		return nil, ErrNaoEncontrado
	}
	return u, nil
}

func (r *UserRepoEmMemoria) BuscarPorID(id user.UserID) (*user.Usuario, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	u, ok := r.dados[id]
	if !ok {
		return nil, ErrNaoEncontrado
	}
	return u, nil
}
