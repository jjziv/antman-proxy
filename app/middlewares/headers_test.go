package middlewares

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestHeaders(t *testing.T) {
	tests := []struct {
		name            string
		path            string
		expectedHeaders map[string]string
		expectedStatus  int
	}{
		{
			name: "Basic service Headers",
			path: "/basic",
			expectedHeaders: map[string]string{
				"Cache-Control":          "public, max-age=31536000, immutable",
				"CDN-Cache-Control":      "max-age=31536000",
				"Vary":                   "Accept-Encoding",
				"X-Content-Type-Options": "nosniff",
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "CORS Headers",
			path: "/cors",
			expectedHeaders: map[string]string{
				"Cache-Control":               "public, max-age=31536000, immutable",
				"CDN-Cache-Control":           "max-age=31536000",
				"Vary":                        "Accept-Encoding",
				"X-Content-Type-Options":      "nosniff",
				"Access-Control-Allow-Origin": "*",
				"Content-Type":                "image/png",
			},
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gin.SetMode(gin.TestMode)
			router := gin.New()
			router.Use(Headers())

			router.GET("/basic", func(c *gin.Context) {
				c.Status(http.StatusOK)
			})

			router.GET("/cors", func(c *gin.Context) {
				c.Header("Content-Type", "image/png")
				c.Status(http.StatusOK)
			})

			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", tt.path, nil)
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			assert.Equal(t, len(tt.expectedHeaders), len(w.Header()))

			for key, value := range tt.expectedHeaders {
				assert.Equal(t, value, w.Header().Get(key))
			}
		})
	}
}
