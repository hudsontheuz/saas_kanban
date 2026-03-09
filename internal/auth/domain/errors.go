package auth

import "errors"

var (
	ErrNomeObrigatorio      = errors.New("nome obrigatório")
	ErrEmailObrigatorio     = errors.New("email obrigatório")
	ErrSenhaObrigatoria     = errors.New("senha obrigatória")
	ErrCredenciaisInvalidas = errors.New("credenciais inválidas")
	ErrEmailJaCadastrado    = errors.New("email já cadastrado")
)
