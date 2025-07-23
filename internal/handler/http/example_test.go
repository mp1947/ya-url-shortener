package handlehttp_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/mp1947/ya-url-shortener/internal/dto"
)

func ExampleHandlerService_DeleteUserURLs() {
	r := gin.New()

	r.Use(func(ctx *gin.Context) {
		ctx.Set("user_id", uuid.NewString())
		ctx.Next()
	})

	r.DELETE("/delete", hs.DeleteUserURLs)

	idsToDeleteReq := `["123431"]`

	req := httptest.NewRequest(
		http.MethodDelete,
		"/delete",
		strings.NewReader(idsToDeleteReq),
	)

	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	fmt.Println(w.Code)

	// Output: 202

}

func ExampleHandlerService_GetOriginalURLByID() {
	r := gin.New()
	r.GET("/:id", hs.GetOriginalURLByID)

	id := "8a56ef86"

	req := httptest.NewRequest(http.MethodGet, "/"+id, nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	fmt.Println(w.Code)

	// Output: 400

}

func ExampleHandlerService_GetUserURLs() {
	r := gin.New()
	r.Use(func(ctx *gin.Context) {
		ctx.Set("user_id", uuid.NewString())
		ctx.Next()
	})
	r.GET("/urls", hs.GetUserURLs)
	req := httptest.NewRequest(http.MethodGet, "/urls", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	fmt.Println(w.Code)

	// Output: 204
}

func ExampleHandlerService_BatchShortenURL() {
	r := gin.New()
	r.Use(func(ctx *gin.Context) {
		ctx.Set("user_id", uuid.NewString())
		ctx.Next()
	})
	correlationID := "aaabbb"
	r.POST("/shorten/batch", hs.BatchShortenURL)
	urls := []dto.BatchShortenRequest{
		{
			CorrelationID: correlationID,
			OriginalURL:   "https://ya.com",
		},
		{
			CorrelationID: correlationID,
			OriginalURL:   "https://goooogle.com",
		},
	}

	urlsByte, _ := json.Marshal(&urls)

	req := httptest.NewRequest(
		http.MethodPost,
		"/shorten/batch",
		bytes.NewReader(urlsByte),
	)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	fmt.Println(w.Code)

	// Output: 201

}

func ExampleHandlerService_JSONShortenURL() {
	r := gin.New()
	r.Use(func(ctx *gin.Context) {
		ctx.Set("user_id", uuid.NewString())
		ctx.Next()
	})
	r.POST("/shorten", hs.JSONShortenURL)

	data := dto.ShortenRequest{
		URL: "https://whatever.com",
	}

	dataByte, _ := json.Marshal(&data)

	req := httptest.NewRequest(http.MethodPost, "/shorten", bytes.NewReader(dataByte))

	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	fmt.Println(w.Code)

	// Output: 201

}
