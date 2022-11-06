package dates

import (
	"time"

	"github.com/rickar/cal/v2/us"
)

func GetHolidays(year int, month time.Month) map[int]string {
	result := map[int]string{}

	for _, h := range us.Holidays {
		if h.Month != month {
			continue
		}
		t, o := h.Calc(year)
		result[t.Day()] = h.Name
		if t != o {
			result[o.Day()] = ""
		}
	}

	return result
}
