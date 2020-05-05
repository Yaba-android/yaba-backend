package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
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
	nb := rand.Int63n(100)
	title := concatStringInt("The Book", nb)
	book := &Book{Name: title, Author: "John Doe", DatePub: "10/09/2015", Filename: "camus_la_peste.epub"}

	w.Header().Set("Content-type", "application/json")
	redisSetNewBook(client, book)
	w.Write([]byte("Book added"))
}

func getAllBooksLink(w http.ResponseWriter, r *http.Request, client *redis.Client) {
	books := redisGetAllBooks(client)

	w.Header().Set("Content-type", "application/json")
	if books == nil {
		w.Write([]byte("error"))
	}
	json.NewEncoder(w).Encode(books)
}

func startRouter(client *redis.Client) {
	if redisIsClientConnected(client) == nil {
		router := mux.NewRouter().StrictSlash(true)

		router.
			PathPrefix(EbooksDir).
			Handler(http.StripPrefix(EbooksDir, http.FileServer(http.Dir("."+EbooksDir))))
		router.HandleFunc("/", homeLink)
		router.HandleFunc("/addNewBook", func(w http.ResponseWriter, r *http.Request) {
			addNewBookLink(w, r, client)
		})
		router.HandleFunc("/getAllBooks", func(w http.ResponseWriter, r *http.Request) {
			getAllBooksLink(w, r, client)
		})
		log.Fatal(http.ListenAndServe(MuxRouterPort, router))
	}
}
