package yaba

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-redis/redis"
	"github.com/gorilla/mux"
)

/********************************
*
* 		ROUTER GORILLA MUX
*
********************************/

func homeLink(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "welcome")
}

func addNewBookLink(w http.ResponseWriter, r *http.Request, client *redis.Client) {
	book := &Book{Name: "The Book", Author: "John Doe", DatePub: "10/09/2015", Path: "epub/thebook.epub"}
	redisSetNewBook(client, book)
}

func getAllBooksLink(w http.ResponseWriter, r *http.Request, client *redis.Client) {
	redisGetAllBooks(client)
}

func startRouter(client *redis.Client) {
	if redisIsClientConnected(client) == nil {
		router := mux.NewRouter().StrictSlash(true)
		router.HandleFunc("/", homeLink)
		router.HandleFunc("/addNewBook", func(w http.ResponseWriter, r *http.Request) {
			addNewBookLink(w, r, client)
		})
		router.HandleFunc("/getAllBooks", func(w http.ResponseWriter, r *http.Request) {
			getAllBooksLink(w, r, client)
		})
		log.Fatal(http.ListenAndServe(":8080", router))
	}
}
