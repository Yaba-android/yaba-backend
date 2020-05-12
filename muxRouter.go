package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

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
	var book Book

	w.Header().Set("Content-type", "application/json")
	err := json.NewDecoder(r.Body).Decode(&book)
	if err != nil {
		fmt.Println("error: ", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Error bad request"))
	} else {
		err = redisSetNewBook(client, &book)
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Error internal server"))
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Book added"))
	}
}

func getAllBooksLink(w http.ResponseWriter, r *http.Request, client *redis.Client) {
	books := redisGetAllBooks(client)

	w.Header().Set("Content-type", "application/json")
	if books == nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error bad request"))
	}
	json.NewEncoder(w).Encode(books)
}

func startRouter(client *redis.Client) {
	if redisIsClientConnected(client) == nil {
		router := mux.NewRouter().StrictSlash(true)

		router.
			PathPrefix(EbooksDir).
			Handler(http.StripPrefix(EbooksDir, http.FileServer(http.Dir("."+EbooksDir))))

		router.HandleFunc("/", homeLink).Methods("GET")

		router.HandleFunc("/addNewBook", func(w http.ResponseWriter, r *http.Request) {
			addNewBookLink(w, r, client)
		}).Methods("POST")

		router.HandleFunc("/getAllBooks", func(w http.ResponseWriter, r *http.Request) {
			getAllBooksLink(w, r, client)
		}).Methods("GET")

		log.Fatal(http.ListenAndServe(":"+os.Getenv(MuxRouterPort), router))
	}
}
