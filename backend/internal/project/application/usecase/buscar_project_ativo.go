package usecase

import (
	projectdto "github.com/hudsontheuz/saas_kanban/internal/project/application/dto"
	projectports "github.com/hudsontheuz/saas_kanban/internal/project/application/ports"
	team "github.com/hudsontheuz/saas_kanban/internal/team/domain"
)

type BuscarProjectAtivoUseCase struct {
	projects projectports.ProjectRepository
}

func NovoBuscarProjectAtivoUseCase(projects projectports.ProjectRepository) *BuscarProjectAtivoUseCase {
	return &BuscarProjectAtivoUseCase{projects: projects}
}

func (uc *BuscarProjectAtivoUseCase) Executar(req projectdto.BuscarProjectAtivoRequest) (projectdto.BuscarProjectAtivoResponse, error) {
	p, err := uc.projects.BuscarAtivoPorTeamID(team.TeamID(req.TeamID))
	if err != nil {
		return projectdto.BuscarProjectAtivoResponse{}, err
	}

	return projectdto.BuscarProjectAtivoResponse{
		ProjectID:                   string(p.ID()),
		TeamID:                      string(p.TeamID()),
		Nome:                        p.Nome(),
		Status:                      string(p.Status()),
		PermitirSoltarDoingParaTodo: p.Settings().PermitirSoltarDoingParaTodo,
	}, nil
}
