package ports

import "github.com/hudsontheuz/saas_kanban/internal/user/domain"

type TokenIssuer interface {
	Emitir(userID user.UserID) (string, error)
}
