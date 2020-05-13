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

func addNewAuthorLink(w http.ResponseWriter, r *http.Request, client *redis.Client) {
	var author Author

	w.Header().Set("Content-type", "application/json")
	err := json.NewDecoder(r.Body).Decode(&author)
	if err != nil {
		fmt.Println("error: ", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Error bad request"))
	} else {
		authorAdded := redisSetNewAuthor(client, &author)
		if authorAdded == nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Error internal server"))
			return
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(authorAdded)
	}
}

func getAllAuthorsLink(w http.ResponseWriter, r *http.Request, client *redis.Client) {
	authors := redisGetAllAuthors(client)

	w.Header().Set("Content-type", "application/json")
	if authors == nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Error not found"))
		return
	}
	json.NewEncoder(w).Encode(authors)
}

func getAuthorLink(w http.ResponseWriter, r *http.Request, client *redis.Client) {
	remoteId := mux.Vars(r)
	author := redisGetAuthorById(client, remoteId["RemoteId"])

	w.Header().Set("Content-type", "application/json")
	if author == nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Error not found"))
		return
	}
	json.NewEncoder(w).Encode(author)
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
		bookAdded := redisSetNewBook(client, &book)
		if bookAdded == nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Error internal server"))
			return
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(bookAdded)
	}
}

func getAllBooksLink(w http.ResponseWriter, r *http.Request, client *redis.Client) {
	books := redisGetAllBooks(client)

	w.Header().Set("Content-type", "application/json")
	if books == nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Error not found"))
		return
	}
	json.NewEncoder(w).Encode(books)
}

func getBookLink(w http.ResponseWriter, r *http.Request, client *redis.Client) {
	remoteId := mux.Vars(r)
	book := redisGetBookById(client, remoteId["RemoteId"])

	w.Header().Set("Content-type", "application/json")
	if book == nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Error not found"))
		return
	}
	json.NewEncoder(w).Encode(book)
}

func addNewGenreLink(w http.ResponseWriter, r *http.Request, client *redis.Client) {
	var genre Genre

	w.Header().Set("Content-type", "application/json")
	err := json.NewDecoder(r.Body).Decode(&genre)
	if err != nil {
		fmt.Println("error: ", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Error bad request"))
	} else {
		genreAdded := redisSetNewGenre(client, &genre)
		if genreAdded == nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Error internal server"))
			return
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(genreAdded)
	}
}

func getAllGenresLink(w http.ResponseWriter, r *http.Request, client *redis.Client) {
	genres := redisGetAllGenres(client)

	w.Header().Set("Content-type", "application/json")
	if genres == nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Error not found"))
		return
	}
	json.NewEncoder(w).Encode(genres)
}

func startRouter(client *redis.Client) {
	if redisIsClientConnected(client) == nil {
		router := mux.NewRouter().StrictSlash(true)

		router.
			PathPrefix(EbooksDir).
			Handler(http.StripPrefix(EbooksDir, http.FileServer(http.Dir("."+EbooksDir))))

		router.HandleFunc("/", homeLink).Methods("GET")

		router.HandleFunc("/author", func(w http.ResponseWriter, r *http.Request) {
			addNewAuthorLink(w, r, client)
		}).Methods("POST")

		router.HandleFunc("/authors", func(w http.ResponseWriter, r *http.Request) {
			getAllAuthorsLink(w, r, client)
		}).Methods("GET")

		router.HandleFunc("/author/{RemoteId}", func(w http.ResponseWriter, r *http.Request) {
			getAuthorLink(w, r, client)
		}).Methods("GET")

		router.HandleFunc("/book", func(w http.ResponseWriter, r *http.Request) {
			addNewBookLink(w, r, client)
		}).Methods("POST")

		router.HandleFunc("/books", func(w http.ResponseWriter, r *http.Request) {
			getAllBooksLink(w, r, client)
		}).Methods("GET")

		router.HandleFunc("/book/{RemoteId}", func(w http.ResponseWriter, r *http.Request) {
			getBookLink(w, r, client)
		}).Methods("GET")

		router.HandleFunc("/genre", func(w http.ResponseWriter, r *http.Request) {
			addNewGenreLink(w, r, client)
		}).Methods("POST")

		router.HandleFunc("/genres", func(w http.ResponseWriter, r *http.Request) {
			getAllGenresLink(w, r, client)
		}).Methods("GET")

		log.Fatal(http.ListenAndServe(":"+os.Getenv(MuxRouterPort), router))
	}
}
