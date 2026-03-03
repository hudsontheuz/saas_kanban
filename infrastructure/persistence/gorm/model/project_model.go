package model

import "time"

type Projeto struct {
	ID                  int64      `gorm:"column:id;primaryKey"`
	EquipeID            int64      `gorm:"column:equipe_id;not null"`
	Nome                string     `gorm:"column:nome;not null"`
	Status              string     `gorm:"column:status;not null"`
	PermitirSoltarDoing bool       `gorm:"column:permitir_soltar_doing;not null"`
	FechadoEm           *time.Time `gorm:"column:fechado_em"`
	DeletedAt           *time.Time `gorm:"column:deleted_at"`
	CreatedAt           time.Time  `gorm:"column:created_at"`
	UpdatedAt           time.Time  `gorm:"column:updated_at"`
}

func (Projeto) TableName() string { return "projeto" }
