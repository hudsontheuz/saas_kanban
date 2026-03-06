package repo

import (
	"strconv"
	"strings"
	"time"

	"github.com/hudsontheuz/saas_kanban/internal/project/domain"
	shared "github.com/hudsontheuz/saas_kanban/internal/shared/errors"
	taskports "github.com/hudsontheuz/saas_kanban/internal/task/application/ports"
	"github.com/hudsontheuz/saas_kanban/internal/task/domain"
	"github.com/hudsontheuz/saas_kanban/internal/task/infrastructure/persistence/gorm/model"
	"github.com/hudsontheuz/saas_kanban/internal/team/domain"
	"gorm.io/gorm"
)

type TaskRepo struct{ db *gorm.DB }

var _ taskports.TaskRepository = (*TaskRepo)(nil)

func NewTaskRepo(db *gorm.DB) *TaskRepo { return &TaskRepo{db: db} }

func parseID(s string) (int64, error) {
	s = strings.TrimSpace(s)
	return strconv.ParseInt(s, 10, 64)
}

func (r *TaskRepo) Salvar(tk *task.Task) error {
	projID, err := parseID(string(tk.ProjectID()))
	if err != nil {
		return err
	}

	var assigneeID *int64
	if a := tk.Assignee(); a != nil {
		uid, err := parseID(string(*a))
		if err != nil {
			return err
		}
		assigneeID = &uid
	}

	var outcome *string
	if o := tk.Outcome(); o != nil {
		s := string(*o)
		outcome = &s
	}

	var deletedAt *time.Time = tk.DeletedAt()

	if strings.TrimSpace(string(tk.ID())) == "" {
		m := model.Tarefa{
			ProjetoID:          projID,
			Titulo:             tk.Titulo(),
			Status:             string(tk.Status()),
			UsuarioAtribuidoID: assigneeID,
			Pausada:            tk.IsPaused(),
			Outcome:            outcome,
			DeletedAt:          deletedAt,
		}
		if err := r.db.Create(&m).Error; err != nil {
			return err
		}
		return tk.DefinirID(task.TaskID(strconv.FormatInt(m.ID, 10)))
	}

	id, err := parseID(string(tk.ID()))
	if err != nil {
		return err
	}

	updates := map[string]any{
		"projeto_id":           projID,
		"titulo":               tk.Titulo(),
		"status":               string(tk.Status()),
		"usuario_atribuido_id": assigneeID,
		"pausada":              tk.IsPaused(),
		"outcome":              outcome,
		"deleted_at":           deletedAt,
		"updated_at":           time.Now(),
	}

	return r.db.Model(&model.Tarefa{}).Where("id = ?", id).Updates(updates).Error
}

func (r *TaskRepo) BuscarPorID(id task.TaskID) (*task.Task, error) {
	i, err := parseID(string(id))
	if err != nil {
		return nil, err
	}

	var m model.Tarefa
	if err := r.db.First(&m, i).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, shared.ErrNaoEncontrado
		}
		return nil, err
	}

	var domAssignee *team.UserID
	if m.UsuarioAtribuidoID != nil {
		u := team.UserID(strconv.FormatInt(*m.UsuarioAtribuidoID, 10))
		domAssignee = &u
	}

	var domOutcome *task.OutcomeTask
	if m.Outcome != nil {
		o := task.OutcomeTask(*m.Outcome)
		domOutcome = &o
	}

	return task.HidratarTask(
		task.TaskID(strconv.FormatInt(m.ID, 10)),
		project.ProjectID(strconv.FormatInt(m.ProjetoID, 10)),
		m.Titulo,
		task.StatusTask(m.Status),
		domAssignee,
		m.Pausada,
		domOutcome,
		m.DeletedAt,
		nil,
	), nil
}

func (r *TaskRepo) ExisteDoingAtivaParaUser(userID team.UserID) (bool, error) {
	uid, err := parseID(string(userID))
	if err != nil {
		return false, err
	}

	var count int64
	err = r.db.Model(&model.Tarefa{}).
		Where("usuario_atribuido_id = ? AND status = 'DOING' AND pausada = false AND deleted_at IS NULL", uid).
		Count(&count).Error
	return count > 0, err
}
