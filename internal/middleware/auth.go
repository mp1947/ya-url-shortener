package middleware

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/mp1947/ya-url-shortener/internal/auth"
	"go.uber.org/zap"
)

func AuthMiddleware(log *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		authToken := c.GetHeader("Authorization")
		cookie, err := c.Cookie("Authorization")

		var token string

		isCookieExists := errors.Is(err, http.ErrNoCookie)

		if isCookieExists && authToken != "" {
			token = authToken
			log.Info("auth token exists")
		} else if !isCookieExists && authToken == "" {
			token = cookie
			log.Info("auth cookie exists")
		}

		ok, userID := auth.Validate(token)

		if token == "" || !ok {
			generatedUserID := uuid.New()
			token, err := auth.CreateToken(generatedUserID)
			if err != nil {
				log.Warn("error creating new cookie", zap.Error(err))
			}
			c.Header("Authorization", token)
			c.SetCookie("Authorization", token, int(time.Second)*3600, "/", "*", true, true)
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
