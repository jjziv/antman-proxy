package managers

import (
	"image"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	cacheManagerMock "antman-proxy/managers/cache/mock_manager"
)

func getAllowedDomains() []string {
	return []string{
		"imgur.com",
		"i.imgur.com",
		"flickr.com",
		"staticflickr.com",
		"images.unsplash.com",
	}
}

func TestNewManager(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

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
			name: "cfg.AlllowedDomains is nil",
			config: &Config{
				AllowedDomains: nil,
				CacheManager:   cacheManagerMock.NewMockManager(ctrl),
			},
			expected: "cfg.AllowedDomains is nil!",
		},
		{
			name: "cfg.CacheManager is nil",
			config: &Config{
				AllowedDomains: getAllowedDomains(),
				CacheManager:   nil,
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

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

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
		{
			name:     "allowed i.imgur domain",
			url:      "https://i.imgur.com/image.jpg",
			expected: true,
		},
		{
			name:     "allowed flickr domain",
			url:      "https://flickr.com/image.jpg",
			expected: true,
		},
		{
			name:     "allowed staticflickr.com domain",
			url:      "https://staticflickr.com/image.jpg",
			expected: true,
		},
		{
			name:     "allowed images.unsplash.com domain",
			url:      "https://images.unsplash.com/image.jpg",
			expected: true,
		},
		{
			name:     "disallowed domain",
			url:      "https://example.com/image.jpg",
			expected: false,
		},
		{
			name:     "invalid url",
			url:      "not-a-url",
			expected: false,
		},
	}

	cacheManager := cacheManagerMock.NewMockManager(ctrl)
	manager, err := NewManager(&Config{
		AllowedDomains: getAllowedDomains(),
		CacheManager:   cacheManager,
	})

	if err != nil {
		t.FailNow()
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
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
		quality       int
		expectedError bool
		expectedPath  string
	}{
		{
			name:          "successful processing with cache hit",
			url:           "https://public-assets.getluna.com/images/social/facebook.png",
			width:         100,
			height:        100,
			format:        "jpeg",
			quality:       50,
			expectedError: false,
			expectedPath:  "https://public-assets.getluna.com/images/social/facebook.png",
		},
		{
			name:          "successful processing without cache hit",
			url:           "https://public-assets.getluna.com/images/social/facebook.png",
			width:         100,
			height:        100,
			format:        "jpeg",
			quality:       50,
			expectedError: false,
			expectedPath:  "https://public-assets.getluna.com/images/social/facebook.png",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			cacheManager := cacheManagerMock.NewMockManager(ctrl)
			cacheManager.EXPECT().Get(gomock.Any(), tt.format).Return("")
			cacheManager.EXPECT().Set(gomock.Any(), gomock.Any(), tt.format).Return(tt.expectedPath, nil)

			manager, err := NewManager(&Config{
				AllowedDomains: getAllowedDomains(),
				CacheManager:   cacheManager,
			})
			if err != nil {
				t.FailNow()
			}

			path, err := manager.ProcessImage(tt.url, tt.width, tt.height, tt.format, tt.quality)

			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedPath, path)
			}
		})
	}
}

func TestImageManager_generateCacheKey(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cacheManager := cacheManagerMock.NewMockManager(ctrl)

	manager, err := NewManager(&Config{
		AllowedDomains: getAllowedDomains(),
		CacheManager:   cacheManager,
	})
	if err != nil {
		t.FailNow()
	}

	key := manager.generateCacheKey("http://example.com/image.jpg", 100, 100, "jpeg")
	assert.NotEmpty(t, key)

	// Test consistency
	key2 := manager.generateCacheKey("http://example.com/image.jpg", 100, 100, "jpeg")
	assert.Equal(t, key, key2)
}

// Helper function to create test image
func createTestImage() image.Image {
	img := image.NewRGBA(image.Rect(0, 0, 100, 100))
	return img
}
