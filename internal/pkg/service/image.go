package service

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"
)

type Image struct {
	imageDir string
}

func NewImage(imageDir string) *Image {
	return &Image{imageDir: imageDir}
}

// Upload TODO: add .jpg check.
func (i *Image) Upload(content []byte) (string, error) {
	var (
		date     = time.Now().Format(time.DateOnly)
		dirPath  = filepath.Join(i.imageDir, date)
		filePath = filepath.Join(dirPath, fmt.Sprintf("%s.jpg", uuid.New().String()))
	)

	if err := os.MkdirAll(dirPath, os.ModePerm); err != nil {
		return "", fmt.Errorf("create directories: %w", err)
	}

	if err := os.WriteFile(filePath, content, 0644); err != nil {
		return "", fmt.Errorf("write file: %w", err)
	}

	return filePath, nil
}
