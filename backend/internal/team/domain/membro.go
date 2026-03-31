package team

import user "github.com/hudsontheuz/saas_kanban/internal/user/domain"

type Membro struct {
	UserID user.UserID
	Nome   string
	Email  string
	Role   string
}
