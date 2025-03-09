package main

import (
	"github.com/gin-gonic/gin"
	"github.com/mp1947/ya-url-shortener/internal/app"
)

func main() {

	urls := &app.Urls{IDToURL: map[string]string{}}
	r := gin.Default()

	r.Any("/", urls.HandleOriginal)
	r.Any("/:id", urls.HandleShort)

	if err := r.Run(":8080"); err != nil {
		panic(err)
	}
}
