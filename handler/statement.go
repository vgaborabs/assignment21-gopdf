package handler

import (
	"assignment21-gopdf/pkg/pdf"
	"assignment21-gopdf/types"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

// HandleStatementPdfStream generates the PDF statement for the types.Statement in the request body and writes it to the response
func HandleStatementPdfStream(w http.ResponseWriter, r *http.Request) {
	var statement types.Statement
	err := json.NewDecoder(r.Body).Decode(&statement)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	reader, err := pdf.GenerateStatement(statement, nil)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	http.ServeContent(w, r, getFileName(statement), time.Time{}, reader)
}

// HandleStatementPdf generates the PDF statement for the types.Statement in the request body and saves it to a file, returning the relative path where it can be later retrieved
func HandleStatementPdf(w http.ResponseWriter, r *http.Request) {
	var statement types.Statement
	err := json.NewDecoder(r.Body).Decode(&statement)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	fileName := fmt.Sprintf("public/%s", getFileName(statement))

	file, err := os.Create(fileName)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	_, err = pdf.GenerateStatement(statement, file)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = file.Close()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	_, err = w.Write([]byte(fileName))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func getFileName(statement types.Statement) string {
	return fmt.Sprintf("%s_%s-%s_%s.pdf", statement.FullName, statement.StartDate.Format("2006-01-02"), statement.EndDate.Format("2006-01-02"), time.Now().Format("20060102150405"))
}
