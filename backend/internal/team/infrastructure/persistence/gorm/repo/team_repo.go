package repo

import (
	"errors"
	"strconv"
	"strings"
	"time"

	shared "github.com/hudsontheuz/saas_kanban/internal/shared/errors"
	teamports "github.com/hudsontheuz/saas_kanban/internal/team/application/ports"
	team "github.com/hudsontheuz/saas_kanban/internal/team/domain"
	"github.com/hudsontheuz/saas_kanban/internal/team/infrastructure/persistence/gorm/model"
	user "github.com/hudsontheuz/saas_kanban/internal/user/domain"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type TeamRepo struct {
	db *gorm.DB
}

var _ teamports.TeamRepository = (*TeamRepo)(nil)

func NewTeamRepo(db *gorm.DB) *TeamRepo {
	return &TeamRepo{db: db}
}

func parseTeamID(s string) (int64, error) {
	s = strings.TrimSpace(s)
	return strconv.ParseInt(s, 10, 64)
}

func (r *TeamRepo) Salvar(t *team.Team) error {
	leaderID, err := parseTeamID(string(t.LeaderID()))
	if err != nil {
		return err
	}

	if strings.TrimSpace(string(t.ID())) == "" {
		return r.db.Transaction(func(tx *gorm.DB) error {
			m := model.Equipe{
				Nome:           t.Nome(),
				LiderUsuarioID: leaderID,
			}
			if err := tx.Create(&m).Error; err != nil {
				return err
			}

			if err := tx.Clauses(clause.OnConflict{DoNothing: true}).Create(&model.EquipeMembro{
				EquipeID:  m.ID,
				UsuarioID: leaderID,
				Role:      "LEADER",
			}).Error; err != nil {
				return err
			}

			return t.DefinirID(team.TeamID(strconv.FormatInt(m.ID, 10)))
		})
	}

	id, err := parseTeamID(string(t.ID()))
	if err != nil {
		return err
	}

	return r.db.Transaction(func(tx *gorm.DB) error {
		updates := map[string]any{
			"nome":             t.Nome(),
			"lider_usuario_id": leaderID,
			"updated_at":       time.Now(),
		}

		if err := tx.Model(&model.Equipe{}).Where("id = ?", id).Updates(updates).Error; err != nil {
			return err
		}

		if err := tx.Clauses(clause.OnConflict{DoNothing: true}).Create(&model.EquipeMembro{
			EquipeID:  id,
			UsuarioID: leaderID,
			Role:      "LEADER",
		}).Error; err != nil {
			return err
		}

		return nil
	})
}

func (r *TeamRepo) BuscarPorID(id team.TeamID) (*team.Team, error) {
	tid, err := parseTeamID(string(id))
	if err != nil {
		return nil, err
	}

	var m model.Equipe
	if err := r.db.First(&m, tid).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, shared.ErrNaoEncontrado
		}
		return nil, err
	}

	return team.HidratarTeam(
		team.TeamID(strconv.FormatInt(m.ID, 10)),
		m.Nome,
		user.UserID(strconv.FormatInt(m.LiderUsuarioID, 10)),
	), nil
}

func (r *TeamRepo) ListarPorUsuarioID(userID user.UserID) ([]*team.Team, error) {
	uid, err := parseTeamID(string(userID))
	if err != nil {
		return nil, err
	}

	var models []model.Equipe
	err = r.db.
		Table("equipe").
		Joins("JOIN equipe_membro ON equipe_membro.equipe_id = equipe.id").
		Where("equipe_membro.usuario_id = ?", uid).
		Order("equipe.id ASC").
		Find(&models).Error
	if err != nil {
		return nil, err
	}

	items := make([]*team.Team, 0, len(models))
	for _, m := range models {
		items = append(items, team.HidratarTeam(
			team.TeamID(strconv.FormatInt(m.ID, 10)),
			m.Nome,
			user.UserID(strconv.FormatInt(m.LiderUsuarioID, 10)),
		))
	}

	return items, nil
}

type membroRow struct {
	UsuarioID int64  `gorm:"column:usuario_id"`
	Nome      string `gorm:"column:nome"`
	Email     string `gorm:"column:email"`
	Role      string `gorm:"column:role"`
}

func (r *TeamRepo) ListarMembros(id team.TeamID) ([]team.Membro, error) {
	tid, err := parseTeamID(string(id))
	if err != nil {
		return nil, err
	}

	var rows []membroRow
	err = r.db.
		Table("equipe_membro").
		Select("equipe_membro.usuario_id, usuario.nome, usuario.email, equipe_membro.role").
		Joins("JOIN usuario ON usuario.id = equipe_membro.usuario_id").
		Where("equipe_membro.equipe_id = ?", tid).
		Order("equipe_membro.id ASC").
		Scan(&rows).Error
	if err != nil {
		return nil, err
	}

	membros := make([]team.Membro, 0, len(rows))
	for _, row := range rows {
		membros = append(membros, team.Membro{
			UserID: user.UserID(strconv.FormatInt(row.UsuarioID, 10)),
			Nome:   row.Nome,
			Email:  row.Email,
			Role:   row.Role,
		})
	}

	return membros, nil
}
