package main

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/jung-kurt/gofpdf"
)

type pdfRequest struct {
	Size        string `json:"size"`
	Orientation string `json:"orientation"`
	Pattern     string `json:"pattern"`
}

func handlerPDF(w http.ResponseWriter, r *http.Request) {
	var req pdfRequest
	err := json.NewDecoder(r.Body).Decode(&req)
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

func createPDF(writer io.Writer, size, orientation, pattern string) error {
	pdf := gofpdf.New(orientation, "mm", size, "")
	pdf.SetAutoPageBreak(false, 0)

	err := pdf.Output(writer)
	if err != nil {
		return err
	}

	return nil
}
