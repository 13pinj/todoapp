package session

import "github.com/gin-gonic/gin"

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
func FromContext(c *gin.Context) Session {
	return nil
}
