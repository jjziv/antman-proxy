package managers

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"
)

const (
	DefaultCacheDir = "image_cache"
	DefaultMaxAge   = int64(86400)
)

type Config struct {
	CacheDir string
	MaxAge   int64
}

type CacheManager struct {
	cacheDir string
	maxAge   int64
}

func (m *CacheManager) ensureCacheDir() error {
	return os.MkdirAll(m.cacheDir, 0755)
}

func NewManager(cfg *Config) (*CacheManager, error) {
	if cfg == nil {
		return nil, fmt.Errorf("CacheManager config is nil!")
	}

	if cfg.CacheDir == "" {
		cfg.CacheDir = DefaultCacheDir
	}

	if cfg.MaxAge == 0 {
		cfg.MaxAge = DefaultMaxAge
	}

	m := &CacheManager{
		cacheDir: cfg.CacheDir,
		maxAge:   cfg.MaxAge,
	}

	err := m.ensureCacheDir()
	if err != nil {
		return nil, err
	}

	return m, nil
}

func (m *CacheManager) Get(key string, format string) string {
	path := m.GetPath(key, format)

	info, err := os.Stat(path)
	if err == nil {
		if time.Now().Unix()-info.ModTime().Unix() <= m.maxAge {
			return path
		}

		err = os.Remove(path)
		if err != nil {
			log.Println(err)
			return ""
		}
	}

	return ""
}

func (m *CacheManager) Set(key string, img []byte, format string) (string, error) {
	path := m.GetPath(key, format)
	return path, os.WriteFile(path, img, 0644)
}

func (m *CacheManager) GetPath(key string, format string) string {
	ext := "jpg"
	if format != "jpeg" {
		ext = format
	}
	return filepath.Join(m.cacheDir, fmt.Sprintf("%s.%s", key, ext))
}
