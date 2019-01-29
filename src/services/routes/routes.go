package routes

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/nasrat_v/maktaba-android-mvp/src/controllers/ebook"
	"github.com/nasrat_v/maktaba-android-mvp/src/controllers/library"
	"github.com/nasrat_v/maktaba-android-mvp/src/controllers/user"
)

// initUserRoutes initialize user endpoints for the REST API
func initUserRoutes(router *mux.Router) {
	router.HandleFunc("/api/users", user.GetUsers).Methods("GET")
	router.HandleFunc("/api/users/{id}", user.GetUser).Methods("GET")
	router.HandleFunc("/api/users", user.CreateUser).Methods("POST")
	router.HandleFunc("/api/users/{id}", user.DeleteUser).Methods("DELETE")
	router.HandleFunc("/api/users/{id}", user.UpdateUser).Methods("PUT")
}

// initEbookRoutes initialize ebooks endpoints for the REST API
func initEbookRoutes(router *mux.Router) {
	router.HandleFunc("/api/ebooks", ebook.GetEbooks).Methods("GET")
	router.HandleFunc("/api/ebooks/{id}", ebook.GetEbook).Methods("GET")
}

func initLibraryRoutes(router *mux.Router) {
	router.HandleFunc("/api/users/{id}/library", library.GetLibrary).Methods("GET")
	router.HandleFunc("/api/users/{id}/library/ebook/{ebook_id}", library.GetEbook).Methods("GET")
	router.HandleFunc("/api/users/{id}/library/ebook/{ebook_id}", library.AddEbook).Methods("PUT")
}

// InitRouterForControllers initialize all endpoints for the REST API
func InitRouterForControllers() {
	router := mux.NewRouter()

	initUserRoutes(router)
	initEbookRoutes(router)
	initLibraryRoutes(router)
	log.Fatal(http.ListenAndServe(":8000", router))
}
