package project

import (
	"time"

	"github.com/hudsontheuz/saas_kanban/internal/domain/shared"
	"github.com/hudsontheuz/saas_kanban/internal/domain/team"
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
	status    StatusProject
	settings  ConfiguracoesProject
	fechadoEm *time.Time
}

func NovoProject(teamID team.TeamID, settings ConfiguracoesProject) (*Project, error) {
	if teamID == "" {
		return nil, ErrTeamObrigatoria
	}

	return &Project{
		id:       ProjectID(shared.NovoID()),
		teamID:   teamID,
		status:   ProjectAtivo,
		settings: settings,
	}, nil
}

func (p *Project) ID() ProjectID                     { return p.id }
func (p *Project) TeamID() team.TeamID               { return p.teamID }
func (p *Project) Settings() ConfiguracoesProject    { return p.settings }
func (p *Project) EstaFechado() bool                 { return p.status == ProjectFechado }

func (p *Project) Fechar(agora time.Time) {
	if p.EstaFechado() {
		return
	}
	p.status = ProjectFechado
	p.fechadoEm = &agora
}
