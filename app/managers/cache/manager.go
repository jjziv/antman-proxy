package managers

type CacheManager struct {
}

func NewManager() *CacheManager {
	return &CacheManager{}
}

func (m *CacheManager) Get(key string, format string) string {
	return ""
}

func (m *CacheManager) Set(key string, img []byte, format string) (string, error) {

	return "", nil
}

func (m *CacheManager) GetPath(key string, format string) string {
	return ""
}
