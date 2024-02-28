// go-yandex-shorter-url/cmd/shortener/main.go
package main

import (
	"github.com/Daniilgitpush/go-yandex-shorter-url/internal/app"
	"net/http"
    "github.com/gorilla/mux"
)

func main() {
	shortener := app.NewShortener()
	router := mux.NewRouter()
	router.HandleFunc("/", shortener.PostHandler).Methods("POST")
	router.HandleFunc("/{id}", shortener.GetHandler).Methods("GET")
	if err := http.ListenAndServe(`:8080`, router); err != nil {
		panic(err)
	}
}