package library

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx"
)

// Library library
type Library struct {
	userID   string `json:"user_id"`
	ebooksID string `json:"ebooks_id"`
}

// DbHandler database connection handler
var DbHandler *pgx.Conn

// GetLibrary find library of the user_id
func GetLibrary(response http.ResponseWriter, request *http.Request) {
}

// GetEbook find ebook in the library of the user_id
func GetEbook(response http.ResponseWriter, request *http.Request) {
}

// AddEbook add a new ebook in the library of the user_id
func AddEbook(response http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)

	response.Header().Set("Content-type", "application/json")
	// check si le ebook existe bien et l'user id aussi, et que l'ebook n'est pas deux fois, plus a ajouter avec virgules
	_, err := DbHandler.Exec("INSERT INTO libraries(ebooks_id) values($1) WHERE user_id=$2", params["ebook_id"], params["id"])
	if err != nil {
		fmt.Println("Unable to exec CreateUser query: ", err)
		json.NewEncoder(response).Encode(http.StatusBadRequest)
		return
	}
	json.NewEncoder(response).Encode(http.StatusAccepted)
}
