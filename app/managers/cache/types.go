package managers

type Manager interface {
	Get(key string, format string) string
	Set(key string, img []byte, format string) (string, error)
	GetPath(key string, format string) string
}
