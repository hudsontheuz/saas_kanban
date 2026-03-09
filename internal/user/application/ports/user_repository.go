package ports

import "github.com/hudsontheuz/saas_kanban/internal/user/domain"

type UserRepository interface {
	Salvar(u *user.Usuario) error
	BuscarPorEmail(email string) (*user.Usuario, error)
	BuscarPorID(id user.UserID) (*user.Usuario, error)
}
