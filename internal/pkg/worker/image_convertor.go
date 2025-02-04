package worker

import (
	"context"
	"fmt"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
	"time"
)

type imageConvertor struct {
	imageDir string
	chDone   chan struct{}
}

func NewImageConvertor(imageDir string) Worker {
	return &imageConvertor{imageDir: imageDir, chDone: make(chan struct{})}
}

func (ic *imageConvertor) Start(ctx context.Context) error {
	ticker := time.NewTicker(10 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			ic.convertImages()
		case <-ic.chDone:
			return nil
		}
	}
}

// Stop TODO add better graceful worker shutdown
func (ic *imageConvertor) Stop() error {
	ic.chDone <- struct{}{}
	ic.chDone <- struct{}{}
	return nil
}

func (ic *imageConvertor) convertImages() {
	// рекурсивно проходит по всем файлам и папкам в директории imageDir

	walkFunc := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// обработка изображений
		// проверяет, является ли файл изображением в формате JPEG
		if filepath.Ext(path) == ".jpg" {
			file, err := os.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()
			// декодирует содержимое JPEG-файла в изображение
			img, err := jpeg.Decode(file)
			if err != nil {
				return err
			}

			// cоздание PNG-файла
			// создает новый файл с тем же именем, но с расширением .png
			pngFile, err := os.Create(changeExt(path, "png"))
			if err != nil {
				return err
			}
			defer pngFile.Close()

			// кодирует изображение в формате PNG и записывает его в новый файл
			return png.Encode(pngFile, img)
		}
		return nil
	}

	err := filepath.Walk(ic.imageDir, walkFunc)
	if err != nil {
		fmt.Printf("Ошибка при конвертации изображений: %v\n", err)
	}
}

// вспомогательная функция
// принимает путь к файлу и новое расширение. Заменяет текущее расширение на новое, сохраняя при этом имя файла.
// таким образом превращает image.jpg в image.png.
func changeExt(path, newExt string) string {
	// TODO: это так не работает, нужно перекодировать весь файл из jpg в png.
	return path[:len(path)-len(filepath.Ext(path))-1] + "." + newExt
}
