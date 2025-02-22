package app

import (
	"fmt"
	"net/http"

	"github.com/mp1947/ya-url-shortener/config"
)

func generateShortURL(req *http.Request, c config.Config, id string) string {
	requestHost := req.Host
	basePath := *c.BasePath
	trailingSlash := "/"

	if string(basePath[len(basePath)-1]) != trailingSlash {
		basePath += "/"
	}

	return fmt.Sprintf("http://%s%s%s", requestHost, basePath, id)
}
