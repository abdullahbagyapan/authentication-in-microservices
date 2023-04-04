package ports

type AppPort interface {
	GenerateToken(id string) (string, error)
	ValidateToken(token string) (string, bool, error)
}
