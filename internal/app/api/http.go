package api

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type HTTP struct {
	port    int
	handler *Handler
}

func NewHTTP(port int, handler *Handler) *HTTP {
	return &HTTP{port: port, handler: handler}
}

func (api *HTTP) ServeHTTP() error {
	r := gin.Default()
	r.GET("/", ping)
	// r.GET("/images", api.handler.) TODO.
	r.POST("/upload", api.handler.UploadImage)

	if err := http.ListenAndServe(fmt.Sprintf(":%d", api.port), r); err != nil {
		return fmt.Errorf("http.ListenAndServe: %w", err)
	}

	return nil
}

// TODO: instead of healtz check.
func ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}
