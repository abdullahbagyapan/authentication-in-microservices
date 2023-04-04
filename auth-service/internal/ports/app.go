package ports

type AppPort interface {
	Register(name, username, password, email string) (string, error)
	Login(username, password string) (string, bool, error)
}
