package app

import (
	"fmt"
	"io"
	"net/http"
)

// HandleOriginal converts provided url to the shorten by generating random Id.
// Returns 400 status code if user sent incorrect request's body, method or content-type.
func (urls *Urls) HandleOriginal(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		body, err := io.ReadAll(r.Body)

		if err != nil || string(body) == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		shortID := generateURLID(randomIDStringLength)

		urls.IDToURL[shortID] = string(body)

		w.WriteHeader(http.StatusCreated)

		w.Write([]byte(fmt.Sprintf("http://localhost:8080/%s", shortID)))
		return
	}

	w.WriteHeader(http.StatusBadRequest)
}

// HandleShort retrieves id from the GET request, looks for
// corresponding url in urls.IDToURL map and redirects to this url via 307 HTTP to the new Location.
// If url wasn't found in urls.IDToURL map or request is incorrect - returns 400 HTTP status code.
func (urls *Urls) HandleShort(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	id := r.PathValue("id")

	originalURL := urls.IDToURL[id]
	if originalURL != "" {
		w.Header().Set("Location", originalURL)
		w.WriteHeader(http.StatusTemporaryRedirect)
		return
	}
	w.WriteHeader(http.StatusBadRequest)
}
