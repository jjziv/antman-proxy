package managers

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestCacheManager_Get(t *testing.T) {
	tests := []struct {
		name     string
		key      string
		format   string
		expected string
	}{
		{
			name:     "successful processing with cache hit",
			key:      "my_key",
			format:   "jpg",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			manager := NewManager()

			result := manager.Get(tt.key, tt.format)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestCacheManager_Set(t *testing.T) {
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
			img:           make([]byte, 1),
			format:        "jpg",
			expectedError: false,
			expectedPath:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			manager := NewManager()
			path, err := manager.Set(tt.key, tt.img, tt.format)

			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedPath, path)
			}
		})
	}
}

func TestCacheManager_GetPath(t *testing.T) {
	tests := []struct {
		name     string
		key      string
		format   string
		expected string
	}{
		{
			name:     "successful processing",
			key:      "my_key",
			format:   "jpg",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			manager := NewManager()

			result := manager.GetPath(tt.key, tt.format)
			assert.Equal(t, tt.expected, result)
		})
	}
}
