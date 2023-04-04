package ports

type CorePort interface {
	Hash(value string) string
}
