package team

import "errors"

var (
	ErrNomeObrigatorio   = errors.New("nome da team é obrigatório")
	ErrLeaderObrigatorio = errors.New("leader é obrigatório")
)