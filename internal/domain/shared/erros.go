package shared

import "errors"

var (
	ErrIDInvalido    = errors.New("id inválido")
	ErrNaoEncontrado = errors.New("não encontrado")
)
