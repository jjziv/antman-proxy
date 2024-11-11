package managers

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestImageManager_IsURLAllowed(t *testing.T) {
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

			manager := NewManager()

			result := manager.IsURLAllowed(tt.url)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestImageManager_ProcessImage(t *testing.T) {
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

			manager := NewManager()
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
