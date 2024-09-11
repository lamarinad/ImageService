package main

import (
	"archive/zip"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

// задаем переменную, содержащую путь к директории, куда будут сохраняться загруженные изображения
var imageDir = "./images/"

// uploadHandler — обработчик для загрузки изображений
func uploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Метод не разрешен", http.StatusMethodNotAllowed)
		return
	}

	// Разбирает тело запроса как мультичастный формат, позволяя загружать файлы
	err := r.ParseMultipartForm(10 << 20) // Максимально 10 МБ
	if err != nil {
		http.Error(w, "Ошибка при разборе формы", http.StatusBadRequest)
		return
	}

	file, fileHeader, err := r.FormFile("image")
	if err != nil {
		http.Error(w, "Ошибка при получении файла", http.StatusBadRequest)
		return
	}

	fmt.Println(fileHeader.Filename)

	// Откладывает закрытие файла до конца функции, чтобы избежать утечки ресурсов.
	defer file.Close()

	date := time.Now().Format(time.DateOnly)
	filePath := filepath.Join(imageDir, date)

	err = os.MkdirAll(filePath, os.ModePerm)
	if err != nil {
		http.Error(w, "Ошибка при создании директории", http.StatusInternalServerError)
		return
	}

	// Создание уникального имени файла
	fileName := filepath.Join(filePath, fmt.Sprintf("%d.jpg", time.Now().UnixNano()))

	// Чтение содержимого файла
	content, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, "Ошибка при чтении файла", http.StatusInternalServerError)
		return
	}

	// Запись файла на диск
	err = os.WriteFile(fileName, content, 0644)
	if err != nil {
		http.Error(w, "Ошибка при сохранении файла", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Файл успешно загружен: %s", fileName)
}

// imagesHandler - обработчик для получения изображений за определённый день
func imagesHandler(w http.ResponseWriter, r *http.Request) {
	// Извлекает параметр "date" из URL-запроса
	date := r.URL.Query().Get("date")
	if date == "" {
		http.Error(w, "Дата не указана", http.StatusBadRequest)
		return
	}

	zipfilename := "images_" + date + ".zip"

	archive, err := os.Create(zipfilename)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// сразу закрываем
	archive.Close()

	zipDirectory(fmt.Sprintf("images/%s", date), zipfilename)

	archive, err = os.Open(zipfilename)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer archive.Close()

	content, err := io.ReadAll(archive)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(content)

}

// files, err := os.ReadDir(filepath.Join(imageDir, date))
// if err != nil {
// 	http.Error(w, "Ошибка при чтении директории", http.StatusInternalServerError)
// 	return
// }

// if len(files) != 0 {
// 	buf, err := os.ReadFile(filepath.Join(imageDir, date, files[0].Name()))
// 	if err != nil {
// 		w.WriteHeader(http.StatusInternalServerError)
// 		fmt.Fprint(w, err.Error())
// 		return
// 	}

// 	mimeType := http.DetectContentType(buf)
// 	w.Header().Set("Content-Type", mimeType)
// }

// for _, f := range files {
// 	buf, err := os.ReadFile(filepath.Join(imageDir, date, f.Name()))
// 	if err != nil {
// 		w.WriteHeader(http.StatusInternalServerError)
// 		fmt.Fprint(w, err.Error())
// 		return
// 	}

//		mimeType := http.DetectContentType(buf)
//		w.Header().Set("Content-Type", mimeType)
//		_, err = w.Write(buf)
//		if err != nil {
//			w.WriteHeader(http.StatusInternalServerError)
//			fmt.Fprint(w, err.Error())
//			return
//		}
//	}
func zipDirectory(filepath, zipfilename string) error {
	outFile, err := os.Open(zipfilename)
	if err != nil {
		return err
	}

	w := zip.NewWriter(outFile)

	if err := addFilesToZip(w, zipfilename, ""); err != nil {
		_ = outFile.Close()
		return err
	}

	if err := w.Close(); err != nil {
		_ = outFile.Close()
		return errors.New("Warning: closing zipfile writer failed: " + err.Error())
	}

	if err := outFile.Close(); err != nil {
		return errors.New("Warning: closing zipfile failed: " + err.Error())
	}

	return nil
}

func addFilesToZip(w *zip.Writer, basePath, baseInZip string) error {
	files, err := os.ReadDir(basePath)
	if err != nil {
		return err
	}

	for _, file := range files {
		fullfilepath := filepath.Join(basePath, file.Name())
		if _, err := os.Stat(fullfilepath); os.IsNotExist(err) {
			// ensure the file exists. For example a symlink pointing to a non-existing location might be listed but not actually exist
			continue
		}

		if file.Type().Perm()&os.ModeSymlink != 0 {
			// ignore symlinks alltogether
			continue
		}

		if file.IsDir() {
			if err := addFilesToZip(w, fullfilepath, filepath.Join(baseInZip, file.Name())); err != nil {
				return err
			}
		} else if file.Type().IsRegular() {
			dat, err := os.ReadFile(fullfilepath)
			if err != nil {
				return err
			}
			f, err := w.Create(filepath.Join(baseInZip, file.Name()))
			if err != nil {
				return err
			}
			_, err = f.Write(dat)
			if err != nil {
				return err
			}
		} else {
			// we ignore non-regular files because they are scary
		}
	}
	return nil
}
