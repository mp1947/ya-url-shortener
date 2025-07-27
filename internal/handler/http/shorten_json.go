package handlehttp

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mp1947/ya-url-shortener/internal/dto"
	shrterr "github.com/mp1947/ya-url-shortener/internal/errors"
)

// JSONShortenURL handles the shortening of a URL provided in JSON format.
//
// @Summary      Shorten URL (JSON)
// @Description  Accepts a JSON payload containing a URL and returns a shortened URL.
// @Tags         shortener
// @Accept       json
// @Produce      json
// @Param        request  body      dto.ShortenRequest  true  "URL to shorten"
// @Success      201      {object}  dto.ShortenResponse "Shortened URL created"
// @Conflict     409      {object}  dto.ShortenResponse "URL already shortened"
// @Failure      400      {object}  dto.ShortenResponse "Invalid request"
// @Failure      500      {object}  gin.H               "Internal server error"
// @Router       /api/shorten [post]
// @Security     ApiKeyAuth
//
// It expects a JSON body with the original URL, validates the input, and returns the shortened URL.
// If the URL has already been shortened, it returns a 409 Conflict with the existing short URL.
// On invalid input, it returns a 400 Bad Request. On server errors, it returns a 500 Internal Server Error.
func (s HandlerService) JSONShortenURL(c *gin.Context) {
	var request dto.ShortenRequest
	rawRequest, err := c.GetRawData()

	userID, _ := c.Get("user_id")

	if err != nil {
		c.JSON(
			http.StatusBadRequest,
			dto.ShortenResponse{Result: requestBodyGetErr},
		)
		return
	}

	if unmarshalErr := json.Unmarshal(rawRequest, &request); unmarshalErr != nil {

		c.JSON(
			http.StatusBadRequest,
			dto.ShortenResponse{Result: requestBindingErr},
		)
		return
	}

	shortURL, err := s.Service.ShortenURL(
		c.Request.Context(),
		string(request.URL),
		fmt.Sprintf("%s", userID),
	)

	if errors.Is(err, shrterr.ErrOriginalURLAlreadyExists) {
		c.JSON(http.StatusConflict, dto.ShortenResponse{Result: shortURL})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "internal server error while shorten url",
		})
		return
	}

	c.JSON(http.StatusCreated, dto.ShortenResponse{Result: shortURL})

}
