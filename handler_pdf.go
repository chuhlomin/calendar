package main

import (
	"encoding/hex"
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
		return nil, fmt.Errorf("error decoding request: %w", err)
	}

	if !localizer.HasLanguage(in.request.Language) {
		return nil, fmt.Errorf("language not supported")
	}

	in.textColor, err = decodeColorHex(in.request.TextColor)
	if err != nil {
		return nil, fmt.Errorf("error parsing text color: %w", err)
	}

	in.weekendColor, err = decodeColorHex(in.request.WeekendColor)
	if err != nil {
		return nil, fmt.Errorf("error parsing weekend color: %w", err)
	}

	in.weeknumbersColor, err = decodeColorHex(in.request.WeeknumbersColor)
	if err != nil {
		return nil, fmt.Errorf("error parsing weeknumbers color %q", in.request.WeeknumbersColor)
	}

	in.monthColor, err = decodeColorHex(in.request.MonthColor)
	if err != nil {
		return nil, fmt.Errorf("error parsing month color %q: %w", in.request.MonthColor, err)
	}

	in.weekdaysColor, err = decodeColorHex(in.request.WeekdaysColor)
	if err != nil {
		return nil, fmt.Errorf("error parsing weekdays color %q: %w", in.request.WeekdaysColor, err)
	}

	in.inactiveColor, err = decodeColorHex(in.request.InactiveColor)
	if err != nil {
		return nil, fmt.Errorf("error parsing inactive color %q: %w", in.request.InactiveColor, err)
	}

	if in.request.Month < -1 || in.request.Month > 11 {
		return nil, fmt.Errorf("invalid month %d", in.request.Month)
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
				return fmt.Errorf("error creating PDF month: %w", err)
			}
		}
	} else {
		err := createPDFMonth(writer, pdf, in, time.Month(in.request.Month+1))
		if err != nil {
			return fmt.Errorf("error creating PDF month: %w", err)
		}
	}

	err := pdf.Output(writer)
	if err != nil {
		return fmt.Errorf("error writing PDF: %w", err)
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
	Inactive bool
	Weekend  bool
}

type days []dayInfo

func getDays(in *input, year int32, month time.Month) (days days, weekNumbers []int) {
	firstDay := time.Weekday(in.request.FirstDay)
	lastDay := firstDay - 1
	if lastDay < 0 {
		lastDay = 6
	}

	t := date(int(year), int(month), 1) // 2020 October 1
	end := t.AddDate(0, 1, 0)           // 2020 November 1
	if end.Weekday() == firstDay && in.request.ShowInactiveDays {
		end = end.AddDate(0, 0, 7)
	}

	if t.Weekday() != firstDay {
		if t.Weekday() == lastDay {
			t = t.AddDate(0, 0, -6)
		} else {
			t = t.AddDate(0, 0, -int(t.Weekday()-firstDay))
		}
	}

	var row int32
	var column int32

	// day shift to get week number
	shift := t
	shift = shift.Add(3 * 24 * time.Hour)
	_, weekNumber := shift.ISOWeek()

	for t.Unix() < end.Unix() {
		inactive := t.Month() != month

		if t.Weekday() == firstDay {
			weekNumbers = append(weekNumbers, weekNumber)
			weekNumber++
			if weekNumber > 52 {
				weekNumber = 1
			}
		}

		// skip inactive days in the first week
		if row == 0 && inactive && !in.request.ShowInactiveDays {
			t = t.AddDate(0, 0, 1)
			column++
			continue
		}

		days = append(days, dayInfo{
			Date:     t,
			X:        column*in.request.DaysXStep + in.request.DaysX - in.request.DaysXStep,
			Y:        row*in.request.DaysYStep + in.request.DaysY - in.request.DaysYStep,
			Inactive: inactive,
			Weekend:  t.Weekday() == time.Saturday || t.Weekday() == time.Sunday,
		})

		t = t.AddDate(0, 0, 1)

		if column == 6 {
			row++
			column = 0
		} else {
			column++
		}
	}

	// finish last row
	for t.Weekday() != firstDay && in.request.ShowInactiveDays {
		days = append(days, dayInfo{
			Date:     t,
			X:        column*in.request.DaysXStep + in.request.DaysX - in.request.DaysXStep,
			Y:        row*in.request.DaysYStep + in.request.DaysY - in.request.DaysYStep,
			Inactive: t.Month() != month,
			Weekend:  t.Weekday() == time.Saturday || t.Weekday() == time.Sunday,
		})
		t = t.AddDate(0, 0, 1)
		column++
	}

	return
}

func decodeColorHex(colorStr string) (color.Color, error) {
	b, err := hex.DecodeString(normalizeColor(colorStr))
	if err != nil {
		return nil, err
	}

	return color.RGBA{b[0], b[1], b[2], b[3]}, nil
}

func normalizeColor(colorStr string) string {
	colorStr = strings.TrimPrefix(colorStr, "#")

	if len(colorStr) == 3 {
		colorStr = fmt.Sprintf("%s%s%s%s%s%s", colorStr[0:1], colorStr[0:1], colorStr[1:2], colorStr[1:2], colorStr[2:3], colorStr[2:3])
	}

	if len(colorStr) == 6 {
		colorStr += "FF"
	}

	return colorStr
}
