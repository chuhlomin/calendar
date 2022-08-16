package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/jung-kurt/gofpdf"
)

type pdfRequest struct {
	Size        string  `json:"size"`
	Orientation string  `json:"orientation"`
	Pattern     pattern `json:"pattern"`
}

type pattern struct {
	Name string `json:"name"`
	Size string `json:"size"`
	size float64
}

func handlerPDF(w http.ResponseWriter, r *http.Request) {
	var req pdfRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// convert pattern size from string to float64
	req.Pattern.size, err = strconv.ParseFloat(req.Pattern.Size, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", "attachment; filename=grid.pdf")

	err = createPDF(
		w,
		req.Size,
		req.Orientation,
		req.Pattern,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func createPDF(writer io.Writer, size, orientation string, pattern pattern) error {
	pdf := gofpdf.New(orientation, "mm", size, "")
	pdf.SetAutoPageBreak(false, 0)

	pdf.AddPage()
	pdf.SetDrawColor(0, 0, 0)
	pdf.SetFillColor(255, 255, 255)

	err := drawPattern(pdf, pattern)
	if err != nil {
		return err
	}

	err = pdf.Output(writer)
	if err != nil {
		return err
	}

	return nil
}

func drawPattern(pdf gofpdf.Pdf, pattern pattern) error {
	w, h := pdf.GetPageSize()
	patternSize := pattern.size

	switch pattern.Name {
	case "rect":
		w, h := pdf.GetPageSize()
		for x := 0.0; x < w; x += patternSize {
			for y := 0.0; y < h; y += patternSize {
				pdf.Rect(x, y, patternSize, patternSize, "D")
			}
		}
	case "dot":
		pdf.SetFillColor(0, 0, 0)
		for x := 0.0; x < w; x += patternSize {
			for y := 0.0; y < h; y += patternSize {
				pdf.Circle(x, y, 0.2, "F")
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
		patternHeight := patternSize * 0.6
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
	default:
		return fmt.Errorf("unknown pattern: %s", pattern)
	}

	return nil
}
