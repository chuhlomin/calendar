package dates

import (
	"testing"
	"time"

	"gotest.tools/assert"
)

var getCalendarTests = []struct {
	year     int
	month    time.Month
	calendar [][7]int
}{
	{
		2020, time.October,
		[][7]int{
			{28, 29, 30, 1, 2, 3, 4},
			{5, 6, 7, 8, 9, 10, 11},
			{12, 13, 14, 15, 16, 17, 18},
			{19, 20, 21, 22, 23, 24, 25},
			{26, 27, 28, 29, 30, 31, 1},
		},
	},
	{
		2020, time.November,
		[][7]int{
			{26, 27, 28, 29, 30, 31, 1},
			{2, 3, 4, 5, 6, 7, 8},
			{9, 10, 11, 12, 13, 14, 15},
			{16, 17, 18, 19, 20, 21, 22},
			{23, 24, 25, 26, 27, 28, 29},
			{30, 1, 2, 3, 4, 5, 6},
		},
	},
}

func Test(t *testing.T) {
	for _, tt := range getCalendarTests {
		assert.DeepEqual(
			t, tt.calendar, GetCalendar(tt.year, tt.month),
		)
	}
}