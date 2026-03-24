package repo

import (
	"errors"
	"strconv"
	"strings"
	"time"

	shared "github.com/hudsontheuz/saas_kanban/internal/shared/errors"
	userports "github.com/hudsontheuz/saas_kanban/internal/user/application/ports"
	"github.com/hudsontheuz/saas_kanban/internal/user/domain"
	"github.com/hudsontheuz/saas_kanban/internal/user/infrastructure/persistence/gorm/model"
	"gorm.io/gorm"
)

type UserRepo struct{ db *gorm.DB }

var _ userports.UserRepository = (*UserRepo)(nil)

func NewUserRepo(db *gorm.DB) *UserRepo { return &UserRepo{db: db} }

func parseID(s string) (int64, error) {
	s = strings.TrimSpace(s)
	return strconv.ParseInt(s, 10, 64)
}

func (r *UserRepo) Salvar(u *user.Usuario) error {
	if strings.TrimSpace(string(u.ID())) == "" {
		m := model.Usuario{
			Nome:      u.Nome(),
			Email:     u.Email(),
			SenhaHash: u.SenhaHash(),
		}
		if err := r.db.Create(&m).Error; err != nil {
			return err
		}
		return u.DefinirID(user.UserID(strconv.FormatInt(m.ID, 10)))
	}

	id, err := parseID(string(u.ID()))
	if err != nil {
		return err
	}

	updates := map[string]any{
		"nome":       u.Nome(),
		"email":      u.Email(),
		"senha_hash": u.SenhaHash(),
		"updated_at": time.Now(),
	}

	return r.db.Model(&model.Usuario{}).Where("id = ?", id).Updates(updates).Error
}

func (r *UserRepo) BuscarPorEmail(email string) (*user.Usuario, error) {
	email = strings.TrimSpace(strings.ToLower(email))

	var m model.Usuario
	if err := r.db.Where("email = ?", email).First(&m).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, shared.ErrNaoEncontrado
		}
		return nil, err
	}

	return user.HidratarUsuario(
		user.UserID(strconv.FormatInt(m.ID, 10)),
		m.Nome,
		m.Email,
		m.SenhaHash,
	), nil
}

func (r *UserRepo) BuscarPorID(id user.UserID) (*user.Usuario, error) {
	uid, err := parseID(string(id))
	if err != nil {
		return nil, err
	}

	var m model.Usuario
	if err := r.db.First(&m, uid).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, shared.ErrNaoEncontrado
		}
		return nil, err
	}

	return user.HidratarUsuario(
		user.UserID(strconv.FormatInt(m.ID, 10)),
		m.Nome,
		m.Email,
		m.SenhaHash,
	), nil
}
