package dto

type LoginRequest struct {
	Email string `json:"email"`
	Senha string `json:"senha"`
}
