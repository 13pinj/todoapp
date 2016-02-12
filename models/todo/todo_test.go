package todo

import (
	"fmt"
	"math/rand"
	"os"
	"testing"

	"github.com/13pinj/todoapp/models"
)

func TestMain(m *testing.M) {
	models.DB.LogMode(false)
	os.Exit(m.Run())
}

func TestNew(t *testing.T) {
	lbl1 := "First todo"
	lbl2 := "Second todo"

	todo1 := New(lbl1)
	todo2 := New(lbl2)
	empty := New("")

	if todo1 == nil || todo2 == nil || empty == nil {
		t.Fatal("New() не должно возвращать nil")
	}

	if todo1.ID != 0 || todo2.ID != 0 || empty.ID != 0 {
		t.Error("ID несохраненных дел не должны быть присвоены")
	}

	if todo1.TodoListID != 0 || todo2.TodoListID != 0 || empty.TodoListID != 0 {
		t.Error("ID списков несохраненных дел не должны быть присвоены")
	}

	if todo1.Done != false || todo2.Done != false || empty.Done != false {
		t.Error("Свежесозданные дела должны быть отмечены невыполнеными")
	}

	if todo1.Label != lbl1 || todo2.Label != lbl2 || empty.Label != "" {
		t.Error("New() должно сохранять аргумент label в поля Todo")
	}
}

func TestTodo_Save(t *testing.T) {
	lbl1 := "First todo"
	lbl2 := "Second todo"

	todo1 := New(lbl1)
	todo2 := New(lbl2)
	empty := New("")

	if todo1 == nil || todo2 == nil || empty == nil {
		t.FailNow()
	}

	todo1.TodoListID = 12
	if err := todo1.Save(); err != nil {
		t.Error("Сохранение валидного Todo не должно вызывать ошибок")
	}
	todo2.Done = true
	todo2.TodoListID = 14
	if err := todo2.Save(); err != nil {
		t.Error("Сохранение валидного Todo не должно вызывать ошибок")
	}
	if err := empty.Save(); err == nil {
		t.Error("Сохранение Todo с пустым текстом не должно быть допущено")
	}

	if todo1.ID == 0 || todo2.ID == 0 {
		t.Error("После сохранения Todo должен быть присвоен ID")
	}
	if empty.ID != 0 {
		t.Error("Сохранение с ошибкой не должно изменять ID структуры")
	}

	if t.Failed() {
		t.FailNow()
	}

	todo1_0, ok := Find(todo1.ID)
	if !ok {
		t.Fatal("Успешно сохраненное дело должно быть найдено по ID")
	}
	if todo1_0 == nil {
		t.FailNow()
	}

	todo2_0, ok := Find(todo2.ID)
	if !ok {
		t.Fatal("Успешно сохраненное дело должно быть найдено по ID")
	}
	if todo2_0 == nil {
		t.FailNow()
	}

	if todo1_0.Done != todo1.Done || todo2_0.Done != todo2.Done {
		t.Error("Поле Done должно успешно сохраняться")
	}
	if todo1_0.Label != todo1.Label || todo2_0.Label != todo2.Label {
		t.Error("Поле Label должно успешно сохраняться")
	}
	if todo1_0.TodoListID != todo1.TodoListID || todo2_0.TodoListID != todo2.TodoListID {
		t.Error("Поле TodoListID должно успешно сохраняться")
	}

	if t.Failed() {
		t.FailNow()
	}

	todo2_0.Label = ""
	if err := todo2_0.Save(); err == nil {
		t.Error("Сохранение созданного дело с пустым текстом не должно быть допущено")
	}
	todo2_1, _ := Find(todo2_0.ID)
	if todo2_1.Label != lbl2 {
		t.Error("В случае безуспешного сохранения, структура в базе не должна быть изменена")
	}

	todo2_1.Label = "Changed text"
	todo2_2, _ := Find(todo2_1.ID)
	if todo2_2.Label != lbl2 {
		t.Error("В случае несохранения, структура в базе не должна быть изменена")
	}

	todo2_1.Done = !todo2_1.Done
	todo2_1.Save()
	todo2_2, _ = Find(todo2_1.ID)
	if todo2_2.Label != todo2_1.Label || todo2_2.Done != todo2_1.Done {
		t.Errorf("При пересохранении структуру, данные должны обновляться. Ожидалось %#v, получено %#v.", todo2_1, todo2_2)
	}
}

func TestFind(t *testing.T) {
	lbl1 := "First todo"
	lbl2 := "Second todo"

	todo1 := New(lbl1)
	todo2 := New(lbl2)

	if todo1 == nil || todo2 == nil {
		t.FailNow()
	}

	todo2.TodoListID = 12
	if err := todo1.Save(); err != nil {
		t.FailNow()
	}
	todo2.TodoListID = 14
	todo2.Done = true
	if err := todo2.Save(); err != nil {
		t.FailNow()
	}

	todo1_0, ok := Find(todo1.ID)
	if !ok {
		t.Fatal("Find() должно находить сохраненные структуры")
	}
	if todo1_0 == nil {
		t.Fatal("Find() не должно возвращать nil в случае успешной находки")
	}
	if todo1_0.ID != todo1.ID ||
		todo1_0.Done != todo1.Done ||
		todo1_0.Label != todo1.Label ||
		todo1_0.TodoListID != todo1.TodoListID {
		t.Fatal("Todo, найденное через Find(), должно быть эквивалентно искомому")
	}

	todo2_0, ok := Find(todo2.ID)
	if !ok {
		t.Fatal("Find() должно находить сохраненные структуры")
	}
	if todo2_0 == nil {
		t.Fatal("Find() не должно возвращать nil в случае успешной находки")
	}
	if todo2_0.ID != todo2.ID ||
		todo2_0.Done != todo2.Done ||
		todo2_0.Label != todo2.Label ||
		todo2_0.TodoListID != todo2.TodoListID {
		t.Fatal("Todo, найденного через Find() должно быть эквивалентно искомому")
	}

	randomID := uint(1337)
	undef, ok := Find(randomID)
	if undef != nil || ok {
		t.Error("Find() должно отвечать nil и false при поиске несуществующих структур")
	}
}

func TestTodo_Destroy(t *testing.T) {
	todo1 := New("First todo")
	todo2 := New("Second todo")
	unsaved := New("Unsaved")

	if todo1 == nil || todo2 == nil {
		t.FailNow()
	}
	if err := todo1.Save(); err != nil {
		t.FailNow()
	}
	todo2.Done = true
	if err := todo2.Save(); err != nil {
		t.FailNow()
	}

	id1 := todo1.ID
	id2 := todo2.ID
	todo1.Destroy()
	if _, ok := Find(id1); ok {
		t.Error("Структура уничтоженная с помощью Destroy(), не должна быть найдена в базе")
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

func TestFindByList(t *testing.T) {
	type findMap map[uint]bool
	testLists := make(map[uint]findMap)

	for i := 130; i < 135; i++ {
		ui := uint(i)

		testLists[ui] = make(findMap)
		n := rand.Int() % 10

		for j := 0; j < n; j++ {
			todo := New(fmt.Sprintf("Todo #%v", j))
			if todo == nil {
				t.Fatal("New() должно работать корректно (nil)")
			}

			todo.TodoListID = ui
			err := todo.Save()
			if err != nil {
				t.Fatalf("Save() должно работать корректно: %v", err)
			}

			testLists[ui][todo.ID] = false
		}
	}

	for listID, fm := range testLists {
		todos := FindByList(listID)
		for _, td := range todos {
			_, ok := fm[td.ID]
			if !ok {
				t.Errorf("Для списка %v найден лишний элемент %#v", listID, td)
				break
			}
			fm[td.ID] = true
		}

		for id, b := range fm {
			if !b {
				t.Errorf("Для списка %v не найден элемент %#v", listID, id)
			}
		}
	}
}
