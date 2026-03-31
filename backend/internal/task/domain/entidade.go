package task

import (
	"strings"
	"time"

	project "github.com/hudsontheuz/saas_kanban/internal/project/domain"
	shared "github.com/hudsontheuz/saas_kanban/internal/shared/errors"
	user "github.com/hudsontheuz/saas_kanban/internal/user/domain"
)

type Task struct {
	id                TaskID
	projectID         project.ProjectID
	titulo            string
	descricao         string
	comentarioReview  string
	comentarioEntrega string

	status   StatusTask
	assignee *user.UserID
	isPaused bool
	outcome  *OutcomeTask

	deletedAt *time.Time
	deletedBy *user.UserID
}

func NovaTask(projectID project.ProjectID, titulo, descricao string) (*Task, error) {
	titulo = strings.TrimSpace(titulo)
	descricao = strings.TrimSpace(descricao)

	if projectID == "" {
		return nil, ErrProjetoObrigatorio
	}
	if titulo == "" {
		return nil, ErrTituloObrigatorio
	}
	if descricao == "" {
		return nil, ErrDescricaoObrigatoria
	}

	return &Task{
		id:                "",
		projectID:         projectID,
		titulo:            titulo,
		descricao:         descricao,
		comentarioReview:  "",
		comentarioEntrega: "",
		status:            ToDo,
	}, nil
}

func (t *Task) ID() TaskID                   { return t.id }
func (t *Task) ProjectID() project.ProjectID { return t.projectID }
func (t *Task) Titulo() string               { return t.titulo }
func (t *Task) Descricao() string            { return t.descricao }
func (t *Task) ComentarioReview() string     { return t.comentarioReview }
func (t *Task) ComentarioEntrega() string    { return t.comentarioEntrega }
func (t *Task) Status() StatusTask           { return t.status }
func (t *Task) IsPaused() bool               { return t.isPaused }
func (t *Task) Assignee() *user.UserID       { return t.assignee }
func (t *Task) Outcome() *OutcomeTask        { return t.outcome }
func (t *Task) DeletedAt() *time.Time        { return t.deletedAt }

func (t *Task) ValidarInvariantes() error {
	if t.outcome != nil && t.status != Done {
		return ErrOutcomeSomenteEmDone
	}
	return nil
}

func (t *Task) SelfAssign(userID user.UserID) error {
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

func (t *Task) PodePausarOuRetomar(userID user.UserID) error {
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

func (t *Task) MoverParaInReview(comentarioEntrega string) error {
	if t.status != Doing {
		return ErrMoverParaInReviewSomenteDoing
	}
	if t.isPaused {
		return ErrTaskPausadaNaoPodeIrReview
	}

	comentarioEntrega = strings.TrimSpace(comentarioEntrega)
	if comentarioEntrega == "" {
		return ErrComentarioEntregaObrigatorio
	}

	t.status = InReview
	t.comentarioEntrega = comentarioEntrega
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
	t.comentarioReview = ""

	return nil
}

func (t *Task) ReprovarParaAjustes(motivo string) error {
	if t.status != InReview {
		return ErrReprovarSomenteInReview
	}

	motivo = strings.TrimSpace(motivo)
	if motivo == "" {
		return ErrMotivoReviewObrigatorio
	}

	t.status = ToDo
	t.isPaused = false
	t.comentarioReview = motivo

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

func (t *Task) SoftDelete(quando time.Time, por user.UserID) error {
	if por == "" {
		return shared.ErrIDInvalido
	}

	t.deletedAt = &quando
	u := por
	t.deletedBy = &u

	return nil
}

func (t *Task) PodeMoverParaInReview(userID user.UserID) error {
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

func (t *Task) DefinirID(id TaskID) error {
	if strings.TrimSpace(string(id)) == "" {
		return shared.ErrIDInvalido
	}
	t.id = id
	return nil
}
