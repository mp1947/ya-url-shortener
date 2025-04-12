package handler

import (
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	shrterr "github.com/mp1947/ya-url-shortener/internal/errors"
)

func (s HandlerService) ShortenURL(c *gin.Context) {

	userID, exists := c.Get("user_id")

	if !exists {
		userID = uuid.New().String()
	}

	if c.Request.Method != http.MethodPost {
		c.Data(http.StatusBadRequest, contentTypePlain, nil)
		return
	}

	body, err := io.ReadAll(c.Request.Body)

	if err != nil || string(body) == "" {
		c.Data(http.StatusBadRequest, contentTypePlain, nil)
		return
	}

	shortURL, err := s.Service.ShortenURL(
		c.Request.Context(),
		s.Cfg,
		string(body),
		fmt.Sprintf("%s", userID),
	)

	if errors.Is(err, shrterr.ErrOriginalURLAlreadyExists) {
		c.Data(http.StatusConflict, contentTypePlain, []byte(shortURL))
		return
	} else if err != nil {
		c.Data(
			http.StatusInternalServerError,
			contentTypePlain,
			[]byte("internal server error while shorten url"),
		)
		return
	}

	c.Data(http.StatusCreated, contentTypePlain, []byte(shortURL))

}
