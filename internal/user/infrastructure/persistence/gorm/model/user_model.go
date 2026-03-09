package model

import "time"

type Usuario struct {
	ID        int64     `gorm:"column:id;primaryKey;autoIncrement"`
	Nome      string    `gorm:"column:nome"`
	Email     string    `gorm:"column:email"`
	SenhaHash string    `gorm:"column:senha_hash"`
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}

func (Usuario) TableName() string { return "usuario" }
