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
