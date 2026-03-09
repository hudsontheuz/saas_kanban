package jwt

import (
	"errors"
	"fmt"
	"strings"
	"time"

	jwtlib "github.com/golang-jwt/jwt/v5"
	"github.com/hudsontheuz/saas_kanban/internal/user/domain"
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

func (v *Validador) ValidarEObterIDUsuario(tokenStr string) (user.UserID, error) {
	if tokenStr == "" {
		return "", ErrTokenInvalido
	}

	interpretador := jwtlib.NewParser(
		jwtlib.WithValidMethods([]string{jwtlib.SigningMethodHS256.Alg()}),
		jwtlib.WithIssuer(v.emissor),
		jwtlib.WithLeeway(30*time.Second),
	)

	claims := &jwtlib.RegisteredClaims{}

	token, err := interpretador.ParseWithClaims(tokenStr, claims, func(t *jwtlib.Token) (any, error) {
		if t.Method != jwtlib.SigningMethodHS256 {
			return nil, fmt.Errorf("%w: método de assinatura não suportado", ErrTokenInvalido)
		}
		return v.segredo, nil
	})
	if err != nil {
		if errors.Is(err, jwtlib.ErrTokenExpired) {
			return "", ErrTokenExpirado
		}
		return "", ErrTokenInvalido
	}
	if !token.Valid {
		return "", ErrTokenInvalido
	}

	sub := strings.TrimSpace(claims.Subject)
	if sub == "" {
		return "", ErrSubInvalido
	}

	return user.UserID(sub), nil
}
