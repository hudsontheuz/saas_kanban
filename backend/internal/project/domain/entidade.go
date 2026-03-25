package project

import (
	"strings"
	"time"

	shared "github.com/hudsontheuz/saas_kanban/internal/shared/errors"
	team "github.com/hudsontheuz/saas_kanban/internal/team/domain"
)

type ProjectID string
type StatusProject string

const (
	ProjectAtivo   StatusProject = "ACTIVE"
	ProjectFechado StatusProject = "CLOSED"
)

type ConfiguracoesProject struct {
	PermitirSoltarDoingParaTodo bool
}

type Project struct {
	id        ProjectID
	teamID    team.TeamID
	nome      string
	status    StatusProject
	settings  ConfiguracoesProject
	fechadoEm *time.Time
}

func NovoProject(teamID team.TeamID, nome string, settings ConfiguracoesProject) (*Project, error) {
	if teamID == "" {
		return nil, ErrTeamObrigatoria
	}

	nome = strings.TrimSpace(nome)
	if nome == "" {
		return nil, ErrNomeObrigatorio // vou te falar abaixo onde criar esse erro
	}

	return &Project{
		id:       "", // Será gerado pelo repositório
		teamID:   teamID,
		nome:     nome,
		status:   ProjectAtivo,
		settings: settings,
	}, nil
}

func (p *Project) ID() ProjectID                  { return p.id }
func (p *Project) TeamID() team.TeamID            { return p.teamID }
func (p *Project) Nome() string                   { return p.nome }
func (p *Project) Settings() ConfiguracoesProject { return p.settings }
func (p *Project) EstaFechado() bool              { return p.status == ProjectFechado }
func (p *Project) FechadoEm() *time.Time          { return p.fechadoEm }
func (p *Project) Status() StatusProject          { return p.status }

func (p *Project) Fechar(agora time.Time) error {
	if p.EstaFechado() {
		return ErrProjetoFechado
	}
	p.status = ProjectFechado
	p.fechadoEm = &agora
	return nil
}

func (p *Project) DefinirID(id ProjectID) error {
	if id == "" {
		return shared.ErrIDInvalido
	}
	p.id = id
	return nil
}

func (p *Project) AtualizarSettings(settings ConfiguracoesProject) error {
	if p.EstaFechado() {
		return ErrProjetoFechado
	}
	p.settings = settings
	return nil
}
