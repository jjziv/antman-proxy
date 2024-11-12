package managers

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestNewManager(t *testing.T) {
	t.Parallel()

	t.Run("Config is nil", func(t *testing.T) {
		_, err := NewManager(nil)
		assert.Equal(t, "CacheManager config is nil!", err.Error())
	})

	t.Run("cfg.CacheDir is blank", func(t *testing.T) {
		m, err := NewManager(&Config{
			CacheDir: "",
			MaxAge:   DefaultMaxAge,
		})

		if err != nil {
			t.FailNow()
		}

		assert.Equal(t, DefaultCacheDir, m.cacheDir)
	})

	t.Run("cfg.MaxAge is 0", func(t *testing.T) {
		m, err := NewManager(&Config{
			CacheDir: DefaultCacheDir,
			MaxAge:   0,
		})

		if err != nil {
			t.FailNow()
		}

		assert.Equal(t, DefaultMaxAge, m.maxAge)
	})
}

func TestCacheManager_Get(t *testing.T) {
	t.Parallel()

	tempDir := t.TempDir()

	// Create test file
	testPath := filepath.Join(tempDir, "my_key.jpg")
	err := os.WriteFile(testPath, []byte("test"), 0644)
	assert.NoError(t, err)

	tests := []struct {
		name         string
		key          string
		format       string
		wait         time.Duration
		expectedPath string
	}{
		{
			name:         "successful processing with cache hit",
			key:          "my_key",
			format:       "jpeg",
			wait:         0,
			expectedPath: testPath,
		},
		{
			name:         "expired cache",
			key:          "my_key",
			format:       "jpeg",
			wait:         time.Second * 2,
			expectedPath: "",
		},
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	manager, err := NewManager(&Config{CacheDir: tempDir, MaxAge: 1})
	if err != nil {
		t.FailNow()
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.wait > 0 {
				time.Sleep(tt.wait)
			}

			result := manager.Get(tt.key, tt.format)
			assert.Equal(t, tt.expectedPath, result)
		})
	}
}

func TestCacheManager_Set(t *testing.T) {
	t.Parallel()

	tempDir := t.TempDir()
	testData := []byte("test image data")

	// Create test file
	testPath := filepath.Join(tempDir, "my_key.jpg")

	tests := []struct {
		name          string
		key           string
		img           []byte
		format        string
		expectedError bool
		expectedPath  string
	}{
		{
			name:          "successful processing with cache set",
			key:           "my_key",
			img:           testData,
			format:        "jpeg",
			expectedError: false,
			expectedPath:  testPath,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			manager, err := NewManager(&Config{CacheDir: tempDir, MaxAge: 1})
			if err != nil {
				t.FailNow()
			}

			path, err := manager.Set(tt.key, tt.img, tt.format)

			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedPath, path)

				savedData, err := os.ReadFile(path)
				assert.NoError(t, err)
				assert.Equal(t, testData, savedData)
			}

		})
	}
}

func TestCacheManager_GetPath(t *testing.T) {
	t.Parallel()

	testDir := "test_cache"

	tests := []struct {
		name     string
		key      string
		format   string
		expected string
	}{
		{
			name:     "jpeg format",
			key:      "test123",
			format:   "jpeg",
			expected: filepath.Join(testDir, "test123.jpg"),
		},
		{
			name:     "png format",
			key:      "test123",
			format:   "png",
			expected: filepath.Join(testDir, "test123.png"),
		},
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	manager, err := NewManager(&Config{CacheDir: testDir, MaxAge: 3600})
	if err != nil {
		t.FailNow()
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := manager.GetPath(tt.key, tt.format)
			assert.Equal(t, tt.expected, result)
		})
	}
}
