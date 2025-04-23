package handler

import (
	"errors"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	shrterr "github.com/mp1947/ya-url-shortener/internal/errors"
)

func (s HandlerService) ShortenURL(c *gin.Context) {

	userID, _ := c.Get("user_id")

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
		string(body),
		userID.(string),
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
