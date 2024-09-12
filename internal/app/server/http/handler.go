package http

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"

	"ImageService/internal/pkg/service" // TODO
)

var (
	errInvalidImageBody = fmt.Errorf("invalid image body")
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

	fileHeaders := form.File["image"]
	if len(fileHeaders) != 1 {
		c.Status(http.StatusBadRequest)
		_ = c.Error(errInvalidImageBody)
		return
	}

	// Загружаем файл
	file, err := fileHeaders[0].Open()
	if err != nil {
		slog.Log(c, slog.LevelError, err.Error())
		c.Status(http.StatusInternalServerError)
		return
	}
	defer file.Close()

	fileContent := make([]byte, fileHeaders[0].Size)
	_, err = file.Read(fileContent)
	if err != nil {
		slog.Log(c, slog.LevelError, err.Error())
		c.Status(http.StatusInternalServerError)
		return
	}

	name, err := h.imageSVC.Upload(fileContent)
	if err != nil {
		slog.Log(c, slog.LevelError, err.Error())

		c.String(http.StatusInternalServerError, errInternal.Error())
	}

	c.String(http.StatusOK, "file loaded: %v", name)
}
