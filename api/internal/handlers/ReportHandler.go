package handlers

import (
	"bookstore/api/api/internal/model"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

type ReportHandler struct {
	ReportsDir string
}

func NewReportHandler(reportsDir string) *ReportHandler {
	return &(ReportHandler{ReportsDir: reportsDir})
}

func (h *ReportHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var mergedReports []model.ReportModel

	err := filepath.Walk(h.ReportsDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && filepath.Ext(path) == ".json" {
			file, err := os.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()

			var report model.ReportModel
			bytes, err := io.ReadAll(file)
			if err != nil {
				return err
			}
			if err := json.Unmarshal(bytes, &report); err != nil {
				return err
			}
			mergedReports = append(mergedReports, report)
		}
		return nil
	})

	if err != nil {
		http.Error(w, "Failed to read reports", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(mergedReports); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}
