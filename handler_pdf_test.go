package main

import (
	"strings"
	"testing"
	"time"

	"github.com/go-test/deep"
)

func (d days) Strings(firstDayOfWeek int32) []string {
	lastDayOfWeek := (firstDayOfWeek + 6) % 7

	var result []string
	var lineHasContent bool

	line := make([]string, 7)
	for _, day := range d {
		weekday := int32(day.Date.Weekday())

		// column = 0 if the first day of the week is monday and day is monday
		// column = 6 if the first day of the week is monday and day is sunday
		// column = 0 if the first day of the week is sunday and day is sunday
		// column = 6 if the first day of the week is sunday and day is saturday
		column := (weekday - firstDayOfWeek + 7) % 7

		dayStr := day.Date.Format("2")
		if day.Inactive {
			line[column] = "_" + dayStr + "_"
		} else if day.Weekend {
			line[column] = "*" + dayStr + "*"
		} else {
			line[column] = dayStr
		}
		lineHasContent = true

		if weekday == lastDayOfWeek {
			result = append(result, strings.Join(line, " "))
			line = make([]string, 7)
			lineHasContent = false
		}
	}
	if lineHasContent {
		result = append(result, strings.Join(line, " "))
	}

	return result
}

func TestGetDays(t *testing.T) {
	var tests = []struct {
		name                string
		year                int32
		month               time.Month
		firstDay            int32
		showInactiveDays    bool
		expectedDays        []string
		expectedWeekNumbers []int
	}{
		{
			name:             "sunday",
			firstDay:         0, // sunday
			showInactiveDays: false,
			year:             2022,
			month:            12,
			expectedDays: []string{
				"    1 2 *3*",
				"*4* 5 6 7 8 9 *10*",
				"*11* 12 13 14 15 16 *17*",
				"*18* 19 20 21 22 23 *24*",
				"*25* 26 27 28 29 30 *31*",
			},
			expectedWeekNumbers: []int{48, 49, 50, 51, 52},
		},
		{
			name:             "sunday, show inactive days",
			firstDay:         0, // sunday
			showInactiveDays: true,
			year:             2022,
			month:            12,
			expectedDays: []string{
				"_27_ _28_ _29_ _30_ 1 2 *3*",
				"*4* 5 6 7 8 9 *10*",
				"*11* 12 13 14 15 16 *17*",
				"*18* 19 20 21 22 23 *24*",
				"*25* 26 27 28 29 30 *31*",
				"_1_ _2_ _3_ _4_ _5_ _6_ _7_",
			},
			expectedWeekNumbers: []int{48, 49, 50, 51, 52, 1},
		},
		{
			name:             "monday",
			firstDay:         1, // monday
			showInactiveDays: false,
			year:             2022,
			month:            12,
			expectedDays: []string{
				"   1 2 *3* *4*",
				"5 6 7 8 9 *10* *11*",
				"12 13 14 15 16 *17* *18*",
				"19 20 21 22 23 *24* *25*",
				"26 27 28 29 30 *31* ",
			},
			expectedWeekNumbers: []int{48, 49, 50, 51, 52},
		},
		{
			name:             "monday, show inactive days",
			firstDay:         1, // monday
			showInactiveDays: true,
			year:             2022,
			month:            12,
			expectedDays: []string{
				"_28_ _29_ _30_ 1 2 *3* *4*",
				"5 6 7 8 9 *10* *11*",
				"12 13 14 15 16 *17* *18*",
				"19 20 21 22 23 *24* *25*",
				"26 27 28 29 30 *31* _1_",
			},
			expectedWeekNumbers: []int{48, 49, 50, 51, 52},
		},
	}

	for _, tt := range tests {
		in := &input{
			request: Request{
				FirstDay:         tt.firstDay,
				ShowInactiveDays: tt.showInactiveDays,
			},
		}
		days, weekNumbers := getDays(in, tt.year, tt.month)
		if diff := deep.Equal(days.Strings(in.request.FirstDay), tt.expectedDays); diff != nil {
			t.Errorf("%s: getDays(..., %v, %v): days:\n%s", tt.name, tt.year, tt.month, strings.Join(diff, "\n"))
		}

		if diff := deep.Equal(weekNumbers, tt.expectedWeekNumbers); diff != nil {
			t.Errorf("%s: getDays(..., %v, %v): weekNumbers:\n%v", tt.name, tt.year, tt.month, strings.Join(diff, "\n"))
		}
	}
}
