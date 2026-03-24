package jwt

import (
	"time"

	jwtlib "github.com/golang-jwt/jwt/v5"
	"github.com/hudsontheuz/saas_kanban/internal/user/domain"
)

type Issuer struct {
	segredo []byte
	emissor string
	ttl     time.Duration
}

func NovoIssuer(segredo, emissor string, ttl time.Duration) (*Issuer, error) {
	if segredo == "" {
		return nil, ErrSecretObrigato
	}
	if emissor == "" {
		emissor = "saas_kanban"
	}
	if ttl <= 0 {
		ttl = 24 * time.Hour
	}

	return &Issuer{segredo: []byte(segredo), emissor: emissor, ttl: ttl}, nil
}

func (i *Issuer) Emitir(userID user.UserID) (string, error) {
	agora := time.Now()
	claims := jwtlib.RegisteredClaims{
		Issuer:    i.emissor,
		Subject:   string(userID),
		IssuedAt:  jwtlib.NewNumericDate(agora),
		ExpiresAt: jwtlib.NewNumericDate(agora.Add(i.ttl)),
	}

	token := jwtlib.NewWithClaims(jwtlib.SigningMethodHS256, claims)
	return token.SignedString(i.segredo)
}
