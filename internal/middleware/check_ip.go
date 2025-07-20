package middleware

import (
	"net"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mp1947/ya-url-shortener/config"
	"go.uber.org/zap"
)

// WithAuthorizedIP is a middleware that restricts access to requests coming from authorized IP addresses.
// It checks the "X-Real-IP" header of the incoming request and verifies if the IP is within the trusted subnet
// specified in the configuration. If the IP is authorized, the request proceeds to the next handler.
// Otherwise, the middleware logs the unauthorized access attempt and responds with HTTP 401 Unauthorized.
//
// Parameters:
//
//	l      - zap.Logger for logging unauthorized access attempts.
//	cfg    - config.Config containing the trusted subnet information.
//	handler - gin.HandlerFunc to be executed if the IP is authorized.
//
// Returns:
//
//	gin.HandlerFunc - a middleware function for Gin that enforces IP-based authorization.
func WithAuthorizedIP(
	l *zap.Logger,
	cfg config.Config,
	handler gin.HandlerFunc,
) gin.HandlerFunc {
	return func(c *gin.Context) {
		headerData := c.GetHeader("X-Real-IP")

		ipAddr := net.ParseIP(headerData)

		if ipAddr != nil && cfg.TrustedSubnet.Contains(ipAddr) {
			handler(c)
			return
		}

		l.Info("received unauthorized request from ip", zap.String("ip", headerData))

		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": "not authorized to access",
		})
	}
}
