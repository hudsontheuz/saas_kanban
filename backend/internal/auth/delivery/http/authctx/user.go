package authctx

import (
	"context"

	"github.com/hudsontheuz/saas_kanban/internal/user/domain"
)

type chaveContexto string

const chaveIDUsuario chaveContexto = "auth.idUsuario"

func ComIDUsuario(ctx context.Context, idUsuario user.UserID) context.Context {
	return context.WithValue(ctx, chaveIDUsuario, idUsuario)
}

func IDUsuarioDoContexto(ctx context.Context) (user.UserID, bool) {
	valor := ctx.Value(chaveIDUsuario)
	if valor == nil {
		return "", false
	}

	id, ok := valor.(user.UserID)
	return id, ok
}
