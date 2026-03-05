package authctx

import "context"

type chaveContexto string

const chaveIDUsuario chaveContexto = "auth.idUsuario"

func ComIDUsuario(ctx context.Context, idUsuario int64) context.Context {
	return context.WithValue(ctx, chaveIDUsuario, idUsuario)
}

func IDUsuarioDoContexto(ctx context.Context) (int64, bool) {
	valor := ctx.Value(chaveIDUsuario)
	if valor == nil {
		return 0, false
	}

	id, ok := valor.(int64)
	return id, ok
}
