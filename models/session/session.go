package session

import (
	"crypto/rand"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Session interface {
	ID() string

	GetInt(key string) int
	SetInt(key string, val int)

	// TODO: Добавить поддержку разных типов данных.
	// GetFloat(key string) float64
	// SetFloat(key string, val float64)
	//
	// GetBytes(key string) []byte
	// SetBytes(key string, val []byte)
	//
	// GetString(key string) string
	// SetString(key string, val string)
}

// Получает сессию из контекста gin.
// Если данные о сессии не содержатся в куки запроса, создает новую сессию.
// Иначе возвращает созданную ранее.

var mapSession = make(map[string]*imSession)

type imSession struct {
	id string
	mp map[string]int
}

func (s *imSession) ID() string {
	return s.id
}

func (s *imSession) GetInt(key string) int {
	return s.mp[key]
}

func (s *imSession) SetInt(key string, val int) {
	if key == "" {
		return
	}
	s.mp[key] = val
}

func sessionInit(c *gin.Context) *imSession {
	session := imSession{}

	hash := make([]byte, 6)
	rand.Read(hash)

	session.id = fmt.Sprintf("%x", hash)
	session.mp = make(map[string]int)

	mapSession[session.id] = &session

	http.SetCookie(c.Writer, &http.Cookie{
		Name:  "SessID",
		Value: session.id,
	})
	return &session
}

func hasSession(c *gin.Context) (string, bool) {
	st, err := c.Request.Cookie("SessID")
	if err != nil {
		return "", false
	}
	return st.Value, true
}

func FromContext(c *gin.Context) Session {
	if c == nil {
		return nil
	}
	key, ok := hasSession(c)
	if ok {
		ses, ok := mapSession[key]
		if ok {
			return ses
		} else {
			return sessionInit(c)
		}
	}
	return sessionInit(c)
}
