package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"sync"
	"time"
	"strings"
	"github.com/gorilla/mux"
)

type Shortener struct {
	shortLinkMap map[string]string
	mu           sync.Mutex
}

func NewShortener() *Shortener {
	rand.Seed(time.Now().UnixNano())
	return &Shortener{
		shortLinkMap: make(map[string]string),
	}
}

// Создает строку, случайной дилины, состоящей из случайных символов
func (s *Shortener) GenerateRandomShortURL() string {
	lenght := rand.Intn(9-4+1) + 4
	text := make([]byte, lenght)
	for i := range text {
		text[i] = byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"[rand.Intn(52)])
	}
	return string(text)
}

// Проверка есть ли ссылка в базе
// Проверка сгенерированный url не совпадает с другими
func (s *Shortener) checkLinkShortURL(link string) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, exists := s.shortLinkMap[link]; exists {
		return "", errors.New("Данная ссылка уже существует")
	}
	var newShortURL string
	for {
		newShortURL = s.GenerateRandomShortURL()
		if _, exists := s.shortLinkMap[newShortURL]; !exists {
			break
		}
	}
	return newShortURL, nil
}

func (s *Shortener) PostHandler(w http.ResponseWriter, r *http.Request) {
	//Проверка на POST запрос
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	//Достаем ссылку из запроса
	responseData, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}
	link := string(responseData)
	//Генерируем сокращенную ссылку, и проверяем
	shortURL, err := s.checkLinkShortURL(link)
	if err != nil {
		w.Write([]byte("Ошибка при создании короткой ссылки"))
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

//Проверяет короткий URL из GET запроса
//Возвращает ключ(ссылку в начальном виде)
func (s *Shortener) checkGetShortURL(id string) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for key, value := range s.shortLinkMap {
		if value == id {
			fmt.Println(key)
			return key, nil
		}
	}
	return "", errors.New("Данная ссылка отсутствует")
}

func (s *Shortener) GetHandler(w http.ResponseWriter, r *http.Request) {
	//Проверка на GET запрос
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	link, err := s.checkGetShortURL(strings.TrimPrefix(r.URL.Path, "/"))
	if err != nil {
		fmt.Fprint(w, "Данная ссылка отсутствует")
		return
	}
	w.Header().Set("Location", link)
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func main() {
	shortener := NewShortener()
	router := mux.NewRouter()
	router.HandleFunc("/", shortener.PostHandler).Methods("POST")
	router.HandleFunc("/{id}", shortener.GetHandler).Methods("GET")
	if err := http.ListenAndServe(`:8080`, router); err != nil {
		panic(err)
	}
}