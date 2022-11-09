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
	FontSizeDays     int    `json:"fontSizeDays"`
	TextColor        string `json:"textColor"`
	WeekendColor     string `json:"weekendColor"`
	DaysXShift       int    `json:"daysXShift"`
	DaysXStep        int    `json:"daysXStep"`
	DaysYShift       int    `json:"daysYShift"`
	DaysYStep        int    `json:"daysYStep"`
	ShowInactiveDays bool   `json:"showInactiveDays"`
	InactiveColor    string `json:"inactiveColor"`

	// Month
	ShowMonth     bool   `json:"showMonth"`
	FontSizeMonth int    `json:"fontSizeMonth"`
	MonthColor    string `json:"monthColor"`

	// Weekdays
	ShowWeekdays     bool   `json:"showWeekdays"`
	FontSizeWeekdays int    `json:"fontSizeWeekdays"`
	WeekdaysColor    string `json:"weekdaysColor"`
	WeekdaysXShift   int    `json:"weekdaysXShift"`
	WeekdaysXStep    int    `json:"weekdaysXStep"`
	WeekdaysYShift   int    `json:"weekdaysYShift"`

	// WeekNumbers
	ShowWeekNumbers     bool   `json:"showWeekNumbers"`
	FontSizeWeekNumbers int    `json:"fontSizeWeekNumbers"`
	WeeknumbersColor    string `json:"weeknumbersColor"`
	WeeknumbersXShift   int    `json:"weeknumbersXShift"`
	WeeknumbersYShift   int    `json:"weeknumbersYShift"`
	WeeknumbersYStep    int    `json:"weeknumbersYStep"`

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

	if req.Month < 0 || req.Month > 11 {
		return nil, errors.Errorf("invalid month %d", req.Month)
	}

	return &req, nil
}

func createPDF(writer io.Writer, req *pdfRequest) error {
	pdf := gofpdf.New("P", "mm", req.Size, "")
	pdf.SetAutoPageBreak(false, 0)
	pdf.AddPage()

	err := drawPage(pdf, req)
	if err != nil {
		return errors.Wrap(err, "error drawing month")
	}

	err = pdf.Output(writer)
	if err != nil {
		return errors.Wrap(err, "error writing pdf")
	}

	return nil
}

func drawPage(pdf *gofpdf.Fpdf, req *pdfRequest) error {
	pdf.SetFillColor(255, 255, 255)

	pdf.AddUTF8Font("numbers", "", "static/iosevka-regular.ttf")
	pdf.AddUTF8Font("month", "", "static/iosevka-aile-regular.ttf")

	year := req.Year
	month := time.Month(req.Month + 1)

	days, weekNumbers := getDays(req, year, month)

	setTextColor(pdf, req.textColor)

	if req.ShowMonth {
		drawMonth(pdf, req, year, month)
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
	pdf.SetFont("month", "", float64(req.FontSizeWeekdays*3))

	for day := 0; day < 7; day++ {
		x := day*req.WeekdaysXStep + req.WeekdaysYShift - req.WeekdaysXStep

		day := time.Weekday((day + 1) % 7).String()
		day = day[:3]

		pdf.MoveTo(float64(x), 0)
		pdf.CellFormat(float64(req.WeekdaysXShift), float64(req.WeekdaysYShift), day, "0", 0, "RB", false, 0, "")
	}
}

func drawWeekNumbers(pdf *gofpdf.Fpdf, req *pdfRequest, weekNumbers []int) {
	setTextColor(pdf, req.weeknumbersColor)
	pdf.SetFont("numbers", "", float64(req.FontSizeWeekNumbers*3))

	line := 0
	for _, weekNumber := range weekNumbers {
		y := line*req.WeeknumbersYStep + req.WeeknumbersYShift - req.WeeknumbersYStep

		pdf.MoveTo(0, float64(y))
		pdf.CellFormat(float64(req.WeeknumbersXShift), float64(req.WeeknumbersYStep), strconv.Itoa(weekNumber), "0", 0, "RB", false, 0, "")

		line++
	}
}

func drawDays(pdf *gofpdf.Fpdf, req *pdfRequest, days []dayInfo) {
	pdf.SetFont("numbers", "", float64(req.FontSizeDays*3))

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
			"0", 0, "RB", false, 0, "",
		)
	}
}

func drawMonth(pdf *gofpdf.Fpdf, req *pdfRequest, year int, month time.Month) {
	var x float64 = 5
	var y float64 = 250
	var w float64 = 200
	var h float64 = 20

	setTextColor(pdf, req.monthColor)
	pdf.SetFont("month", "", float64(req.FontSizeMonth*3))

	pdf.MoveTo(x, y)
	pdf.CellFormat(
		// w, h, fmt.Sprintf("%s %d", i18n(month.String()), year),
		w, h, fmt.Sprintf("%s %d", month.String(), year),
		"0", 0, "C", false, 0, "",
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
			X:        column*req.DaysXStep + req.DaysXShift - req.DaysXStep,
			Y:        row*req.DaysYStep + req.DaysYShift - req.DaysYStep,
			Weekend:  t.Weekday() == time.Saturday || t.Weekday() == time.Sunday,
			Inactive: t.Month() != month,
		})

		weekNumbers = append(weekNumbers, weekNumber)
		weekNumber++
		if weekNumber > 52 {
			weekNumber = 1
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
			X:        column*req.DaysXStep + req.DaysXShift - req.DaysXStep,
			Y:        row*req.DaysYStep + req.DaysYShift - req.DaysYStep,
			Weekend:  t.Weekday() == time.Saturday || t.Weekday() == time.Sunday,
			Inactive: t.Month() != month,
		})
		column++
		t = t.AddDate(0, 0, 1)
	}

	return
}
