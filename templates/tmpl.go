package tmpl

import "html/template"

var fm = template.FuncMap{
	// pages создает массив номеров страниц, которые следует вывести
	// виджету pager, если пользователь находится на странице номер cur
	// из max штук.
	"pages": func(cur, max int) (ps []int) {
		ps = make([]int, max)
		for i := range ps {
			ps[i] = i + 1
		}
		return
	},
	// rplural возвращает подходящую форму множественного числа
	// для русского слова в сочетании с числом count.
	// Требует три формы слова на вход:
	// one: (один) мяч
	// few: (два) мяча
	// many: (много) мячей
	"rplural": func(count int, one, few, many string) string {
		d := count % 10
		dd := count % 100
		switch {
		case dd >= 11 && dd <= 20:
			return many
		case d == 1:
			return one
		case d >= 2 && d <= 4:
			return few
		default:
			return many
		}
	},
}

func Funcs() template.FuncMap {
	return fm
}
