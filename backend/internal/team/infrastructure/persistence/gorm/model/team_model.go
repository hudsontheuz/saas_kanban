package model

import "time"

type Equipe struct {
	ID             int64     `gorm:"column:id;primaryKey"`
	Nome           string    `gorm:"column:nome;not null"`
	LiderUsuarioID int64     `gorm:"column:lider_usuario_id;not null"`
	CreatedAt      time.Time `gorm:"column:created_at"`
	UpdatedAt      time.Time `gorm:"column:updated_at"`
}

func (Equipe) TableName() string { return "equipe" }

type EquipeMembro struct {
	ID        int64     `gorm:"column:id;primaryKey"`
	EquipeID  int64     `gorm:"column:equipe_id;not null"`
	UsuarioID int64     `gorm:"column:usuario_id;not null"`
	Role      string    `gorm:"column:role;not null"`
	CreatedAt time.Time `gorm:"column:created_at"`
}

func (EquipeMembro) TableName() string { return "equipe_membro" }
