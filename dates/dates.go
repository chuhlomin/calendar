package dates

import (
	"time"
)

func GetTime(year int, month time.Month) time.Time {
	return time.Date(year, month, 1, 0, 0, 0, 0, time.Local)
}

func Date(year, month, day int) time.Time {
	return time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
}

type DayInfo struct {
	Date     time.Time
	X        int
	Y        int
	Weekend  bool
	Inactive bool
}

func GetCalendar(
	year int,
	month time.Month,
	daysXStep,
	daysXShift int,
	daysYStep,
	daysYShift int,
) [][7]DayInfo {
	var calendar [][7]DayInfo

	start := Date(year, int(month), 1) // 2020 October 1
	end := start.AddDate(0, 1, 0)      // 2020 November 1

	weekStartsOn := time.Monday
	weekday := start.Weekday() // Thursday
	daysSinceWeekStartToBeginningOfMonth := int(weekday) - int(weekStartsOn)
	if daysSinceWeekStartToBeginningOfMonth == -1 {
		// Sun (0) - Mon (1) = -1
		daysSinceWeekStartToBeginningOfMonth = 6
	}

	sheetStart := start.AddDate(0, 0, -daysSinceWeekStartToBeginningOfMonth)
	t := sheetStart

	rowData := [7]DayInfo{}
	row := 0
	column := 0

	for t.Unix() < end.Unix() {
		rowData[column] = DayInfo{
			Date:     t,
			X:        column*daysXStep + daysXShift - daysXStep,
			Y:        row*daysYStep + daysYShift - daysYStep,
			Weekend:  t.Weekday() == time.Saturday || t.Weekday() == time.Sunday,
			Inactive: t.Month() != month,
		}

		if column == 6 {
			row++
			column = 0
			calendar = append(calendar, rowData)
			rowData = [7]DayInfo{}
		} else {
			column++
		}

		t = t.AddDate(0, 0, 1)
	}

	// finish
	for column < 7 {
		rowData[column] = DayInfo{
			Date:     t,
			X:        column*daysXStep + daysXShift - daysXStep,
			Y:        row*daysYStep + daysYShift - daysYStep,
			Weekend:  t.Weekday() == time.Saturday || t.Weekday() == time.Sunday,
			Inactive: t.Month() != month,
		}
		column++
		t = t.AddDate(0, 0, 1)
	}
	calendar = append(calendar, rowData)

	return calendar
}
