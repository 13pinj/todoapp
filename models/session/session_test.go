package session

import (
	"net/http"
	"net/http/cookiejar"
	"os"
	"testing"

	"github.com/13pinj/todoapp/core/apptesting"
	"github.com/gin-gonic/gin"
)

var (
	server  *apptesting.Server
	session Session
)

func TestMain(m *testing.M) {
	hf := func(c *gin.Context) {
		session = FromContext(c)
	}
	server = apptesting.NewServer(hf)
	os.Exit(m.Run())
}

// Отправляет с локального клиента запрос на сервер и возвращает
// полученную для него сессию.
func retrieveSession(client *http.Client) Session {
	if client.Jar == nil {
		jar, _ := cookiejar.New(nil)
		client.Jar = jar
	}

	resp, err := client.Get(server.URL.String())
	if err != nil {
		panic(err)
	}
	resp.Body.Close()
	client.Jar.SetCookies(server.URL, resp.Cookies())

	return session
}

func TestCookieStoring(t *testing.T) {

	sess := FromContext(nil)
	if sess != nil {
		t.Error("FromContext должна возвращать нулевую сессию для нулевого контекста.")
	}

	client := &http.Client{}

	sess1 := retrieveSession(client)
	sess2 := retrieveSession(client)

	client.Jar = nil
	sess3 := retrieveSession(client)

	// Проверки
	if sess1 == nil || sess3 == nil {
		t.Fatal("Сессия должна быть создана, если не найден ключ сессии в куки.")
	}
	if sess2 == nil {
		t.Fatal("Сессия должна быть получена, если ключ найден в куки.")
	}

	if sess1.ID() != sess2.ID() {
		t.Error("Сессии от одинаковых куки должны быть одинаковыми.")
	}
	if sess1.ID() == sess3.ID() || sess2.ID() == sess3.ID() {
		t.Error("Сессии от разных куки должны быть разными.")
	}

}

// Проверяет равенство ключа key в сессии s значению expected.
// Иначе сообщает о нарушении условия cond.
func assertIntKey(t *testing.T, cond string, s Session, key string, expected int) {
	actual := s.GetInt(key)
	if actual != expected {
		t.Errorf("%s Ключ %q: ожидалось %v, получено %v.", cond, key, expected, actual)
	}
}

func TestDataStoring(t *testing.T) {
	client := &http.Client{}

	ints := map[string]int{
		"int1":       12,
		"second int": 34,
	}
	missingInt := "missInt"

	// TODO: Тесты для всех поддерживаемых типов.
	// floats := map[string]float64{
	// 	"float1":       12.04,
	// 	"second float": 3.14,
	// }
	// missingFloat := "missFloat"
	//
	// bytes := map[string][]byte{
	// 	"slice1":       {1, 2, 3, 12},
	// 	"second slice": {},
	// }
	// missingBytes := "missBytes"
	//
	// strings := map[string]string{
	// 	"hello":    "world",
	// 	"username": "onotole",
	// }
	// missingString := "missString"

	sess := retrieveSession(client)
	if sess == nil {
		t.Fatal("Сессия не должна быть nil.")
	}
	for key, val := range ints {
		assertIntKey(t, "Пустая сессия не должна содержать значений int.", sess, key, 0)
		sess.SetInt(key, val)
	}
	sess.SetInt("", 89)

	sess = retrieveSession(client)
	if sess == nil {
		t.Fatal("Сессия не должна быть nil.")
	}
	for key, val := range ints {
		assertIntKey(t, "Сессия должна сохранять значения int.", sess, key, val)
	}
	assertIntKey(t, "Сессия должна возвращать 0 для незаданных значений.", sess, missingInt, 0)
	assertIntKey(t, "Значения с пустыми ключами не должны сохраняться.", sess, "", 0)

	// Обнуление куки.
	client.Jar = nil

	sess = retrieveSession(client)
	if sess == nil {
		t.Fatal("Сессия не должна быть nil.")
	}
	sess.SetInt(missingInt, 12)

	sess = retrieveSession(client)
	if sess == nil {
		t.Fatal("Сессия не должна быть nil.")
	}
	for key := range ints {
		assertIntKey(t, "Сессия не должна возвращать значения чужой сессии.", sess, key, 0)
	}
	assertIntKey(t, "Сессия должна сохранять значения int.", sess, missingInt, 12)

}
