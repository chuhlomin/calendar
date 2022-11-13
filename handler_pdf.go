package main

import (
	"encoding/json"
	"fmt"
	"image/color"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/jung-kurt/gofpdf"
	"github.com/lucasb-eyer/go-colorful"
	"github.com/pkg/errors"
)

type pdfRequest struct {
	Size     string `json:"size"`
	FirstDay int    `json:"firstDay"` // 0 - Sunday, 1 - Monday, ...
	Year     int    `json:"year"`
	Month    int    `json:"month"`

	// Days
	DaysFontSize     int    `json:"daysFontSize"`
	DaysFontFamily   string `json:"daysFontFamily"`
	TextColor        string `json:"textColor"`
	WeekendColor     string `json:"weekendColor"`
	DaysX            int    `json:"daysX"`
	DaysY            int    `json:"daysY"`
	DaysXStep        int    `json:"daysXStep"`
	DaysYStep        int    `json:"daysYStep"`
	ShowInactiveDays bool   `json:"showInactiveDays"`
	InactiveColor    string `json:"inactiveColor"`

	// Month
	ShowMonth       bool   `json:"showMonth"`
	MonthFontFamily string `json:"monthFontFamily"`
	MonthFontSize   int    `json:"monthFontSize"`
	MonthColor      string `json:"monthColor"`
	MonthY          int    `json:"monthY"`

	// Weekdays
	ShowWeekdays       bool   `json:"showWeekdays"`
	WeekdaysFontFamily string `json:"weekdaysFontFamily"`
	WeekdaysFontSize   int    `json:"weekdaysFontSize"`
	WeekdaysColor      string `json:"weekdaysColor"`
	WeekdaysX          int    `json:"weekdaysX"`
	WeekdaysY          int    `json:"weekdaysY"`

	// WeekNumbers
	ShowWeekNumbers       bool   `json:"showWeekNumbers"`
	WeeknumbersFontFamily string `json:"weeknumbersFontFamily"`
	WeeknumbersFontSize   int    `json:"weeknumbersFontSize"`
	WeeknumbersColor      string `json:"weeknumbersColor"`
	WeeknumbersX          int    `json:"weeknumbersX"`
	WeeknumbersY          int    `json:"weeknumbersY"`

	textColor        color.Color
	weekendColor     color.Color
	weeknumbersColor color.Color
	monthColor       color.Color
	weekdaysColor    color.Color
	inactiveColor    color.Color
}

func handlerPDF(w http.ResponseWriter, r *http.Request) {
	req, err := parsePDFRequest(r)
	if err != nil {
		log.Printf("error parsing PDF request: %s", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", "attachment; filename=calendar.pdf")

	err = createPDF(w, req)
	if err != nil {
		log.Printf("error creating PDF: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Sorry, something went wrong"))
		return
	}
}

func parsePDFRequest(r *http.Request) (*pdfRequest, error) {
	var req pdfRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return nil, errors.Wrap(err, "error decoding json")
	}

	req.textColor, err = colorful.Hex(req.TextColor)
	if err != nil {
		return nil, errors.Wrapf(err, "error parsing text color %q", req.TextColor)
	}

	req.weekendColor, err = colorful.Hex(req.WeekendColor)
	if err != nil {
		return nil, errors.Wrapf(err, "error parsing weekend color %q", req.WeekendColor)
	}

	req.weeknumbersColor, err = colorful.Hex(req.WeeknumbersColor)
	if err != nil {
		return nil, errors.Wrapf(err, "error parsing weeknumbers color %q", req.WeeknumbersColor)
	}

	req.monthColor, err = colorful.Hex(req.MonthColor)
	if err != nil {
		return nil, errors.Wrapf(err, "error parsing month color %q", req.MonthColor)
	}

	req.weekdaysColor, err = colorful.Hex(req.WeekdaysColor)
	if err != nil {
		return nil, errors.Wrapf(err, "error parsing weekdays color %q", req.WeekdaysColor)
	}

	req.inactiveColor, err = colorful.Hex(req.InactiveColor)
	if err != nil {
		return nil, errors.Wrapf(err, "error parsing inactive color %q", req.InactiveColor)
	}

	if req.Month < -1 || req.Month > 11 {
		return nil, errors.Errorf("invalid month %d", req.Month)
	}

	return &req, nil
}

func createPDF(writer io.Writer, req *pdfRequest) error {
	pdf := gofpdf.New("P", "mm", req.Size, "")
	pdf.SetAutoPageBreak(false, 0)
	pdf.SetFillColor(255, 255, 255)

	pdf.AddUTF8Font("daysFontFamily", "", fmt.Sprint("fonts/", req.DaysFontFamily, ".ttf"))
	pdf.AddUTF8Font("monthFontFamily", "", fmt.Sprint("fonts/", req.MonthFontFamily, ".ttf"))
	pdf.AddUTF8Font("weekdaysFontFamily", "", fmt.Sprint("fonts/", req.WeekdaysFontFamily, ".ttf"))
	pdf.AddUTF8Font("weeknumbersFontFamily", "", fmt.Sprint("fonts/", req.WeeknumbersFontFamily, ".ttf"))

	if req.Month == -1 {
		for month := 0; month < 12; month++ {
			err := createPDFMonth(writer, pdf, req, time.Month(month+1))
			if err != nil {
				return errors.Wrapf(err, "error creating PDF for month %d", month)
			}
		}
	} else {
		err := createPDFMonth(writer, pdf, req, time.Month(req.Month+1))
		if err != nil {
			return errors.Wrapf(err, "error creating PDF for month %d", req.Month)
		}
	}

	err := pdf.Output(writer)
	if err != nil {
		return errors.Wrap(err, "error writing pdf")
	}

	return nil
}

func createPDFMonth(writer io.Writer, pdf *gofpdf.Fpdf, req *pdfRequest, month time.Month) error {
	pdf.AddPage()

	days, weekNumbers := getDays(req, req.Year, month)

	setTextColor(pdf, req.textColor)

	if req.ShowMonth {
		drawMonth(pdf, req, req.Year, month)
	}

	if req.ShowWeekdays {
		drawWeekdays(pdf, req)
	}

	if req.ShowWeekNumbers {
		drawWeekNumbers(pdf, req, weekNumbers)
	}

	drawDays(pdf, req, days)

	return nil
}

func drawWeekdays(pdf *gofpdf.Fpdf, req *pdfRequest) {
	setTextColor(pdf, req.weekdaysColor)
	pdf.SetFont("weekdaysFontFamily", "", float64(req.WeekdaysFontSize))

	for day := 0; day < 7; day++ {
		x := day*req.DaysXStep + req.WeekdaysX - req.DaysXStep

		day := time.Weekday((day + 1) % 7).String()
		day = day[:3]

		pdf.MoveTo(float64(x), 0)
		pdf.CellFormat(float64(req.DaysXStep), float64(req.WeekdaysY), day, "0", 0, "RA", false, 0, "")
	}
}

func drawWeekNumbers(pdf *gofpdf.Fpdf, req *pdfRequest, weekNumbers []int) {
	setTextColor(pdf, req.weeknumbersColor)
	pdf.SetFont("weeknumbersFontFamily", "", float64(req.WeeknumbersFontSize))

	line := 0
	for _, weekNumber := range weekNumbers {
		y := line*req.DaysYStep + req.WeeknumbersY - req.DaysYStep

		pdf.MoveTo(0, float64(y))
		pdf.CellFormat(float64(req.WeeknumbersX), float64(req.DaysYStep), strconv.Itoa(weekNumber), "0", 0, "RA", false, 0, "")

		line++
	}
}

func drawDays(pdf *gofpdf.Fpdf, req *pdfRequest, days []dayInfo) {
	pdf.SetFont("daysFontFamily", "", float64(req.DaysFontSize))

	var color color.Color

	for _, dayInfo := range days {
		color = req.textColor
		if dayInfo.Inactive {
			color = req.inactiveColor
		} else {
			if dayInfo.Weekend {
				color = req.weekendColor
			}
		}

		if dayInfo.Inactive && !req.ShowInactiveDays {
			continue
		}

		setTextColor(pdf, color)
		pdf.MoveTo(float64(dayInfo.X), float64(dayInfo.Y))
		pdf.CellFormat(
			float64(req.DaysXStep),
			float64(req.DaysYStep),
			fmt.Sprintf("%d", dayInfo.Date.Day()),
			"0", 0, "RA", false, 0, "",
		)
	}
}

func drawMonth(pdf *gofpdf.Fpdf, req *pdfRequest, year int, month time.Month) {
	setTextColor(pdf, req.monthColor)
	pdf.SetFont("monthFontFamily", "", float64(req.MonthFontSize))

	pdf.MoveTo(0, float64(req.MonthY))
	pdf.CellFormat(
		// w, h, fmt.Sprintf("%s %d", i18n(month.String()), year),
		200, 0, fmt.Sprintf("%s %d", month.String(), year),
		"0", 0, "CA", false, 0, "",
	)
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
	X        int
	Y        int
	Weekend  bool
	Inactive bool
}

func getDays(req *pdfRequest, year int, month time.Month) (days []dayInfo, weekNumbers []int) {
	start := date(year, int(month), 1) // 2020 October 1
	end := start.AddDate(0, 1, 0)      // 2020 November 1

	firstDay := time.Weekday(req.FirstDay)
	weekday := start.Weekday() // Thursday
	daysSinceWeekStartToBeginningOfMonth := int(weekday) - int(firstDay)
	if daysSinceWeekStartToBeginningOfMonth == -1 {
		// Sun (0) - Mon (1) = -1
		daysSinceWeekStartToBeginningOfMonth = 6
	}

	sheetStart := start.AddDate(0, 0, -daysSinceWeekStartToBeginningOfMonth)
	t := sheetStart

	row := 0
	column := 0
	_, weekNumber := t.ISOWeek()

	for t.Unix() < end.Unix() {
		days = append(days, dayInfo{
			Date:     t,
			X:        column*req.DaysXStep + req.DaysX - req.DaysXStep,
			Y:        row*req.DaysYStep + req.DaysY - req.DaysYStep,
			Weekend:  t.Weekday() == time.Saturday || t.Weekday() == time.Sunday,
			Inactive: t.Month() != month,
		})

		if int(t.Weekday()) == req.FirstDay {
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

	if !req.ShowInactiveDays {
		return
	}

	// finish
	for column < 7 {
		days = append(days, dayInfo{
			Date:     t,
			X:        column*req.DaysXStep + req.DaysX - req.DaysXStep,
			Y:        row*req.DaysYStep + req.DaysY - req.DaysYStep,
			Weekend:  t.Weekday() == time.Saturday || t.Weekday() == time.Sunday,
			Inactive: t.Month() != month,
		})
		column++
		t = t.AddDate(0, 0, 1)
	}

	return
}
