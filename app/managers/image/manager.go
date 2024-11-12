package managers

import (
	"errors"

	cacheManager "antman-proxy/managers/cache"
)

type Config struct {
	CacheManager cacheManager.Manager
}

type ImageManager struct {
	cacheManager cacheManager.Manager
}

func NewManager(cfg *Config) (*ImageManager, error) {
	if cfg == nil {
		return nil, errors.New("ImageManager Config is nil!")
	}

	if cfg.CacheManager == nil {
		return nil, errors.New("cfg.CacheManager is nil!")
	}

	return &ImageManager{
		cacheManager: cfg.CacheManager,
	}, nil
}

func (m *ImageManager) IsURLAllowed(imageURL string) bool {
	return true
}

func (m *ImageManager) ProcessImage(imageURL string, width, height int, format string) (string, error) {
	return "", nil
}
