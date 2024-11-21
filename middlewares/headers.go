package middlewares

import (
	"strings"

	"github.com/gin-gonic/gin"
)

// Returns the headers the service needs in order to enable Browser + CDN caching. Additionally, includes any necessary security headers.
func getCdnHeaders(c *gin.Context) map[string]string {
	headers := map[string]string{
		"Cache-Control":          "public, max-age=31536000, immutable",
		"CDN-Cache-Control":      "max-age=31536000",
		"Vary":                   "Accept-Encoding",
		"X-Content-Type-Options": "nosniff",
	}

	// Enables CORS
	if strings.HasPrefix(c.Writer.Header().Get("Content-Type"), "image/") {
		headers["Access-Control-Allow-Origin"] = "*"
	}

	return headers
}

func Headers() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		for key, value := range getCdnHeaders(c) {
			c.Header(key, value)
		}
	}
}
