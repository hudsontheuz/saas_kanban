package jwt

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	jwtlib "github.com/golang-jwt/jwt/v5"
)

var (
	ErrTokenInvalido  = errors.New("token inválido")
	ErrTokenExpirado  = errors.New("token expirado")
	ErrSubInvalido    = errors.New("sub inválido no token")
	ErrSecretObrigato = errors.New("jwt secret obrigatório")
)

type Validador struct {
	segredo []byte
	emissor string
}

func NovoValidador(segredo, emissor string) (*Validador, error) {
	if segredo == "" {
		return nil, ErrSecretObrigato
	}
	if emissor == "" {
		emissor = "saas_kanban"
	}

	return &Validador{
		segredo: []byte(segredo),
		emissor: emissor,
	}, nil
}

// ValidarEObterIDUsuario valida o JWT (HS256) e retorna o idUsuario a partir do claim "sub".
// Convenção do projeto: "sub" é string numérica do ID do usuário.
func (v *Validador) ValidarEObterIDUsuario(tokenStr string) (int64, error) {
	if tokenStr == "" {
		return 0, ErrTokenInvalido
	}

	interpretador := jwtlib.NewParser(
		jwtlib.WithValidMethods([]string{jwtlib.SigningMethodHS256.Alg()}),
		jwtlib.WithIssuer(v.emissor),
		jwtlib.WithLeeway(30*time.Second),
	)

	claims := &jwtlib.RegisteredClaims{}

	token, err := interpretador.ParseWithClaims(tokenStr, claims, func(t *jwtlib.Token) (any, error) {
		// HS256 apenas
		if t.Method != jwtlib.SigningMethodHS256 {
			return nil, fmt.Errorf("%w: método de assinatura não suportado", ErrTokenInvalido)
		}
		return v.segredo, nil
	})
	if err != nil {
		if errors.Is(err, jwtlib.ErrTokenExpired) {
			return 0, ErrTokenExpirado
		}
		return 0, ErrTokenInvalido
	}
	if !token.Valid {
		return 0, ErrTokenInvalido
	}

	sub := claims.Subject
	if sub == "" {
		return 0, ErrSubInvalido
	}

	idUsuario, convErr := strconv.ParseInt(sub, 10, 64)
	if convErr != nil || idUsuario <= 0 {
		return 0, ErrSubInvalido
	}

	return idUsuario, nil
}
