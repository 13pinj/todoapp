package session

import (
	"crypto/rand"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/13pinj/todoapp/Godeps/_workspace/src/github.com/gin-gonic/gin"
	"github.com/13pinj/todoapp/models"
)

type dbSession struct {
	SID        string `gorm:"primary_key"`
	CreatedAt  time.Time
	IntMap     map[string]int `sql:"-"`
	EncodedMap string
}

func (s *dbSession) ID() string {
	return s.SID
}

func (s *dbSession) GetInt(key string) int {
	return s.IntMap[key]
}

func (s *dbSession) SetInt(key string, val int) {
	if key == "" {
		return
	}
	s.IntMap[key] = val
	s.encode()
	models.DB.Save(s)
}

func (s *dbSession) encode() {
	val := url.Values{}
	for k, v := range s.IntMap {
		vs := strconv.Itoa(v)
		val.Set(k, vs)
	}
	s.EncodedMap = val.Encode()
}

func (s *dbSession) decode() {
	s.IntMap = make(map[string]int)
	val, err := url.ParseQuery(s.EncodedMap)
	if err != nil {
		return
	}
	for k := range val {
		v, err := strconv.Atoi(val.Get(k))
		if err != nil {
			continue
		}
		s.IntMap[k] = v
	}
}

func dbHasSession(c *gin.Context) (string, bool) {
	st, err := c.Request.Cookie("SessID")
	if err != nil {
		return "", false
	}
	return st.Value, true
}

func dbSessionInit(c *gin.Context) *dbSession {
	session := dbSession{}

	hash := make([]byte, 6)
	rand.Read(hash)

	session.SID = fmt.Sprintf("%x", hash)
	session.IntMap = make(map[string]int)

	models.DB.Create(&session)

	http.SetCookie(c.Writer, &http.Cookie{
		Name:  "SessID",
		Value: session.SID,
		Path:  "/",
	})
	return &session
}

func dbFromContext(c *gin.Context) *dbSession {
	iid, exists := c.Get("session_id")
	if exists {
		id := iid.(string)
		s := &dbSession{}
		err := models.DB.Where("s_id = ?", id).Find(s).Error
		if err == nil {
			s.decode()
			return s
		}
	}
	key, ok := dbHasSession(c)
	if ok {
		s := &dbSession{}
		err := models.DB.Where("s_id = ?", key).Find(s).Error
		if err == nil {
			s.decode()
			c.Set("session_id", key)
			return s
		}
	}
	ses := dbSessionInit(c)
	c.Set("session_id", ses.ID())
	return ses
}

func init() {
	models.DB.AutoMigrate(&dbSession{})
}
