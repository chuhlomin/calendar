package dates

import "time"

func GetTime(year int, month time.Month) time.Time {
	return time.Date(year, month, 1, 0, 0, 0, 0, time.Local)
}

func Date(year, month, day int) time.Time {
	return time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
}

func GetCalendar(year int, month time.Month) [][7]int {
	var calendar [][7]int

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

	// log.Printf("Sheet start %v", sheetStart.Format("2006 Jan 2"))
	// log.Printf("End %v", end.Format("2006 Jan 2"))

	row := [7]int{}
	column := 0
	line := 0
	for t := sheetStart; t.Unix() < end.Unix(); {

		row[column] = t.Day()

		// log.Printf("Current day %v", t.Format("2006 Jan 2"))

		column++

		if column > 6 {
			column = 0
			line++
			calendar = append(calendar, row)
			row = [7]int{}
		}

		t = t.AddDate(0, 0, 1)
	}

	// finish
	nextMonthDay := 1
	for column < 7 {
		row[column] = nextMonthDay
		column++
		nextMonthDay++
	}
	calendar = append(calendar, row)

	return calendar
}
