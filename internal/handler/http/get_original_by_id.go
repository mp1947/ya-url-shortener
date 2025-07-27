package handlehttp

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetOriginalURLByID handles GET requests to retrieve the original URL by its shortened ID.
//
// @Summary      Get original URL by ID
// @Description  Redirects to the original URL corresponding to the provided shortened ID.
// @Tags         url
// @Accept       plain
// @Produce      plain
// @Param        id   path      string  true  "Shortened URL ID"
// @Success      307  {string}  string  "Temporary Redirect to the original URL"
// @Failure      400  {string}  string  "Bad Request"
// @Failure      410  {string}  string  "Gone - URL has been deleted"
// @Failure      500  {string}  string  "Internal Server Error"
// @Router       /{id} [get]
func (s HandlerService) GetOriginalURLByID(c *gin.Context) {
	if c.Request.Method != http.MethodGet {
		c.Data(http.StatusBadRequest, contentTypePlain, nil)
		return
	}

	id := c.Param("id")

	url, err := s.Service.GetOriginalURL(c.Request.Context(), id)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	if url.IsDeleted {
		c.Data(http.StatusGone, contentTypePlain, nil)
		return
	}

	if url.OriginalURL != "" {
		c.Header("Location", url.OriginalURL)
		c.Data(http.StatusTemporaryRedirect, contentTypePlain, nil)
		return
	}
	c.Data(http.StatusBadRequest, contentTypePlain, nil)
}
