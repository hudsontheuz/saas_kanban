package usecase

import (
	"errors"
	"strings"

	"github.com/hudsontheuz/saas_kanban/internal/auth/application/dto"
	authports "github.com/hudsontheuz/saas_kanban/internal/auth/application/ports"
	auth "github.com/hudsontheuz/saas_kanban/internal/auth/domain"
	shared "github.com/hudsontheuz/saas_kanban/internal/shared/errors"
	userports "github.com/hudsontheuz/saas_kanban/internal/user/application/ports"
	"github.com/hudsontheuz/saas_kanban/internal/user/domain"
)

type RegisterUseCase struct {
	users  userports.UserRepository
	hasher authports.PasswordHasher
	issuer authports.TokenIssuer
}

func NovoRegisterUseCase(users userports.UserRepository, hasher authports.PasswordHasher, issuer authports.TokenIssuer) *RegisterUseCase {
	return &RegisterUseCase{users: users, hasher: hasher, issuer: issuer}
}

func (uc *RegisterUseCase) Executar(req dto.RegisterRequest) (dto.AuthResponse, error) {
	req.Nome = strings.TrimSpace(req.Nome)
	req.Email = strings.TrimSpace(strings.ToLower(req.Email))
	req.Senha = strings.TrimSpace(req.Senha)

	if req.Nome == "" {
		return dto.AuthResponse{}, auth.ErrNomeObrigatorio
	}
	if req.Email == "" {
		return dto.AuthResponse{}, auth.ErrEmailObrigatorio
	}
	if req.Senha == "" {
		return dto.AuthResponse{}, auth.ErrSenhaObrigatoria
	}

	_, err := uc.users.BuscarPorEmail(req.Email)
	if err == nil {
		return dto.AuthResponse{}, auth.ErrEmailJaCadastrado
	}
	if !errors.Is(err, shared.ErrNaoEncontrado) {
		return dto.AuthResponse{}, err
	}

	hash, err := uc.hasher.Hash(req.Senha)
	if err != nil {
		return dto.AuthResponse{}, err
	}

	u, err := user.NovoUsuario(req.Nome, req.Email, hash)
	if err != nil {
		return dto.AuthResponse{}, err
	}
	if err := uc.users.Salvar(u); err != nil {
		return dto.AuthResponse{}, err
	}

	token, err := uc.issuer.Emitir(u.ID())
	if err != nil {
		return dto.AuthResponse{}, err
	}

	return dto.AuthResponse{Token: token, UserID: string(u.ID())}, nil
}
