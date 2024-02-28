// go-yandex-shorter-url/internal/app/handlers.go
package app

import (
	"fmt"
	"io"
	"net/http"
	"strings"
)

//POST handler
func (s *Shortener) PostHandler(w http.ResponseWriter, r *http.Request) {
	//Достаем ссылку из запроса
	responseData, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading the request", http.StatusInternalServerError)
		return
	}
	link := string(responseData)
	//Генерируем сокращенную ссылку, и проверяем
	shortURL, err := s.checkLinkShortURL(link)
	if err != nil {
		http.Error(w, "Error creating a short link", http.StatusInternalServerError)
		return
	}

	s.mu.Lock()
	s.shortLinkMap[link] = shortURL
	s.mu.Unlock()

	responeseURL := "http://localhost:8080/" + shortURL

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)

	fmt.Fprint(w, responeseURL)
}

//GET handler
func (s *Shortener) GetHandler(w http.ResponseWriter, r *http.Request) {

	link, err := s.checkGetShortURL(strings.TrimPrefix(r.URL.Path, "/"))
	if err != nil {
		http.Error(w, "Link is missing", http.StatusNotFound)
		return
	}
	w.Header().Set("Location", link)
	w.WriteHeader(http.StatusTemporaryRedirect)
}