package tmpl

import (
	"fmt"
	"html/template"
	"time"
)

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
	"rplural": rplural,
	// ago формирует строку типа "X часов назад" из обозначенной точки во времени.
	"ago": ago,
}

// rplural доступна в шаблонах через fm["rplural"]
func rplural(count int, one, few, many string) string {
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
}

// ago доступна в шаблонах через fm["ago"]
func ago(t time.Time) string {
	now := time.Now()
	diff := now.Sub(t)
	hours := int(diff.Hours())
	days := int(now.Truncate(24*time.Hour).Sub(t.Truncate(24*time.Hour)).Hours() / 24)
	switch {
	case hours == 0:
		m := diff.Minutes()
		return fmt.Sprintf("%d %s назад", int(m), rplural(int(m), "минуту", "минуты", "минут"))
	case days == 0 && hours < 24:
		h := diff.Hours()
		return fmt.Sprintf("%d %s назад", int(h), rplural(int(h), "час", "часа", "часов"))
	case days == 1:
		return "вчера"
	case days == 2:
		return "позавчера"
	default:
		return t.Format("02.01.2006")
	}
}

func Funcs() template.FuncMap {
	return fm
}
