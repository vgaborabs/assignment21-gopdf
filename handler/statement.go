package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/vgaborabs/assignment21-gopdf/pkg/pdf"
	"github.com/vgaborabs/assignment21-gopdf/types"
	"log/slog"
	"net/http"
	"os"
	"time"
)

// HandleStatementPdfStream generates the PDF statement for the types.Statement in the request body and writes it to the response as an octet-stream
func HandleStatementPdfStream(w http.ResponseWriter, r *http.Request) {
	statement, err := getStatement(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	b, err := pdf.GenerateStatement(statement)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	reader := bytes.NewReader(b)
	http.ServeContent(w, r, getFileName(statement), time.Time{}, reader)
}

// HandleStatementPdf generates the PDF statement for the types.Statement in the request body and saves it to a file, returning the relative path where it can be later retrieved
func HandleStatementPdf(w http.ResponseWriter, r *http.Request) {
	statement, err := getStatement(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	pf, err := os.Open("public")
	if os.IsNotExist(err) {
		slog.Info("Public folder not exists, creating... ")

		err = os.Mkdir("public", 0755)
		if err != nil {
			slog.Error("Error creating public folder: ", "err", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		pf, _ = os.Open("public")
	} else if err != nil {
		slog.Error("Error opening public folder: ", "err", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err = pf.Close(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	fileName := fmt.Sprintf("public/%s", getFileName(statement))

	file, err := os.Create(fileName)
	if err != nil {
		slog.Error("Error creating file: ", "file", fileName, "err", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer func() {
		err = file.Close()
		if err != nil {
			slog.Error("Failed to close file", "file", fileName, "err", err)
		}
	}()

	b, err := pdf.GenerateStatement(statement)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	_, err = file.Write(b)
	if err != nil {
		slog.Error("Failed to write file", "file", fileName, "err", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, err = w.Write([]byte(fileName))
	if err != nil {
		slog.Error("Failed to write response", "err", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func getFileName(statement types.Statement) string {
	return fmt.Sprintf("%s_%s-%s_%s.pdf", statement.FullName, statement.StartDate.Format("2006-01-02"), statement.EndDate.Format("2006-01-02"), time.Now().Format("20060102150405"))
}

func getStatement(r *http.Request) (types.Statement, error) {
	var statement types.Statement
	err := json.NewDecoder(r.Body).Decode(&statement)
	if err != nil {
		slog.Error("Failed to decode request body to Statement", "err", err)
	}
	return statement, err
}
