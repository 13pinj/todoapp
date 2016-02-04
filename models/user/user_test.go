package user

import (
	"os"
	"testing"

	"github.com/13pinj/todoapp/core/apptesting"
	"github.com/gin-gonic/gin"
)

type cred struct {
	name string
	pass string
}

var (
	lastOk   bool
	lastUser *User

	credForServer *cred
	loginServer   *apptesting.Server
	fromCtxServer *apptesting.Server
)

func TestMain(m *testing.M) {
	// Тестовый сервер для вызова Login
	credForServer = &cred{}
	hf := func(c *gin.Context) {
		lastUser, lastOk = Login(c, credForServer.name, credForServer.pass)
	}
	loginServer = apptesting.NewServer(hf)

	// Тестовый сервер для вызова FromContext
	hf = func(c *gin.Context) {
		lastUser, lastOk = FromContext(c)
	}
	fromCtxServer = apptesting.NewServer(hf)

	os.Exit(m.Run())
}

// Функция выполнения подключения к loginServer
func loginWith(client *apptesting.Client, c *cred) {
	credForServer = c
	resp, err := client.Get(loginServer.URL.String())
	if err != nil {
		panic(err)
	}
	resp.Body.Close()
}

// Функция подключения к fromCtxServer
func runFromContext(client *apptesting.Client) {
	resp, err := client.Get(fromCtxServer.URL.String())
	if err != nil {
		panic(err)
	}
	resp.Body.Close()
}

func TestRegister(t *testing.T) {
	validCred := cred{"starKiller_1337", "my secure password"}
	u, err := Register(validCred.name, validCred.pass)
	if err != nil {
		t.Fatalf("Register() не должен отвечать ошибкой на корректный ввод (%#v)", validCred)
	}
	if u == nil {
		t.Fatalf("Register() не должен возвращать нулевой указатель в случае успешной регистрации (%#v)", validCred)
	}
	if u.Name != validCred.name {
		t.Errorf("Поле имени в возвращаемой структуре должно быть равно введенному. Ожидалось %q, получено %q.", validCred.name, u.Name)
	}
	if u.PwdHash == validCred.pass {
		t.Error("Пароль не должен храниться в структуре в открытом виде")
	}

	takedCred := cred{"starKiller_1337", "qwerty1234"}
	_, err = Register(validCred.name, validCred.pass)
	if err == nil {
		t.Errorf("Register() не должен регистрировать занятые никнеймы (%q).", takedCred.name)
	}

	credWithSpace := cred{"starKiller 1337", "qwerty1234"}
	_, err = Register(credWithSpace.name, credWithSpace.pass)
	if err == nil {
		t.Errorf("Register() не должен регистрировать никнеймы, содержащие пробелы (%q).", credWithSpace.name)
	}

	shortNameCred := cred{"sta", "qwerty1234"}
	_, err = Register(shortNameCred.name, shortNameCred.pass)
	if err == nil {
		t.Errorf("Register() не должен регистрировать никнеймы короче 4 символов (%q).", shortNameCred.name)
	}

	shortPassCred := cred{"starKiller_1337", "qwert"}
	_, err = Register(shortPassCred.name, shortPassCred.pass)
	if err == nil {
		t.Errorf("Register() не должен регистрировать пароли короче 4 символов (%q).", shortPassCred.pass)
	}
}

func TestLoginFromContext(t *testing.T) {

	// Подготовить данные для тестовых входов
	cr := cred{"login_tester", "qwerty1234567"}
	u, err := Register(cr.name, cr.pass)
	if u == nil || err != nil {
		t.Fatal("Register должно работать корректно")
	}

	cr2 := cred{"login_tester_2", "qwerty1234567"}
	u, err = Register(cr2.name, cr2.pass)
	if u == nil || err != nil {
		t.Fatal("Register должно работать корректно")
	}

	cr3 := cred{}

	client := apptesting.NewClient()

	// Login с корректными данными
	loginWith(client, &cr)
	if !lastOk {
		t.Fatalf("Login() должно сообщать об успешной авторизации при получении корректных данных для входа (%#v)", cr)
	}
	if lastUser == nil {
		t.Fatalf("Login() не должно возвращать нулевой указатель в случае успешного входа (%#v)", cr)
	}
	if lastUser.Name != cr.name {
		t.Errorf("Login() должно возвращать правильную структуру. Ожидалось %q, получено %q.", cr.name, lastUser.Name)
	}

	// FromContext после корректного входа
	runFromContext(client)
	if !lastOk {
		t.Fatal("FromContext() должно сообщать об успехе после Login.")
	}
	if lastUser == nil {
		t.Fatalf("FromContext() не должно возвращать нулевой указатель в случае успешного входа (%#v)", cr)
	}
	if lastUser.Name != cr.name {
		t.Errorf("FromContext() должно возвращать правильную структуру. Ожидалось %q, получено %q.", cr.name, lastUser.Name)
	}

	// FromContext при пустых куки
	client.ClearCookie()
	runFromContext(client)
	if lastOk {
		t.Fatal("FromContext() должно сообщать о неуспехе при пустых куки.")
	}
	if lastUser != nil {
		t.Fatalf("FromContext() не должно возвращать структуру пользователя при пустых куки")
	}

	// Login с неверным паролем
	cr3 = cr
	cr3.pass = "incorrect password"
	loginWith(client, &cr3)
	if lastOk {
		t.Fatal("Login() должно сообщать о неуспехе при неверном пароле.")
	}
	if lastUser != nil {
		t.Fatalf("Login() не должно возвращать структуру пользователя при неверном пароле")
	}

	// FromContext после неверного пароля
	runFromContext(client)
	if lastOk {
		t.Fatal("FromContext() должно сообщать о неуспехе после вызова Login() с неверным паролем.")
	}
	if lastUser != nil {
		t.Fatalf("FromContext() не должно возвращать структуру пользователя после вызова Login() с неверным паролем")
	}

	// Login с неверным никнеймом
	cr3 = cr
	cr3.name = "undefined_user"
	loginWith(client, &cr3)
	if lastOk {
		t.Fatal("Login() должно сообщать о неуспехе при неверном никнейме.")
	}
	if lastUser != nil {
		t.Fatalf("Login() не должно возвращать структуру пользователя при неверном никнейме")
	}

	// Проверка перезаписи данных Login`ом
	loginWith(client, &cr)
	loginWith(client, &cr2)
	if !lastOk {
		t.Fatalf("Login() должно сообщать об успешной авторизации после записи входа (%#v)", cr)
	}
	if lastUser == nil {
		t.Fatalf("Login() не должно возвращать нулевой указатель после перезаписи входа (%#v)", cr)
	}
	if lastUser.Name != cr.name {
		t.Errorf("Login() должно возвращать правильную структуру после перезаписи входа. Ожидалось %q, получено %q.", cr.name, lastUser.Name)
	}
}
