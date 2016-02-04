package user

import "testing"

func TestRegister(t *testing.T) {
	type cred struct {
		name string
		pass string
	}

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
