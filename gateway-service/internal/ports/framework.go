package ports

type RedisPort interface {
	SetToken(token string) error
	GetToken(token string) error
}
