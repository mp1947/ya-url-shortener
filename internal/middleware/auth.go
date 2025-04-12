package middleware

import (
	"errors"
	"net/http"
	"net/http/httputil"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/mp1947/ya-url-shortener/internal/auth"
	"go.uber.org/zap"
)

func AuthMiddleware(log *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenCookie, err := c.Cookie("token")

		log.Info("print all cookies", zap.Any("cookies", c.Request.Cookies()))

		rDump, _ := httputil.DumpRequest(c.Request, true)

		log.Info("print all headers", zap.Any("request data", string(rDump)))

		ok, userID := auth.Validate(tokenCookie)

		if errors.Is(err, http.ErrNoCookie) || !ok {
			generatedUserID := uuid.New()
			newCookie, err := auth.CreateCookie(generatedUserID)
			if err != nil {
				log.Warn("error creating new cookie", zap.Error(err))
			}
			c.SetCookie("token", newCookie, int(time.Second)*3600, "/", "*", true, true)
			c.Set("user_id", generatedUserID.String())
			c.Next()
			return
		}
		userIDStr := userID.String()
		log.Info("processing request from user", zap.String("user_id", userIDStr))
		c.Set("user_id", userIDStr)
		c.Next()
	}
}
