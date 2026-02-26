package usecase

import "errors"

var (
	ErrSomenteLeaderPodeGerenciarProject = errors.New("somente o leader pode criar/fechar project no MVP")
	ErrJaExisteProjectAtivo              = errors.New("já existe project ativo para essa team")
)
