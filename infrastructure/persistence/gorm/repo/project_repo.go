package repo

import (
	"strconv"
	"strings"
	"time"

	"github.com/hudsontheuz/saas_kanban/infrastructure/persistence/gorm/model"
	projectports "github.com/hudsontheuz/saas_kanban/internal/application/project/ports"
	"github.com/hudsontheuz/saas_kanban/internal/domain/project"
	shared "github.com/hudsontheuz/saas_kanban/internal/domain/shared"
	"github.com/hudsontheuz/saas_kanban/internal/domain/team"
	"gorm.io/gorm"
)

type ProjectRepo struct{ db *gorm.DB }

var _ projectports.ProjectRepository = (*ProjectRepo)(nil)

func NewProjectRepo(db *gorm.DB) *ProjectRepo { return &ProjectRepo{db: db} }

func parseID2(s string) (int64, error) {
	s = strings.TrimSpace(s)
	return strconv.ParseInt(s, 10, 64)
}

func (r *ProjectRepo) Salvar(p *project.Project) error {
	tid, err := parseID2(string(p.TeamID()))
	if err != nil {
		return err
	}

	status := "ACTIVE"
	var fechadoEm *time.Time = nil
	if p.EstaFechado() {
		status = "CLOSED"
		// você não tem getter de fechadoEm; por enquanto fica nil
	}

	permitir := p.Settings().PermitirSoltarDoingParaTodo

	if strings.TrimSpace(string(p.ID())) == "" {
		m := model.Projeto{
			EquipeID:            tid,
			Nome:                p.Nome(),
			Status:              status,
			PermitirSoltarDoing: permitir,
			FechadoEm:           fechadoEm,
		}
		if err := r.db.Create(&m).Error; err != nil {
			return err
		}
		return p.DefinirID(project.ProjectID(strconv.FormatInt(m.ID, 10)))
	}

	id, err := parseID2(string(p.ID()))
	if err != nil {
		return err
	}

	updates := map[string]any{
		"equipe_id":             tid,
		"nome":                  p.Nome(),
		"status":                status,
		"permitir_soltar_doing": permitir,
		"fechado_em":            fechadoEm,
		"updated_at":            time.Now(),
	}

	return r.db.Model(&model.Projeto{}).Where("id = ?", id).Updates(updates).Error
}

func (r *ProjectRepo) BuscarPorID(id project.ProjectID) (*project.Project, error) {
	pid, err := parseID2(string(id))
	if err != nil {
		return nil, err
	}

	var m model.Projeto
	if err := r.db.First(&m, pid).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, shared.ErrNaoEncontrado
		}
		return nil, err
	}

	// Você ainda não tem HidratarProject no domínio.
	// -> Crie domain/project/hidratar.go (vou te mandar se quiser já).
	return project.HidratarProject(
		project.ProjectID(strconv.FormatInt(m.ID, 10)),
		team.TeamID(strconv.FormatInt(m.EquipeID, 10)),
		m.Nome,
		project.StatusProject(m.Status),
		project.ConfiguracoesProject{PermitirSoltarDoingParaTodo: m.PermitirSoltarDoing},
		m.FechadoEm,
	), nil
}

func (r *ProjectRepo) BuscarAtivoPorTeamID(teamID team.TeamID) (*project.Project, error) {
	tid, err := parseID2(string(teamID))
	if err != nil {
		return nil, err
	}

	var m model.Projeto
	err = r.db.Where("equipe_id = ? AND status = 'ACTIVE' AND deleted_at IS NULL", tid).First(&m).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, shared.ErrNaoEncontrado
		}
		return nil, err
	}

	return project.HidratarProject(
		project.ProjectID(strconv.FormatInt(m.ID, 10)),
		team.TeamID(strconv.FormatInt(m.EquipeID, 10)),
		m.Nome,
		project.StatusProject(m.Status),
		project.ConfiguracoesProject{PermitirSoltarDoingParaTodo: m.PermitirSoltarDoing},
		m.FechadoEm,
	), nil
}
