package model

import "time"

type Tarefa struct {
	ID                 int64      `gorm:"column:id;primaryKey"`
	ProjetoID          int64      `gorm:"column:projeto_id;not null"`
	Titulo             string     `gorm:"column:titulo;not null"`
	Descricao          *string    `gorm:"column:descricao"`
	ComentarioEntrega  *string    `gorm:"column:comentario_entrega"`
	ComentarioReview   *string    `gorm:"column:comentario_review"`
	Status             string     `gorm:"column:status;not null"`
	UsuarioAtribuidoID *int64     `gorm:"column:usuario_atribuido_id"`
	Pausada            bool       `gorm:"column:pausada;not null"`
	Outcome            *string    `gorm:"column:outcome"`
	DeletedAt          *time.Time `gorm:"column:deleted_at"`
	DeletedBy          *int64     `gorm:"column:deleted_by"`
	CreatedAt          time.Time  `gorm:"column:created_at"`
	UpdatedAt          time.Time  `gorm:"column:updated_at"`
}

func (Tarefa) TableName() string { return "tarefa" }
