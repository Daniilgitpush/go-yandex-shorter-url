// go-yandex-shorter-url/cmd/shortener/main.go
package main

import (
	"net/http"
    "github.com/gorilla/mux"
	"github.com/Daniilgitpush/go-yandex-shorter-url/internal/app/handlers"
    "github.com/Daniilgitpush/go-yandex-shorter-url/internal/app/shortener"
)

func main() {
	shortener := NewShortener()
	router := mux.NewRouter()
	router.HandleFunc("/", shortener.PostHandler).Methods("POST")
	router.HandleFunc("/{id}", shortener.GetHandler).Methods("GET")
	if err := http.ListenAndServe(`:8080`, router); err != nil {
		panic(err)
	}
}