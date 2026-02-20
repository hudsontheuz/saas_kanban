package usecase

import "errors"

var (
	ErrSomenteLeaderPodeDecidir = errors.New("somente o leader pode aprovar/reprovar/rejeitar no MVP")
)
