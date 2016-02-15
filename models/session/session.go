package session

import "github.com/13pinj/todoapp/Godeps/_workspace/src/github.com/gin-gonic/gin"

// Session - интерфейс структуры сессии.
// Все реализации сессий в проекте будут удовлетворять этому интерфейсу.
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

// FromContext получает сессию из контекста gin.
// Если данные о сессии не содержатся в куки запроса, создает новую сессию.
// Иначе возвращает созданную ранее.
func FromContext(c *gin.Context) Session {
	if c == nil {
		return nil
	}
	return imFromContext(c)
}
