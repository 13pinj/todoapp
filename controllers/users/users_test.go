package users

import (
	"net/url"
	"os"
	"testing"

	"github.com/13pinj/todoapp/Godeps/_workspace/src/github.com/gin-gonic/gin"
	"github.com/13pinj/todoapp/core/apptesting"
	"github.com/13pinj/todoapp/models/user"
)

var (
	regServer, loginServer, logoutServer, desServer *apptesting.Server

	cu *user.User
	ok bool
)

func TestMain(m *testing.M) {
	hf := func(c *gin.Context) {
		Register(c)
		cu, ok = user.FromContext(c)
	}
	regServer = apptesting.NewServer(hf)

	hf = func(c *gin.Context) {
		Login(c)
		cu, ok = user.FromContext(c)
	}
	loginServer = apptesting.NewServer(hf)

	hf = func(c *gin.Context) {
		Logout(c)
		cu, ok = user.FromContext(c)
	}
	logoutServer = apptesting.NewServer(hf)

	hf = func(c *gin.Context) {
		Destroy(c)
		cu, ok = user.FromContext(c)
	}
	desServer = apptesting.NewServer(hf)

	os.Exit(m.Run())
}

func TestRegister(t *testing.T) {
	client := apptesting.NewClient()
	creds := url.Values{
		"name":     {"bert_maklin"},
		"password": {"my_password"},
	}
	resp, err := client.PostForm(regServer.URL.String(), creds)
	if err != nil {
		panic(err)
	}
	resp.Body.Close()

	if !ok || cu == nil {
		t.Fatal("Register должен выполнять успешную регистрацию и вход")
	}
	if cu.Name != creds.Get("name") {
		t.Errorf("Register должен выполнять регистрацию под правильным именем. Ожидалось %v, получено %v.", creds.Get("name"), cu.Name)
	}

	creds1 := url.Values{
		"name":     {"bert_maklin_rereg"},
		"password": {"my_password"},
	}
	resp, err = client.PostForm(regServer.URL.String(), creds1)
	if err != nil {
		panic(err)
	}
	resp.Body.Close()

	if !ok || cu == nil {
		t.Fatal("Register не должен перезаписывать сессию")
	}
	if cu.Name != creds.Get("name") {
		t.Errorf("Register не должен перезаписывать сессию. Ожидалось %v, получено %v.", creds.Get("name"), cu.Name)
	}

	client.ClearCookie()

	resp, _ = client.PostForm(loginServer.URL.String(), creds1)
	resp.Body.Close()

	if ok || cu != nil {
		t.Error("Register не должен выполнять регистрацию при попытке перезаписи сессии")
	}

	client.ClearCookie()

	creds2 := url.Values{
		"name":     {""},
		"password": {""},
	}
	resp, err = client.PostForm(regServer.URL.String(), creds2)
	if err != nil {
		panic(err)
	}
	resp.Body.Close()

	if ok || cu != nil {
		t.Fatal("Register не должен выполнять регистрацию ошибочных данных")
	}
}

func TestLogin(t *testing.T) {
	client := apptesting.NewClient()
	creds := url.Values{
		"name":     {"bert_maklin_2"},
		"password": {"my_password"},
	}
	resp, _ := client.PostForm(regServer.URL.String(), creds)
	resp.Body.Close()

	client.ClearCookie()

	resp, err := client.PostForm(loginServer.URL.String(), creds)
	if err != nil {
		panic(err)
	}
	resp.Body.Close()

	if !ok || cu == nil {
		t.Fatal("Login должен выполнять успешный вход")
	}
	if cu.Name != creds.Get("name") {
		t.Errorf("Login должен выполнять вход под правильным именем. Ожидалось %v, получено %v.", creds.Get("name"), cu.Name)
	}

	client.ClearCookie()

	creds2 := url.Values{
		"name":     {"bert_maklin_2"},
		"password": {"my_passworddd"},
	}

	resp, err = client.PostForm(loginServer.URL.String(), creds2)
	if err != nil {
		panic(err)
	}
	resp.Body.Close()

	if ok || cu != nil {
		t.Fatal("Login не должен выполнять вход с ошибочными данными")
	}
}

func TestLogout(t *testing.T) {
	client := apptesting.NewClient()
	creds := url.Values{
		"name":     {"bert_maklin_3"},
		"password": {"my_password"},
	}
	resp, _ := client.PostForm(regServer.URL.String(), creds)
	resp.Body.Close()

	resp, err := client.PostForm(logoutServer.URL.String(), nil)
	if err != nil {
		panic(err)
	}
	resp.Body.Close()

	if ok || cu != nil {
		t.Fatal("Logout должен выполнять успешный выход")
	}

	client.ClearCookie()

	t.Log("Запускаю Logout с пустыми куки...")
	resp, err = client.PostForm(logoutServer.URL.String(), nil)
	if err != nil {
		panic(err)
	}
	resp.Body.Close()

}

func TestDestroy(t *testing.T) {
	client := apptesting.NewClient()
	creds := url.Values{
		"name":     {"bert_maklin_4"},
		"password": {"my_password"},
	}
	resp, _ := client.PostForm(regServer.URL.String(), creds)
	resp.Body.Close()

	resp, err := client.PostForm(desServer.URL.String(), nil)
	if err != nil {
		panic(err)
	}
	resp.Body.Close()
	if ok || cu != nil {
		t.Fatal("Destroy должен стирать данные об аутентификации")
	}

	resp, _ = client.PostForm(loginServer.URL.String(), creds)
	resp.Body.Close()
	if ok || cu != nil {
		t.Fatal("Destroy должен делать последующий вход невозможным")
	}
}
