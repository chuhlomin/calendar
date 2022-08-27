package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/jung-kurt/gofpdf"
	"github.com/pkg/errors"
)

type pdfRequest struct {
	Size        string  `json:"size"`
	Orientation string  `json:"orientation"`
	Pattern     pattern `json:"pattern"`
}

type pattern struct {
	Name      string `json:"name"`
	Size      string `json:"size"`
	Height    string `json:"height"`
	Color     string `json:"color"`
	LineWidth string `json:"lineWidth"`
	size      float64
	height    float64
	lineWidth float64
}

var errUnknownPattern = errors.New("unknown pattern")

func handlerPDF(w http.ResponseWriter, r *http.Request) {
	req, err := parsePDFRequest(r)
	if err != nil {
		log.Printf("error parsing PDF request: %s", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", "attachment; filename=grid.pdf")

	err = createPDF(w, req.Size, req.Orientation, req.Pattern)
	if err != nil {
		log.Printf("error creating PDF: %s", err)
		if errors.Cause(err) == errUnknownPattern {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Sorry, don't know how to draw a " + req.Pattern.Name + " pattern."))
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Sorry, something went wrong"))
		}
		return
	}
}

func parsePDFRequest(r *http.Request) (*pdfRequest, error) {
	var req pdfRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return nil, errors.Wrap(err, "error decoding json")
	}

	// convert pattern size from string to float64
	req.Pattern.size, err = strconv.ParseFloat(req.Pattern.Size, 64)
	if err != nil {
		return nil, errors.Wrap(err, "error converting pattern size")
	}
	req.Pattern.height, err = strconv.ParseFloat(req.Pattern.Height, 64)
	if err != nil {
		return nil, errors.Wrap(err, "error converting pattern height")
	}

	req.Pattern.lineWidth, err = strconv.ParseFloat(req.Pattern.LineWidth, 64)
	if err != nil {
		return nil, errors.Wrap(err, "error converting pattern line width")
	}

	return &req, nil
}

func createPDF(writer io.Writer, size, orientation string, pattern pattern) error {
	pdf := gofpdf.New(orientation, "mm", size, "")
	pdf.SetAutoPageBreak(false, 0)
	pdf.AddPage()

	err := drawPattern(pdf, pattern)
	if err != nil {
		return errors.Wrap(err, "error drawing pattern")
	}

	err = pdf.Output(writer)
	if err != nil {
		return errors.Wrap(err, "error writing pdf")
	}

	return nil
}

func drawPattern(pdf gofpdf.Pdf, pattern pattern) error {
	r, g, b, a := convertColor(pattern.Color)
	pdf.SetFillColor(255, 255, 255)
	pdf.SetDrawColor(r, g, b)
	pdf.SetAlpha(a, "Normal")

	w, h := pdf.GetPageSize()
	patternSize := pattern.size
	patternHeight := pattern.height

	pdf.SetLineWidth(pattern.lineWidth / 1000)

	switch pattern.Name {
	case "rect":
		w, h := pdf.GetPageSize()
		for x := 0.0; x < w; x += patternSize {
			for y := 0.0; y < h; y += patternSize {
				pdf.Rect(x, y, patternSize, patternSize, "D")
			}
		}
	case "lines":
		w, h := pdf.GetPageSize()
		for y := 0.0; y < h; y += patternSize {
			pdf.Line(0, y, w, y)
		}
	case "dot":
		pdf.SetFillColor(r, g, b)
		for x := 0.0; x < w; x += patternSize {
			for y := 0.0; y < h; y += patternSize {
				pdf.Circle(x, y, pattern.lineWidth/1000, "F")
			}
		}
	case "diamond":
		for x := 0.0; x < w; x += patternSize {
			for y := 0.0; y < h; y += patternSize {
				pdf.MoveTo(x+patternSize/2, y)
				pdf.LineTo(x, y+patternSize/2)
				pdf.LineTo(x+patternSize/2, y+patternSize)
				pdf.LineTo(x+patternSize, y+patternSize/2)
				pdf.LineTo(x+patternSize/2, y)
				pdf.ClosePath()
				pdf.DrawPath("D")
			}
		}
	case "rhombus":
		for x := 0.0; x < w; x += patternSize {
			for y := 0.0; y < h; y += patternHeight {
				pdf.MoveTo(x+patternSize/2, y)
				pdf.LineTo(x, y+patternHeight/2)
				pdf.LineTo(x+patternSize/2, y+patternHeight)
				pdf.LineTo(x+patternSize, y+patternHeight/2)
				pdf.LineTo(x+patternSize/2, y)
				pdf.ClosePath()
				pdf.DrawPath("D")
			}
		}
	case "triangles":
		for x := 0.0; x < w; x += patternSize {
			pdf.MoveTo(x, 0)
			pdf.LineTo(x, h)
		}

		for x := 0.0; x < w; x += patternSize {
			for y := 0.0; y < h; y += patternHeight {
				pdf.MoveTo(x+patternSize/2, y)
				pdf.LineTo(x, y+patternHeight/2)
				pdf.LineTo(x+patternSize/2, y+patternHeight)
				pdf.LineTo(x+patternSize/2, y)
				pdf.LineTo(x+patternSize, y+patternHeight/2)
				pdf.LineTo(x+patternSize/2, y+patternHeight)
				pdf.ClosePath()
				pdf.DrawPath("D")
			}
		}

	default:
		return errors.Wrap(errUnknownPattern, pattern.Name)
	}

	return nil
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
