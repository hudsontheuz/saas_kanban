package usecase

import (
	"errors"
	"strings"

	"github.com/hudsontheuz/saas_kanban/internal/auth/application/dto"
	authports "github.com/hudsontheuz/saas_kanban/internal/auth/application/ports"
	auth "github.com/hudsontheuz/saas_kanban/internal/auth/domain"
	shared "github.com/hudsontheuz/saas_kanban/internal/shared/errors"
	userports "github.com/hudsontheuz/saas_kanban/internal/user/application/ports"
)

type LoginUseCase struct {
	users  userports.UserRepository
	hasher authports.PasswordHasher
	issuer authports.TokenIssuer
}

func NovoLoginUseCase(users userports.UserRepository, hasher authports.PasswordHasher, issuer authports.TokenIssuer) *LoginUseCase {
	return &LoginUseCase{users: users, hasher: hasher, issuer: issuer}
}

func (uc *LoginUseCase) Executar(req dto.LoginRequest) (dto.AuthResponse, error) {
	req.Email = strings.TrimSpace(strings.ToLower(req.Email))
	req.Senha = strings.TrimSpace(req.Senha)

	if req.Email == "" {
		return dto.AuthResponse{}, auth.ErrEmailObrigatorio
	}
	if req.Senha == "" {
		return dto.AuthResponse{}, auth.ErrSenhaObrigatoria
	}

	u, err := uc.users.BuscarPorEmail(req.Email)
	if err != nil {
		if errors.Is(err, shared.ErrNaoEncontrado) {
			return dto.AuthResponse{}, auth.ErrCredenciaisInvalidas
		}
		return dto.AuthResponse{}, err
	}

	if err := uc.hasher.Compare(u.SenhaHash(), req.Senha); err != nil {
		return dto.AuthResponse{}, auth.ErrCredenciaisInvalidas
	}

	token, err := uc.issuer.Emitir(u.ID())
	if err != nil {
		return dto.AuthResponse{}, err
	}

	return dto.AuthResponse{Token: token, UserID: string(u.ID())}, nil
}
