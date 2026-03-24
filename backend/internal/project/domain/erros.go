package project

import "errors"

var (
	ErrTeamObrigatoria = errors.New("team é obrigatória")
	ErrProjetoFechado  = errors.New("projeto está fechado")
	ErrNomeObrigatorio = errors.New("nome do projeto é obrigatório")
)
