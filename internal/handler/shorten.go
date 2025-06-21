package handler

import (
	"errors"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	shrterr "github.com/mp1947/ya-url-shortener/internal/errors"
)

// ShortenURL handles the shortening of a given URL sent in the request body.
//
// @Summary      Shorten a URL
// @Description  Accepts a plain text URL in the request body and returns a shortened URL.
// @Tags         shortener
// @Accept       plain
// @Produce      plain
// @Param        url  body      string  true  "Original URL to shorten"
// @Success      201  {string}  string  "Shortened URL"
// @Conflict     409  {string}  string  "Shortened URL already exists"
// @Failure      400  {string}  string  "Invalid request"
// @Failure      500  {string}  string  "Internal server error"
// @Router       /api/shorten [post]
//
// The handler expects a POST request with the original URL in the request body as plain text.
// It returns the shortened URL in plain text format. If the URL has already been shortened,
// it returns a 409 Conflict with the existing shortened URL. For invalid requests or internal
// errors, appropriate HTTP status codes and messages are returned.
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
