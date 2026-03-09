package ports

type PasswordHasher interface {
	Hash(senha string) (string, error)
	Compare(hash, senha string) error
}
