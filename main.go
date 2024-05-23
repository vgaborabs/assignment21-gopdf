package main

import (
	"github.com/vgaborabs/assignment21-gopdf/handler"
	"log"
	"net/http"
)

func main() {

	// Handle the statement PDF as file
	http.HandleFunc("POST /statement/pdf", handler.HandleStatementPdf)
	// Handle the statement PDF as stream
	http.HandleFunc("POST /statement/pdf/stream", handler.HandleStatementPdfStream)
	// Serve the saved files
	http.Handle("/public/", http.StripPrefix("/public", http.FileServer(http.Dir("public"))))

	log.Fatal(http.ListenAndServe(":3000", nil))

}
