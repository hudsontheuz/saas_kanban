package task

import (
	"strings"
	"time"

	"github.com/hudsontheuz/saas_kanban/internal/domain/project"
	"github.com/hudsontheuz/saas_kanban/internal/domain/shared"
	"github.com/hudsontheuz/saas_kanban/internal/domain/team"
)

type Task struct {
	id        TaskID
	projectID project.ProjectID
	titulo    string

	status   StatusTask
	assignee *team.UserID
	isPaused bool
	outcome  *OutcomeTask

	deletedAt *time.Time
	deletedBy *team.UserID
}

func NovaTask(projectID project.ProjectID, titulo string) (*Task, error) {
	titulo = strings.TrimSpace(titulo)

	if projectID == "" {
		return nil, ErrProjetoObrigatorio
	}
	if titulo == "" {
		return nil, ErrTituloObrigatorio
	}

	return &Task{
		id:        TaskID(shared.NovoID()),
		projectID: projectID,
		titulo:    titulo,
		status:    ToDo,
	}, nil
}

func (t *Task) ID() TaskID                   { return t.id }
func (t *Task) ProjectID() project.ProjectID { return t.projectID }
func (t *Task) Titulo() string               { return t.titulo }
func (t *Task) Status() StatusTask           { return t.status }
func (t *Task) IsPaused() bool               { return t.isPaused }
func (t *Task) Assignee() *team.UserID       { return t.assignee }
func (t *Task) Outcome() *OutcomeTask        { return t.outcome }
func (t *Task) DeletedAt() *time.Time        { return t.deletedAt }

func (t *Task) ValidarInvariantes() error {
	if t.outcome != nil && t.status != Done {
		return ErrOutcomeSomenteEmDone
	}
	return nil
}

func (t *Task) SelfAssign(userID team.UserID) error {
	if t.status != ToDo {
		return ErrTransicaoInvalida
	}
	if userID == "" {
		return shared.ErrIDInvalido
	}
	u := userID
	t.assignee = &u
	t.status = Doing
	t.isPaused = false
	return nil
}

func (t *Task) Pausar() error {
	if t.status != Doing {
		return ErrSomenteDoingPodePausar
	}
	if t.isPaused {
		return ErrJaPausada
	}
	t.isPaused = true
	return nil
}

func (t *Task) Retomar() error {
	if t.status != Doing {
		return ErrSomenteDoingPodePausar
	}
	if !t.isPaused {
		return ErrNaoEstaPausada
	}
	t.isPaused = false
	return nil
}

func (t *Task) PodePausarOuRetomar(userID team.UserID) error {
	if t.assignee == nil {
		return ErrSemAssignee
	}
	if userID == "" {
		return shared.ErrIDInvalido
	}
	if *t.assignee != userID {
		return ErrSomenteAssigneePodePausar
	}
	return nil
}

func (t *Task) MoverParaInReview() error {
	if t.status != Doing {
		return ErrTransicaoInvalida
	}
	if t.assignee == nil {
		return ErrSemAssignee
	}
	t.status = InReview
	return nil
}

func (t *Task) Aprovar() error {
	if t.status != InReview {
		return ErrTransicaoInvalida
	}
	out := OutcomeApproved
	t.status = Done
	t.outcome = &out
	t.isPaused = false
	return nil
}

func (t *Task) ReprovarParaAjustes() error {
	if t.status != InReview {
		return ErrReprovarSomenteInReview
	}
	t.status = Doing
	t.isPaused = false
	return nil
}

func (t *Task) RejeitarEmToDo() error {
	if t.status != ToDo {
		return ErrRejeitarSomenteEmTodo
	}
	out := OutcomeRejected
	t.status = Done
	t.outcome = &out
	t.isPaused = false
	return nil
}

func (t *Task) SoftDelete(quando time.Time, por team.UserID) error {
	if por == "" {
		return shared.ErrIDInvalido
	}
	t.deletedAt = &quando
	u := por
	t.deletedBy = &u
	return nil
}

func (t *Task) PodeMoverParaInReview(userID team.UserID) error {
	if t.assignee == nil {
		return ErrSemAssignee
	}
	if userID == "" {
		return shared.ErrIDInvalido
	}
	if *t.assignee != userID {
		return ErrSomenteAssigneePodeMoverInReview
	}
	return nil
}
