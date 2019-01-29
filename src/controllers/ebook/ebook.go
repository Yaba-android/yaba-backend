package ebook

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx"
)

// Ebook ebook
type Ebook struct {
	ID       string `json:"id"`
	Tilte    string `json:"title"`
	Author   string `json:"author"`
	Filepath string `json:"filepath"`
}

// DbHandler database connection handler
var DbHandler *pgx.Conn

// convertEbookDBToEbook iterate throw query result and convert ebook in database to ebook in backend
func convertEbookDBToEbook(rows *pgx.Rows, ebooks *[]Ebook) bool {
	var ebook Ebook
	var ID uint32
	status := false

	for rows.Next() {
		err := rows.Scan(&ID, &ebook.Tilte, &ebook.Author, &ebook.Filepath)
		if err != nil {
			fmt.Println("Unable to ConvertUserDBToUser: ", err)
		}
		ebook.ID = fmt.Sprint(ID) // convert serial (uint32) to string
		*ebooks = append(*ebooks, ebook)
		status = true
	}
	rows.Close()
	return (status)
}

// encodeEbookResponse encode response with result or error
func encodeEbookResponse(rows *pgx.Rows, response *http.ResponseWriter) {
	var ebooks []Ebook

	(*response).Header().Set("Content-type", "application/json")
	if convertEbookDBToEbook(rows, &ebooks) {
		json.NewEncoder(*response).Encode(ebooks)
	} else {
		json.NewEncoder(*response).Encode(http.StatusBadRequest)
	}
}

// GetEbooks find all Ebooks
func GetEbooks(response http.ResponseWriter, request *http.Request) {
	rows, err := DbHandler.Query("SELECT * FROM ebooks")

	if err != nil {
		fmt.Println("Unable to exec GetEbooks query: ", err)
	}
	encodeEbookResponse(rows, &response)
}

// GetEbook find one Ebook
func GetEbook(response http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	rows, err := DbHandler.Query("SELECT * FROM ebooks WHERE ebook_id=$1", params["id"])

	if err != nil {
		fmt.Println("Unable to exec GetEbook query: ", err)
	}
	encodeEbookResponse(rows, &response)
}
