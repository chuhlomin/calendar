package main

import (
	"encoding/json"
	"fmt"
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
	Color     string `json:"color"`
	LineWidth string `json:"lineWidth"`
	size      float64
	lineWidth float64
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
		log.Println("error converting pattern size:", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	req.Pattern.lineWidth, err = strconv.ParseFloat(req.Pattern.LineWidth, 64)
	if err != nil {
		log.Println("error converting line width:", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", "attachment; filename=grid.pdf")

	err = createPDF(w, req.Size, req.Orientation, req.Pattern)
	if err != nil {
		log.Println("error creating pdf:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
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
	pdf.SetFillColor(255, 255, 255)
	pdf.SetDrawColor(convertColor(pattern.Color))

	w, h := pdf.GetPageSize()
	patternSize := pattern.size

	pdf.SetLineWidth(pattern.lineWidth / 1000)

	switch pattern.Name {
	case "rect":
		w, h := pdf.GetPageSize()
		for x := 0.0; x < w; x += patternSize {
			for y := 0.0; y < h; y += patternSize {
				pdf.Rect(x, y, patternSize, patternSize, "D")
			}
		}
	case "dot":
		pdf.SetFillColor(convertColor(pattern.Color))
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
		return fmt.Errorf("unknown pattern: %s", pattern.Name)
	}

	return nil
}

func convertColor(hex string) (int, int, int) {
	var r, g, b int64
	if len(hex) == 7 {
		r, _ = strconv.ParseInt(hex[1:3], 16, 0)
		g, _ = strconv.ParseInt(hex[3:5], 16, 0)
		b, _ = strconv.ParseInt(hex[5:7], 16, 0)
	}
	return int(r), int(g), int(b)
}
