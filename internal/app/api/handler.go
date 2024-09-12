package api

import (
	"ImageService/internal/pkg/service"
	"fmt"
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"
)

var (
	errInvalidImageBody = fmt.Errorf("invalid")
	errInternal         = fmt.Errorf("internal")
)

type Handler struct {
	imageSVC *service.Image
}

func NewHandler(imageSVC *service.Image) *Handler {
	return &Handler{imageSVC: imageSVC}
}

func (h *Handler) UploadImage(c *gin.Context) {
	form, err := c.MultipartForm()
	if err != nil {
		_ = c.Error(err)

		return
	}

	fileContent, ok := form.Value["image"]
	if !ok || len(fileContent) != 1 {
		c.Status(http.StatusBadRequest)
		_ = c.Error(errInvalidImageBody)

		return
	}

	name, err := h.imageSVC.Upload([]byte(fileContent[0]))
	if err != nil {
		slog.Log(c, slog.LevelError, err.Error())

		c.String(http.StatusInternalServerError, errInternal.Error())
	}

	c.String(http.StatusOK, "file loaded: %v", name)
}
