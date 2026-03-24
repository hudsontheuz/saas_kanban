package user

import (
	"strings"

	shared "github.com/hudsontheuz/saas_kanban/internal/shared/errors"
)

type Usuario struct {
	id        UserID
	nome      string
	email     string
	senhaHash string
}

func NovoUsuario(nome, email, senhaHash string) (*Usuario, error) {
	nome = strings.TrimSpace(nome)
	email = strings.TrimSpace(strings.ToLower(email))
	senhaHash = strings.TrimSpace(senhaHash)

	if nome == "" {
		return nil, ErrNomeObrigatorio
	}
	if email == "" {
		return nil, ErrEmailObrigatorio
	}
	if senhaHash == "" {
		return nil, ErrSenhaHashObrigatoria
	}

	return &Usuario{
		nome:      nome,
		email:     email,
		senhaHash: senhaHash,
	}, nil
}

func HidratarUsuario(id UserID, nome, email, senhaHash string) *Usuario {
	return &Usuario{
		id:        id,
		nome:      strings.TrimSpace(nome),
		email:     strings.TrimSpace(strings.ToLower(email)),
		senhaHash: strings.TrimSpace(senhaHash),
	}
}

func (u *Usuario) ID() UserID        { return u.id }
func (u *Usuario) Nome() string      { return u.nome }
func (u *Usuario) Email() string     { return u.email }
func (u *Usuario) SenhaHash() string { return u.senhaHash }

func (u *Usuario) DefinirID(id UserID) error {
	if strings.TrimSpace(string(id)) == "" {
		return shared.ErrIDInvalido
	}
	u.id = id
	return nil
}
