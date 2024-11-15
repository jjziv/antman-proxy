package handlers

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestTemplate(t *testing.T) (string, func()) {
	tmpDir, err := os.MkdirTemp("", "template-test")
	require.NoError(t, err)

	templateContent := `
<!DOCTYPE html>
<html>
<head><title>Test Template</title></head>
<body><h1>Test Content</h1></body>
</html>`

	tmpFile := filepath.Join(tmpDir, "index.html")
	err = os.WriteFile(tmpFile, []byte(templateContent), 0666)
	require.NoError(t, err)

	cleanup := func() {
		os.RemoveAll(tmpDir)
	}

	return tmpDir, cleanup
}

func TestHtmlHandler_HandleIndex(t *testing.T) {
	tmpDir, cleanup := setupTestTemplate(t)
	defer cleanup()

	tests := []struct {
		name           string
		expectedStatus int
		expectedHeader map[string]string
		expectBody     string
	}{
		{
			name:           "Successfully renders index page",
			expectedStatus: http.StatusOK,
			expectedHeader: map[string]string{
				"Cache-Control": "public, max-age=300",
			},
			expectBody: "Test Content",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gin.SetMode(gin.TestMode)
			router := gin.New()

			router.LoadHTMLGlob(filepath.Join(tmpDir, "*.html"))

			handler, err := NewHandler()
			require.NoError(t, err)
			require.NotNil(t, handler)

			router.GET("/", handler.HandleIndex)

			w := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, "/", nil)
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			for key, value := range tt.expectedHeader {
				assert.Equal(t, value, w.Header().Get(key))
			}

			if tt.expectBody != "" {
				assert.Contains(t, w.Body.String(), tt.expectBody)
			}
		})
	}
}

func TestHtmlHandler_NewHandler(t *testing.T) {
	tests := []struct {
		name          string
		expectedError bool
	}{
		{
			name:          "Successfully creates new handler",
			expectedError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler, err := NewHandler()

			if tt.expectedError {
				assert.Error(t, err)
				assert.Nil(t, handler)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, handler)

				_, ok := handler.(*HtmlHandler)
				assert.True(t, ok)
			}
		})
	}
}
