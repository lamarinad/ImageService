package main

import (
	"fmt"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
	"time"
)

// создаем новый таймер, который будет срабатывать каждые 10 минут
func startImageConverter() {
	ticker := time.NewTicker(10 * time.Minute)
	// гарантирует остановку таймера при выходе из функции, чтобы освободить ресурсы.
	defer ticker.Stop()

	// бесконечный цикл, ожидающий сигнала от таймера,
	// после чего вызывает функцию convertImages, которая выполняет конвертацию изображений
	for {
		<-ticker.C
		convertImages()
	}
}

func convertImages() {
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

	err := filepath.Walk(imageDir, walkFunc)
	if err != nil {
		fmt.Printf("Ошибка при конвертации изображений: %v\n", err)
	}
}

// вспомогательная функция
// принимает путь к файлу и новое расширение. Заменяет текущее расширение на новое, сохраняя при этом имя файла.
// таким образом превращает image.jpg в image.png.
func changeExt(path, newExt string) string {
	return path[:len(path)-len(filepath.Ext(path))-1] + "." + newExt
}
