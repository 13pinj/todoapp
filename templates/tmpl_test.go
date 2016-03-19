package tmpl

import (
	"testing"
	"time"
)

func TestAgo(t *testing.T) {
	now := time.Now()
	cases := []struct {
		t   time.Time
		res string
	}{
		{now.Add(-1 * time.Minute), "1 минуту назад"},
		{now.Add(-2 * time.Minute), "2 минуты назад"},
		{now.Add(-59 * time.Minute), "59 минут назад"},
		{now.Add(-1 * time.Hour), "1 час назад"},
		{now.Add(-2 * time.Hour), "2 часа назад"},
		{now.Add(-time.Duration(now.Hour()+1) * time.Hour), "вчера"},
		{now.AddDate(0, 0, -1), "вчера"},
		{now.AddDate(0, 0, -2), "позавчера"},
		{time.Date(2000, time.March, 4, 12, 8, 0, 0, time.Local), "04.03.2000"},
	}

	for _, c := range cases {
		res := ago(c.t)
		if res != c.res {
			t.Errorf("ago(%v): получено %v, ожидалось %v", c.t, res, c.res)
		}
	}
}
