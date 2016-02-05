package todolist

import (
	"testing"

	"github.com/13pinj/todoapp/models/todo"
)

func TestNew(t *testing.T) {
	title1 := "First todolist"
	title2 := "Second todolist"

	list1 := New(title1)
	list2 := New(title2)
	empty := New("")

	if list1 == nil || list2 == nil || empty == nil {
		t.Fatal("New() не должно возвращать nil")
	}

	if list1.ID != 0 || list2.ID != 0 || empty.ID != 0 {
		t.Error("ID несохраненных списков не должны быть присвоены")
	}

	if list1.Title != title1 || list2.Title != title2 || empty.Title != "" {
		t.Error("New() должно сохранять заголовок списка")
	}
}

func TestTodoList_Save(t *testing.T) {
	title1 := "First todolist"
	title2 := "Second todolist"

	list1 := New(title1)
	list2 := New(title2)
	empty := New("")

	if list1 == nil || list2 == nil || empty == nil {
		t.FailNow()
	}

	if err := list1.Save(); err != nil {
		t.Error("Сохранение валидного списка не должно вызывать ошибок")
	}
	if err := list2.Save(); err != nil {
		t.Error("Сохранение валидного списка не должно вызывать ошибок")
	}
	if err := empty.Save(); err == nil {
		t.Error("Сохранение списка с пустым заголовком не должно быть допущено")
	}

	if list1.ID == 0 || list2.ID == 0 {
		t.Error("После сохранения списку должен быть присвоен ID")
	}
	if empty.ID != 0 {
		t.Error("Сохранение с ошибкой не должно изменять ID списка")
	}

	if t.Failed() {
		t.FailNow()
	}

	list1_0, ok := Find(list1.ID)
	if !ok {
		t.Fatal("Успешно сохраненный список должен быть найдено по ID")
	}
	if list1_0 == nil {
		t.FailNow()
	}

	list2_0, ok := Find(list2.ID)
	if !ok {
		t.Fatal("Успешно сохраненный список должен быть найдено по ID")
	}
	if list2_0 == nil {
		t.FailNow()
	}

	if list1_0.Title != title1 || list2_0.Title != title2 {
		t.Error("Поле заголовка должно успешно сохраняться")
	}

	if t.Failed() {
		t.FailNow()
	}

	list2_0.Title = ""
	if err := list2_0.Save(); err == nil {
		t.Error("Сохранение прежде сохраненного списка с пустым заголовком не должно быть допущено")
	}
	list2_1, _ := Find(list2_0.ID)
	if list2_1.Title != title2 {
		t.Error("В случае безуспешного сохранения, структура в базе не должна быть изменена")
	}

	list2_1.Title = "Changed text"
	list2_2, _ := Find(list2_1.ID)
	if list2_2.Title != title2 {
		t.Error("В случае несохранения, структура в базе не должна быть изменена")
	}

	list2_1.Save()
	list2_2, _ = Find(list2_1.ID)
	if list2_2.Title != list2_1.Title {
		t.Error("После пересохранения, структура в базе должна быть обновлена")
	}
}

func TestFind(t *testing.T) {
	title1 := "First todolist"
	title2 := "Second todolist"

	list1 := New(title1)
	list2 := New(title2)

	if list1 == nil || list2 == nil {
		t.FailNow()
	}
	if err := list1.Save(); err != nil {
		t.FailNow()
	}
	if err := list2.Save(); err != nil {
		t.FailNow()
	}

	list1_0, ok := Find(list1.ID)
	if !ok {
		t.Fatal("Find() должно находить сохраненные структуры")
	}
	if list1_0 == nil {
		t.Fatal("Find() не должно возвращать nil в случае успешной находки")
	}
	if list1_0.ID != list1.ID || list1_0.Title != list1.Title {
		t.Fatal("Todo, найденного через Find() должно быть эквивалентно искомому")
	}

	list2_0, ok := Find(list2.ID)
	if !ok {
		t.Fatal("Find() должно находить сохраненные структуры")
	}
	if list2_0 == nil {
		t.Fatal("Find() не должно возвращать nil в случае успешной находки")
	}
	if list2_0.ID != list2.ID || list2_0.Title != list2.Title {
		t.Fatal("Todo, найденного через Find() должно быть эквивалентно искомому")
	}

	randomID := uint(1337)
	undef, ok := Find(randomID)
	if undef != nil || ok {
		t.Error("Find() должно отвечать nil и false при поиске несуществующих структур")
	}
}

func TestTodoList_Destroy(t *testing.T) {
	title1 := "First todolist"
	title2 := "Second todolist"

	list1 := New(title1)
	list2 := New(title2)
	unsaved := New("Unsaved")

	if list1 == nil || list2 == nil || unsaved == nil {
		t.FailNow()
	}
	if err := list1.Save(); err != nil {
		t.FailNow()
	}
	if err := list2.Save(); err != nil {
		t.FailNow()
	}

	id1 := list1.ID
	id2 := list2.ID
	list1.Add("First")
	list1.Add("Second")
	todos1 := list1.Todos
	list1.Destroy()
	if _, ok := Find(id1); ok {
		t.Error("Структура уничтоженная с помощью Destroy(), не должна быть найдена в базе")
	}

	for _, t1 := range todos1 {
		if _, ok := todo.Find(t1.ID); ok {
			t.Error("Destroy() должно уничтожать все todo списка.")
			break
		}
	}

	if _, ok := Find(id2); !ok {
		t.Error("Destroy() не должно затрагивать другие структуры в базе и их ID")
	}

	if t.Failed() {
		t.FailNow()
	}

	t.Log("Destroy() должно без паники обрабатывать несохраненные структуры")
	unsaved.Destroy()

	if _, ok := Find(id1); ok {
		t.Error("Destroy() на несохраненной структуре не должно затрагивать другие структуры в базе")
	}
	if _, ok := Find(id2); !ok {
		t.Error("Destroy() на несохраненной структуре не должно затрагивать другие структуры в базе")
	}
}

func TestTodoList_LenDoneUndone(t *testing.T) {
	title1 := "First todolist"
	title2 := "Second todolist"

	list1 := New(title1)
	list2 := New(title2)

	if list1 == nil || list2 == nil {
		t.FailNow()
	}
	if err := list1.Save(); err != nil {
		t.FailNow()
	}
	if err := list2.Save(); err != nil {
		t.FailNow()
	}

	list1.Add("Task 1")
	if len(list1.Todos) < 1 {
		t.FailNow()
	}
	list1.Todos[0].Done = true
	list1.Todos[0].Save()
	list1.Add("Task 2")
	list1.Add("Task 3")

	list2.Add("Task 1")
	list2.Add("Task 2")

	if list1.Len() != 3 {
		t.Errorf("Len() должно возвращать длину списка. Ожидалось %v, получено %v.", 3, list1.Len())
	}
	if list2.Len() != 2 {
		t.Errorf("Len() должно возвращать длину списка. Ожидалось %v, получено %v.", 2, list2.Len())
	}

	if list1.LenDone() != 1 {
		t.Errorf("LenDone() должно возвращать длину списка. Ожидалось %v, получено %v.", 1, list1.LenDone())
	}
	if list2.LenDone() != 0 {
		t.Errorf("LenDone() должно возвращать длину списка. Ожидалось %v, получено %v.", 0, list2.LenDone())
	}

	if list1.LenUndone() != 2 {
		t.Errorf("LenUndone() должно возвращать длину списка. Ожидалось %v, получено %v.", 2, list1.LenUndone())
	}
	if list2.LenUndone() != 2 {
		t.Errorf("LenUndone() должно возвращать длину списка. Ожидалось %v, получено %v.", 2, list2.LenUndone())
	}

	done1 := list1.Done()
	if len(done1) != 1 {
		t.Fatalf("Done() должен возвращать корректное значение. Ожидалась длина %v, получено %v.", 1, len(done1))
	}

	done2 := list2.Done()
	if len(done2) != 0 {
		t.Fatalf("Done() должен возвращать корректное значение. Ожидалась длина %v, получено %v.", 0, len(done1))
	}

	undone1 := list1.Undone()
	if len(undone1) != 2 {
		t.Fatalf("Undone() должен возвращать корректное значение. Ожидалась длина %v, получено %v.", 2, len(undone1))
	}

	undone2 := list2.Undone()
	if len(undone2) != 2 {
		t.Fatalf("Undone() должен возвращать корректное значение. Ожидалась длина %v, получено %v.", 2, len(undone2))
	}
}

func TestTodoList_Add(t *testing.T) {
	title1 := "First todolist"
	list1 := New(title1)

	if list1 == nil {
		t.FailNow()
	}

	if err := list1.Add("First todo"); err == nil {
		t.Error("Add() не должен добавлять задания к несохраненному списку.")
	}

	list1.Save()

	if err := list1.Add("First todo"); err != nil {
		t.Error("Add() не должен отвечать ошибкой на корректное добавление.")
	}
	if list1.Len() != 1 || list1.LenUndone() != 1 {
		t.Error("Дело должно быть добавлено в список (Len() или LenUndone() вернули некорретные значения)")
	}

	if err := list1.Add(""); err == nil {
		t.Error("Add() не должен добавлять задания с пустым текстом.")
	}
}
