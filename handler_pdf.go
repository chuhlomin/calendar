package main

import (
	"encoding/json"
	"fmt"
	"image/color"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/jung-kurt/gofpdf"
	"github.com/lucasb-eyer/go-colorful"
	"github.com/pkg/errors"
)

type input struct {
	request          Request
	textColor        color.Color
	weekendColor     color.Color
	weeknumbersColor color.Color
	monthColor       color.Color
	weekdaysColor    color.Color
	inactiveColor    color.Color
}

func handlerPDF(w http.ResponseWriter, r *http.Request) {
	in, err := parsePDFRequest(r)
	if err != nil {
		log.Printf("error parsing PDF request: %s", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", "attachment; filename=calendar.pdf")

	err = createPDF(w, in)
	if err != nil {
		log.Printf("error creating PDF: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Sorry, something went wrong"))
		return
	}
}

func parsePDFRequest(r *http.Request) (*input, error) {
	in := input{}
	err := json.NewDecoder(r.Body).Decode(&in.request)
	if err != nil {
		return nil, errors.Wrap(err, "error decoding json")
	}

	if !localizer.HasLanguage(in.request.Language) {
		return nil, errors.New("language not supported")
	}

	in.textColor, err = colorful.Hex(in.request.TextColor)
	if err != nil {
		return nil, errors.Wrapf(err, "error parsing text color %q", in.request.TextColor)
	}

	in.weekendColor, err = colorful.Hex(in.request.WeekendColor)
	if err != nil {
		return nil, errors.Wrapf(err, "error parsing weekend color %q", in.request.WeekendColor)
	}

	in.weeknumbersColor, err = colorful.Hex(in.request.WeeknumbersColor)
	if err != nil {
		return nil, errors.Wrapf(err, "error parsing weeknumbers color %q", in.request.WeeknumbersColor)
	}

	in.monthColor, err = colorful.Hex(in.request.MonthColor)
	if err != nil {
		return nil, errors.Wrapf(err, "error parsing month color %q", in.request.MonthColor)
	}

	in.weekdaysColor, err = colorful.Hex(in.request.WeekdaysColor)
	if err != nil {
		return nil, errors.Wrapf(err, "error parsing weekdays color %q", in.request.WeekdaysColor)
	}

	in.inactiveColor, err = colorful.Hex(in.request.InactiveColor)
	if err != nil {
		return nil, errors.Wrapf(err, "error parsing inactive color %q", in.request.InactiveColor)
	}

	if in.request.Month < -1 || in.request.Month > 11 {
		return nil, errors.Errorf("invalid month %d", in.request.Month)
	}

	return &in, nil
}

func createPDF(writer io.Writer, in *input) error {
	pdf := gofpdf.New("P", "mm", in.request.Size, "")
	pdf.SetAutoPageBreak(false, 0)
	pdf.SetFillColor(255, 255, 255)

	pdf.AddUTF8Font("daysFontFamily", "", fmt.Sprint("fonts/", in.request.DaysFontFamily, ".ttf"))
	pdf.AddUTF8Font("monthFontFamily", "", fmt.Sprint("fonts/", in.request.MonthFontFamily, ".ttf"))
	pdf.AddUTF8Font("weekdaysFontFamily", "", fmt.Sprint("fonts/", in.request.WeekdaysFontFamily, ".ttf"))
	pdf.AddUTF8Font("weeknumbersFontFamily", "", fmt.Sprint("fonts/", in.request.WeeknumbersFontFamily, ".ttf"))

	if in.request.Month == -1 {
		for month := 0; month < 12; month++ {
			err := createPDFMonth(writer, pdf, in, time.Month(month+1))
			if err != nil {
				return errors.Wrapf(err, "error creating PDF for month %d", month)
			}
		}
	} else {
		err := createPDFMonth(writer, pdf, in, time.Month(in.request.Month+1))
		if err != nil {
			return errors.Wrapf(err, "error creating PDF for month %d", in.request.Month)
		}
	}

	err := pdf.Output(writer)
	if err != nil {
		return errors.Wrap(err, "error writing pdf")
	}

	return nil
}

func createPDFMonth(writer io.Writer, pdf *gofpdf.Fpdf, in *input, month time.Month) error {
	pdf.AddPage()

	days, weekNumbers := getDays(in, in.request.Year, month)

	setTextColor(pdf, in.textColor)

	if in.request.ShowMonth {
		drawMonth(pdf, in, in.request.Year, month)
	}

	if in.request.ShowWeekdays {
		drawWeekdays(pdf, in)
	}

	if in.request.ShowWeekNumbers {
		drawWeekNumbers(pdf, in, weekNumbers)
	}

	drawDays(pdf, in, days)

	return nil
}

func drawWeekdays(pdf *gofpdf.Fpdf, in *input) {
	setTextColor(pdf, in.weekdaysColor)
	pdf.SetFont("weekdaysFontFamily", "", float64(in.request.WeekdaysFontSize))

	var day int32
	for day = 0; day < 7; day++ {
		x := day*in.request.DaysXStep + in.request.WeekdaysX - in.request.DaysXStep

		day := time.Weekday((day + 1) % 7).String()
		day = localizer.I18n(in.request.Language, "weekday_short_"+strings.ToLower(day))

		pdf.MoveTo(float64(x), 0)
		pdf.CellFormat(float64(in.request.DaysXStep), float64(in.request.WeekdaysY), day, "0", 0, "RA", false, 0, "")
	}
}

func drawWeekNumbers(pdf *gofpdf.Fpdf, in *input, weekNumbers []int) {
	setTextColor(pdf, in.weeknumbersColor)
	pdf.SetFont("weeknumbersFontFamily", "", float64(in.request.WeeknumbersFontSize))

	var line int32
	for _, weekNumber := range weekNumbers {
		y := line*in.request.DaysYStep + in.request.WeeknumbersY - in.request.DaysYStep

		pdf.MoveTo(0, float64(y))
		pdf.CellFormat(float64(in.request.WeeknumbersX), float64(in.request.DaysYStep), strconv.Itoa(weekNumber), "0", 0, "RA", false, 0, "")

		line++
	}
}

func drawDays(pdf *gofpdf.Fpdf, in *input, days []dayInfo) {
	pdf.SetFont("daysFontFamily", "", float64(in.request.DaysFontSize))

	var color color.Color

	for _, dayInfo := range days {
		color = in.textColor
		if dayInfo.Inactive {
			color = in.inactiveColor
		} else {
			if dayInfo.Weekend {
				color = in.weekendColor
			}
		}

		if dayInfo.Inactive && !in.request.ShowInactiveDays {
			continue
		}

		setTextColor(pdf, color)
		pdf.MoveTo(float64(dayInfo.X), float64(dayInfo.Y))
		pdf.CellFormat(
			float64(in.request.DaysXStep),
			float64(in.request.DaysYStep),
			fmt.Sprintf("%d", dayInfo.Date.Day()),
			"0", 0, "RA", false, 0, "",
		)
	}
}

func drawMonth(pdf *gofpdf.Fpdf, in *input, year int32, month time.Month) {
	setTextColor(pdf, in.monthColor)
	pdf.SetFont("monthFontFamily", "", float64(in.request.MonthFontSize))

	line := time.Date(int(year), month, 1, 0, 0, 0, 0, time.UTC).Format(in.request.MonthFormat)
	log.Printf("line: %s, format: %s", line, in.request.MonthFormat)
	replacer := strings.NewReplacer(
		"January", localizer.I18n(in.request.Language, "month_january"),
		"February", localizer.I18n(in.request.Language, "month_february"),
		"March", localizer.I18n(in.request.Language, "month_march"),
		"April", localizer.I18n(in.request.Language, "month_april"),
		"May", localizer.I18n(in.request.Language, "month_may"),
		"June", localizer.I18n(in.request.Language, "month_june"),
		"July", localizer.I18n(in.request.Language, "month_july"),
		"August", localizer.I18n(in.request.Language, "month_august"),
		"September", localizer.I18n(in.request.Language, "month_september"),
		"October", localizer.I18n(in.request.Language, "month_october"),
		"November", localizer.I18n(in.request.Language, "month_november"),
		"December", localizer.I18n(in.request.Language, "month_december"),
	)
	line = replacer.Replace(line)

	pdf.MoveTo(0, float64(in.request.MonthY))
	pdf.CellFormat(200, 0, line, "0", 0, "CA", false, 0, "")
}

func setTextColor(pdf *gofpdf.Fpdf, color color.Color) {
	r, g, b, a := color.RGBA()
	pdf.SetTextColor(int(r/0x0101), int(g/0x0101), int(b/0x0101))
	pdf.SetAlpha(float64(a)/0xffff, "Normal")
}

func convertColor(hex string) (int, int, int, float64) {
	var r, g, b, a int64

	if len(hex) == 7 {
		r, _ = strconv.ParseInt(hex[1:3], 16, 0)
		g, _ = strconv.ParseInt(hex[3:5], 16, 0)
		b, _ = strconv.ParseInt(hex[5:7], 16, 0)
		return int(r), int(g), int(b), 1.0
	}

	if len(hex) == 9 {
		r, _ = strconv.ParseInt(hex[1:3], 16, 0)
		g, _ = strconv.ParseInt(hex[3:5], 16, 0)
		b, _ = strconv.ParseInt(hex[5:7], 16, 0)
		a, _ = strconv.ParseInt(hex[7:9], 16, 0)
		return int(r), int(g), int(b), float64(a) / 255.0
	}

	return 0, 0, 0, 0
}

func date(year, month, day int) time.Time {
	return time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
}

type dayInfo struct {
	Date     time.Time
	X        int32
	Y        int32
	Weekend  bool
	Inactive bool
}

func getDays(in *input, year int32, month time.Month) (days []dayInfo, weekNumbers []int) {
	start := date(int(year), int(month), 1) // 2020 October 1
	end := start.AddDate(0, 1, 0)           // 2020 November 1

	firstDay := time.Weekday(in.request.FirstDay)
	weekday := start.Weekday() // Thursday
	daysSinceWeekStartToBeginningOfMonth := int(weekday) - int(firstDay)
	if daysSinceWeekStartToBeginningOfMonth == -1 {
		// Sun (0) - Mon (1) = -1
		daysSinceWeekStartToBeginningOfMonth = 6
	}

	sheetStart := start.AddDate(0, 0, -daysSinceWeekStartToBeginningOfMonth)
	t := sheetStart

	var row int32
	var column int32
	_, weekNumber := t.ISOWeek()

	for t.Unix() < end.Unix() {
		days = append(days, dayInfo{
			Date:     t,
			X:        column*in.request.DaysXStep + in.request.DaysX - in.request.DaysXStep,
			Y:        row*in.request.DaysYStep + in.request.DaysY - in.request.DaysYStep,
			Weekend:  t.Weekday() == time.Saturday || t.Weekday() == time.Sunday,
			Inactive: t.Month() != month,
		})

		if int32(t.Weekday()) == in.request.FirstDay {
			weekNumbers = append(weekNumbers, weekNumber)
			weekNumber++
			if weekNumber > 52 {
				weekNumber = 1
			}
		}

		if column == 6 {
			row++
			column = 0
		} else {
			column++
		}

		t = t.AddDate(0, 0, 1)
	}

	if !in.request.ShowInactiveDays {
		return
	}

	// finish
	for column < 7 {
		days = append(days, dayInfo{
			Date:     t,
			X:        column*in.request.DaysXStep + in.request.DaysX - in.request.DaysXStep,
			Y:        row*in.request.DaysYStep + in.request.DaysY - in.request.DaysYStep,
			Weekend:  t.Weekday() == time.Saturday || t.Weekday() == time.Sunday,
			Inactive: t.Month() != month,
		})
		column++
		t = t.AddDate(0, 0, 1)
	}

	return
}
