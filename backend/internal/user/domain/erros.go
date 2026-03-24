package user

import "errors"

var (
	ErrNomeObrigatorio      = errors.New("nome obrigatório")
	ErrEmailObrigatorio     = errors.New("email obrigatório")
	ErrSenhaHashObrigatoria = errors.New("senha hash obrigatória")
)
