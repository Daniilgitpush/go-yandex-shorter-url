// go-yandex-shorter-url/internal/app/shortener/shortener.go
package shortener

import (
	"errors"
	"math/rand"
	"sync"
	"time"
)

type Shortener struct {
	shortLinkMap map[string]string
	mu           sync.Mutex
	rand         *rand.Rand
}

func NewShortener() *Shortener {
	return &Shortener{
		shortLinkMap: make(map[string]string),
		rand:         rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

func (s *Shortener) GenerateRandomShortURL(rand *rand.Rand) string {
	length := rand.Intn(9-4+1) + 4
	text := make([]byte, length)
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
		return "", errors.New("the link already exists")
	}
	var newShortURL string
	for {
		newShortURL = s.GenerateRandomShortURL(s.rand)
		if _, exists := s.shortLinkMap[newShortURL]; !exists {
			break
		}
	}
	return newShortURL, nil
}

//Проверяет короткий URL из GET запроса
//Возвращает ключ(ссылку в начальном виде)
func (s *Shortener) checkGetShortURL(id string) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for key, value := range s.shortLinkMap {
		if value == id {
			return key, nil
		}
	}
	return "", errors.New("link is missing")
}

