package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx"
)

// User user
type User struct {
	ID       string `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

var users []User
var connection *pgx.Conn

func mockUserData() {
	users = append(users, User{ID: "123", Email: "fake@exemple.com", Password: "lalala"})
	users = append(users, User{ID: "456", Email: "shake@exemple.com", Password: "tititi"})
}

// GetUsers find all users
func GetUsers(response http.ResponseWriter, request *http.Request) {
	var user User
	rows, _ := connection.Query("SELECT * FROM users")

	for rows.Next() {
		rows.Scan(&user.ID, &user.Email, &user.Password)
		users = append(users, user)
	}
	response.Header().Set("Content-type", "application/json")
	json.NewEncoder(response).Encode(users)
}

// GetUser find one user
func GetUser(response http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)

	response.Header().Set("Content-type", "application/json")
	for index := range users {
		if users[index].ID == params["id"] {
			json.NewEncoder(response).Encode(index)
			return
		}
	}
	json.NewEncoder(response).Encode(User{})
}

// CreateUser create a new user
func CreateUser(response http.ResponseWriter, request *http.Request) {
	var newUser User

	connection.Exec("INSERT INTO users")
	response.Header().Set("Content-type", "application/json")
	json.NewDecoder(request.Body).Decode(&newUser)
	users = append(users, newUser)
	json.NewEncoder(response).Encode(newUser)
}

// DeleteUser remove one user
func DeleteUser(response http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)

	response.Header().Set("Content-type", "application/json")
	for index := range users {
		if users[index].ID == params["id"] {
			users = append(users[:index], users[index+1:]...)
			break
		}
	}
	json.NewEncoder(response).Encode(users)
}

// UpdateUser update one user
func UpdateUser(response http.ResponseWriter, request *http.Request) {
	var user User
	params := mux.Vars(request)

	response.Header().Set("Content-type", "application/json")
	json.NewDecoder(request.Body).Decode(&user)
	for index := range users {
		if users[index].ID == params["id"] {
			users[index].Email = user.Email
			users[index].Password = user.Password
			json.NewEncoder(response).Encode(users[index])
			return
		}
	}
	json.NewEncoder(response).Encode(users)
}

func initRoutes(router *mux.Router) {
	router.HandleFunc("/api/users", GetUsers).Methods("GET")
	router.HandleFunc("/api/users/{id}", GetUser).Methods("GET")
	router.HandleFunc("/api/users", CreateUser).Methods("POST")
	router.HandleFunc("/api/users/{id}", DeleteUser).Methods("DELETE")
	router.HandleFunc("/api/users/{id}", UpdateUser).Methods("PUT")
}

func initDatabaseConnection() {
	config, err := pgx.ParseEnvLibpq()
	if err != nil {
		fmt.Println(os.Stderr, "Unable to parse Postgre environment: ", err)
		os.Exit(1)
	}
	connection, err = pgx.Connect(config)
	if err != nil {
		fmt.Println(os.Stderr, "Connection to database failed: ", err)
		os.Exit(1)
	}
}

func main() {
	router := mux.NewRouter()

	initDatabaseConnection()
	//mockUserData()
	initRoutes(router)
	log.Fatal(http.ListenAndServe(":8000", router))
}
