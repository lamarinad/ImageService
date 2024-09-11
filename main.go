package main

import (
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) { w.Write([]byte("helo")) })
	http.HandleFunc("/upload", uploadHandler)
	http.HandleFunc("/images", imagesHandler)

	go startImageConverter()
	go convertImages()

	log.Println("Сервер запущен на порту 8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Ошибка при запуске сервера: %v", err)
	}
}
