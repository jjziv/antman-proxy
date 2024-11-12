package managers

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	cacheManagerMock "antman-proxy/managers/cache/mock_manager"
)

func TestNewManager(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		config   *Config
		expected string
	}{
		{
			name:     "Config is nil",
			config:   nil,
			expected: "ImageManager Config is nil!",
		},
		{
			name: "cfg.CacheDir is nil",
			config: &Config{
				CacheManager: nil,
			},
			expected: "cfg.CacheManager is nil!",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewManager(tt.config)
			assert.Equal(t, tt.expected, err.Error())
		})
	}
}

func TestImageManager_IsURLAllowed(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		url      string
		expected bool
	}{
		{
			name:     "allowed imgur domain",
			url:      "https://imgur.com/image.jpg",
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			cacheManager := cacheManagerMock.NewMockManager(ctrl)

			manager, err := NewManager(&Config{CacheManager: cacheManager})
			if err != nil {
				t.FailNow()
			}

			result := manager.IsURLAllowed(tt.url)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestImageManager_ProcessImage(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		url           string
		width         int
		height        int
		format        string
		expectedError bool
		expectedPath  string
	}{
		{
			name:          "successful processing with cache hit",
			url:           "https://imgur.com/image.jpg",
			width:         100,
			height:        100,
			format:        "jpeg",
			expectedError: false,
			expectedPath:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			cacheManager := cacheManagerMock.NewMockManager(ctrl)

			manager, err := NewManager(&Config{CacheManager: cacheManager})
			if err != nil {
				t.FailNow()
			}

			path, err := manager.ProcessImage(tt.url, tt.width, tt.height, tt.format)

			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedPath, path)
			}
		})
	}
}
