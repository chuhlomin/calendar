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

	"github.com/chuhlomin/calendar/dates"
)

type pdfRequest struct {
	Size         string `json:"size"`
	Orientation  string `json:"orientation"`
	FirstDay     string `json:"firstDay"`
	TextColor    string `json:"textColor"`
	WeekendColor string `json:"weekendColor"`
	Year         int    `json:"year"`
	Month        int    `json:"month"`

	textColor    color.Color
	weekendColor color.Color
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

	// w, h := pdf.GetPageSize()

	year := req.Year
	month := time.Month(req.Month + 1)

	calendar := dates.GetCalendar(year, month)

	setTextColor(pdf, req.textColor)

	drawMonth(pdf, year, month)
	drawWeekdays(pdf, req)
	drawWeekNumbers(pdf, year, month, len(calendar))
	drawDays(pdf, req, calendar)

	return nil
}

func drawWeekdays(pdf *gofpdf.Fpdf, req *pdfRequest) {
	var w float64 = 20
	var h float64 = 20
	var left float64 = 20
	var top float64 = 21
	var marginLeft float64 = 5

	pdf.SetFont("month", "", 12)
	pdf.SetTextColor(200, 200, 200)

	for day := 0; day < 7; day++ {
		x := left + float64(day)*(w+marginLeft)

		day := time.Weekday((day + 1) % 7).String()
		day = day[:3]

		pdf.MoveTo(x, top)
		pdf.CellFormat(w, h, day, "0", 0, "RT", false, 0, "")
	}
}

func drawWeekNumbers(pdf *gofpdf.Fpdf, year int, month time.Month, lines int) {
	var w float64 = 10
	var h float64 = 20
	var left float64 = 9
	var top float64 = 35
	var marginTop float64 = 15

	pdf.SetFont("numbers", "", 12)
	pdf.SetTextColor(200, 200, 200)

	start := dates.Date(year, int(month), 1)
	_, week := start.ISOWeek()

	for line := 0; line < lines; line++ {
		y := top + float64(line)*(h+marginTop)

		pdf.MoveTo(left, y)
		pdf.CellFormat(w, h, strconv.Itoa(week), "0", 0, "RT", false, 0, "")

		week++

		if week > 52 {
			week = 1
		}
	}
}

func drawDays(pdf *gofpdf.Fpdf, req *pdfRequest, calendar [][7]int) {
	year := 2022
	month := time.November
	holidays := dates.GetHolidays(year, month)

	var w float64 = 20
	var h float64 = 20
	var left float64 = 20
	var top float64 = 35
	var marginLeft float64 = 5
	var marginTop float64 = 15

	colorGray, _ := colorful.Hex("#c8c8c8")

	prevDay := 0
	useGrayColor := true
	pdf.SetFont("numbers", "", 36)

	var color color.Color

	for line, row := range calendar {
		for column, day := range row {
			x := left + float64(column)*(w+marginLeft)
			y := top + float64(line)*(h+marginTop)

			if prevDay > day || (prevDay == 0 && day == 1) {
				useGrayColor = !useGrayColor
			}
			prevDay = day

			color = req.textColor
			if useGrayColor {
				color = colorGray
			} else {
				if column > 4 {
					color = req.weekendColor
				}
			}

			if holiday := holidays[day]; !useGrayColor && false {
				// disable holiday rendering till this feature is on frontend
				if holiday != "" {
					color = req.weekendColor

					holiday = strings.ReplaceAll(holiday, " Day", "")
					holiday = strings.ReplaceAll(holiday, "'s", "")

					align := "R"
					if day > 9 {
						align = "C"
					}

					setTextColor(pdf, color)
					pdf.SetFont("numbers", "", 6)
					pdf.MoveTo(x, y+20)
					pdf.CellFormat(
						w, 5, holiday,
						"0", 0, align, false, 0, "",
					)
					// pdf.SetFont("numbers", "", 36)
				}
			}

			setTextColor(pdf, color)
			pdf.MoveTo(x, y)
			pdf.CellFormat(
				w, h, fmt.Sprintf("%d", day),
				"0", 0, "RT", false, 0, "",
			)
		}
	}
}

func drawMonth(pdf *gofpdf.Fpdf, year int, month time.Month) {
	var x float64 = 5
	var y float64 = 250
	var w float64 = 200
	var h float64 = 20

	pdf.SetFont("month", "", 36)
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
