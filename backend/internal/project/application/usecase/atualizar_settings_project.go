package usecase

import (
	projectdto "github.com/hudsontheuz/saas_kanban/internal/project/application/dto"
	projectports "github.com/hudsontheuz/saas_kanban/internal/project/application/ports"
	project "github.com/hudsontheuz/saas_kanban/internal/project/domain"
	teamports "github.com/hudsontheuz/saas_kanban/internal/team/application/ports"
	user "github.com/hudsontheuz/saas_kanban/internal/user/domain"
)

type AtualizarSettingsProjectUseCase struct {
	teams    teamports.TeamRepository
	projects projectports.ProjectRepository
}

func NovoAtualizarSettingsProjectUseCase(
	teams teamports.TeamRepository,
	projects projectports.ProjectRepository,
) *AtualizarSettingsProjectUseCase {
	return &AtualizarSettingsProjectUseCase{
		teams:    teams,
		projects: projects,
	}
}

func (uc *AtualizarSettingsProjectUseCase) Executar(req projectdto.AtualizarSettingsProjectRequest) error {
	p, err := uc.projects.BuscarPorID(project.ProjectID(req.ProjectID))
	if err != nil {
		return err
	}

	tm, err := uc.teams.BuscarPorID(p.TeamID())
	if err != nil {
		return err
	}

	if !tm.EhLeader(user.UserID(req.LeaderID)) {
		return ErrSomenteLeaderPodeGerenciarProject
	}

	if err := p.AtualizarSettings(project.ConfiguracoesProject{
		PermitirSoltarDoingParaTodo: req.PermitirSoltarDoingParaTodo,
	}); err != nil {
		return err
	}

	return uc.projects.Salvar(p)
}
