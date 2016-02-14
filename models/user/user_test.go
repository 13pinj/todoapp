package user

import (
	"os"
	"testing"

	"github.com/13pinj/todoapp/Godeps/_workspace/src/github.com/gin-gonic/gin"
	"github.com/13pinj/todoapp/core/apptesting"
	"github.com/13pinj/todoapp/models"
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

	userForServer   *User
	autoLoginServer *apptesting.Server

	fromCtxServer *apptesting.Server
	logoutServer  *apptesting.Server
)

func TestMain(m *testing.M) {
	// Тестовый сервер для вызова Login
	credForServer = &cred{}
	hf := func(c *gin.Context) {
		lastUser, lastOk = Login(c, credForServer.name, credForServer.pass)
	}
	loginServer = apptesting.NewServer(hf)

	// Тестовый сервер для вызова User.AutoLogin
	userForServer = &User{}
	hf = func(c *gin.Context) {
		userForServer.AutoLogin(c)
	}
	autoLoginServer = apptesting.NewServer(hf)

	// Тестовый сервер для вызова FromContext
	hf = func(c *gin.Context) {
		lastUser, lastOk = FromContext(c)
	}
	fromCtxServer = apptesting.NewServer(hf)

	// Тестовый сервер для вызова Logout
	hf = func(c *gin.Context) {
		Logout(c)
	}
	logoutServer = apptesting.NewServer(hf)
	models.DB.LogMode(false)
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

// Функция выполнения подключения к autoLoginServer
func autoLoginWith(client *apptesting.Client, u *User) {
	userForServer = u
	resp, err := client.Get(autoLoginServer.URL.String())
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

// Функция подключения к logoutServer
func runLogout(client *apptesting.Client) {
	resp, err := client.Get(logoutServer.URL.String())
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
		t.Errorf("Register() не должен регистрировать никнеймы короче 6 символов (%q).", shortNameCred.name)
	}
	shortNameCred = cred{"star", "qwerty1234"}
	_, err = Register(shortNameCred.name, shortNameCred.pass)
	if err != nil {
		t.Errorf("Register() должен успешно регистрировать никнеймы длиной в 4 символа (%q).", shortNameCred.name)
	}
	shortNameCred = cred{"юни", "qwerty1234"}
	_, err = Register(shortNameCred.name, shortNameCred.pass)
	if err == nil {
		t.Errorf("Register() не должен регистрировать никнеймы короче 4 символов Unicode (%q).", shortNameCred.name)
	}

	shortPassCred := cred{"starKiller_1338", "qwert"}
	_, err = Register(shortPassCred.name, shortPassCred.pass)
	if err == nil {
		t.Errorf("Register() не должен регистрировать пароли короче 4 символов (%q).", shortPassCred.pass)
	}
	shortPassCred = cred{"starKiller_1339", "qwerty"}
	_, err = Register(shortPassCred.name, shortPassCred.pass)
	if err != nil {
		t.Errorf("Register() должен успешно регистрировать пароли длиной в 6 символов (%q).", shortPassCred.pass)
	}
	shortPassCred = cred{"starKiller_1340", "парол"}
	_, err = Register(shortPassCred.name, shortPassCred.pass)
	if err == nil {
		t.Errorf("Register() не должен регистрировать пароли короче 6 символов Unicode (%q).", shortPassCred.pass)
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
	if lastUser.Name != cr2.name {
		t.Errorf("Login() должно возвращать правильную структуру после перезаписи входа. Ожидалось %q, получено %q.", cr2.name, lastUser.Name)
	}
}

// С помощью клиента client выполняет регистрация и вход пользователя с данными
// cr на серверах. В случае ошибок в работе Register, Login или FromContext
// прекращает тест.
func prepareTest(t *testing.T, client *apptesting.Client, cr *cred) *User {
	// Проверка Register
	u, err := Register(cr.name, cr.pass)
	if u == nil || err != nil {
		t.Fatal("Register должно работать корректно")
	}

	// Проверка Login
	loginWith(client, cr)
	if !lastOk || lastUser == nil || lastUser.Name != cr.name {
		t.Fatal("Login() и FromContext() должны работать корректно")
	}

	// Проверка FromContext
	runFromContext(client)
	if !lastOk || lastUser == nil || lastUser.Name != cr.name {
		t.Fatal("Login() и FromContext() должны работать корректно")
	}

	return u
}

func TestLogout(t *testing.T) {
	// Подготовить данные для тестовых входов
	cr := cred{"logout_tester", "qwerty1234567"}
	client := apptesting.NewClient()
	prepareTest(t, client, &cr)

	// Вызов Logout после успешного входа
	runLogout(client)
	runFromContext(client)
	if lastOk || lastUser != nil {
		t.Error("Logout() должно стирать данные об аутентификации пользователя")
	}

	// Вызов Logout с пустыми куками
	client.ClearCookie()
	t.Log("Logout() должно корректно реагировать на отсутствие аутентификации")
	runLogout(client)
	runFromContext(client)
	if lastOk || lastUser != nil {
		t.Error("Logout() должно корректно реагировать на отсутствие аутентификации")
	}
}

func TestUser_AutoLogin(t *testing.T) {
	// Подготовить данные для тестовых входов
	cr := cred{"autologin_tester", "qwerty1234567"}
	client := apptesting.NewClient()
	u := prepareTest(t, client, &cr)

	// Запуск AutoLogin на "чистом" клиенте
	client.ClearCookie()
	autoLoginWith(client, u)
	runFromContext(client)
	if !lastOk {
		t.Fatalf("AutoLogin() должно выполнять успешный вход (FromContext() вернуло false)")
	}
	if lastUser == nil {
		t.Fatalf("AutoLogin() должно выполнять успешный вход (FromContext() вернуло nil, true)")
	}
	if lastUser.Name != cr.name {
		t.Errorf("AutoLogin() должно выполнять вход с правильными данными. Ожидалось %q, получено %q.", cr.name, lastUser.Name)
	}

	// Запуск AutoLogin на "занятом" клиенте
	autoLoginWith(client, u)
	runFromContext(client)
	if !lastOk {
		t.Fatalf("AutoLogin() должно выполнять успешный вход с перезаписью (FromContext() вернуло false)")
	}
	if lastUser == nil {
		t.Fatalf("AutoLogin() должно выполнять успешный вход с перезаписью (FromContext() вернуло nil, true)")
	}
	if lastUser.Name != cr.name {
		t.Errorf("AutoLogin() должно выполнять вход с перезаписью с правильными данными. Ожидалось %q, получено %q.", cr.name, lastUser.Name)
	}
}

func TestUser_Destroy(t *testing.T) {
	// Подготовить данные для тестовых входов
	cr := cred{"destroy_tester", "qwerty1234567"}
	client := apptesting.NewClient()
	u := prepareTest(t, client, &cr)

	// Проверка FromContext после удаления пользователя
	u.Destroy()
	runFromContext(client)
	if lastOk || lastUser != nil {
		t.Error("Destroy() должно делать данные об аутентификации пользователя недоступными")
	}

	// Проверка регистрации после удаления
	cr.pass = "New Password"
	client.ClearCookie()
	u, err := Register(cr.name, cr.pass)
	if u == nil || err != nil {
		t.Fatal("Перерегистрация удаленного пользователя должна проходить успешно")
	}

	// Проверка входа после удаления
	u.Destroy()
	loginWith(client, &cr)
	if lastOk || lastUser != nil {
		t.Error("Destroy() должно делать вход невозможным")
	}

	// Проверка AutoLogin после удаления
	client.ClearCookie()
	u, err = Register(cr.name, cr.pass)
	if u == nil || err != nil {
		t.Fatal("Перерегистрация удаленного пользователя должна проходить успешно")
	}
	u.Destroy()
	autoLoginWith(client, u)
	if lastOk || lastUser != nil {
		t.Error("Destroy() должно делать AutoLogin невозможным")
	}
}
