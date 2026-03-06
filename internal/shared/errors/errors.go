package errors

import sterrors "errors"

var (
	ErrIDInvalido    = sterrors.New("id inválido")
	ErrNaoEncontrado = sterrors.New("não encontrado")
)
