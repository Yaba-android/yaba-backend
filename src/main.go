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

var connection *pgx.Conn

// ConvertUserDBToUser iterate throw query result and convert user in database to user in backend
func ConvertUserDBToUser(rows *pgx.Rows, users *[]User) bool {
	var user User
	var ID uint32
	status := false

	for rows.Next() {
		err := rows.Scan(&ID, &user.Email, &user.Password)
		if err != nil {
			fmt.Println("Unable to find user: ", err)
		}
		user.ID = fmt.Sprint(ID) // convert serial (uint32) to string
		*users = append(*users, user)
		status = true
	}
	return (status)
}

// EncodeUserResponse encode response with result or error
func EncodeUserResponse(rows *pgx.Rows, response *http.ResponseWriter) {
	var users []User

	(*response).Header().Set("Content-type", "application/json")
	if ConvertUserDBToUser(rows, &users) {
		json.NewEncoder(*response).Encode(users)
	} else {
		json.NewEncoder(*response).Encode(http.StatusBadRequest)
	}
}

// GetUsers find all users
func GetUsers(response http.ResponseWriter, request *http.Request) {
	rows, err := connection.Query("SELECT * FROM users")

	if err != nil {
		fmt.Println("Unable to exec GetUsers query: ", err)
	}
	EncodeUserResponse(rows, &response)
}

// GetUser find one user
func GetUser(response http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	rows, err := connection.Query("SELECT * FROM users WHERE user_id=$1", params["id"])

	if err != nil {
		fmt.Println("Unable to exec GetUser query: ", err)
	}
	EncodeUserResponse(rows, &response)
}

// CreateUser create a new user
func CreateUser(response http.ResponseWriter, request *http.Request) {
	var newUser User

	response.Header().Set("Content-type", "application/json")
	json.NewDecoder(request.Body).Decode(&newUser)
	if newUser.Email == "" || newUser.Password == "" {
		fmt.Println("Unable CreateUser: Email or Password cannot be empty")
		json.NewEncoder(response).Encode(http.StatusBadRequest)
		return
	}
	_, err := connection.Exec("INSERT INTO users(user_id, email, password) values($1, $2, $3)", newUser.ID, newUser.Email, newUser.Password)
	if err != nil {
		fmt.Println("Unable to exec CreateUser query: ", err)
		json.NewEncoder(response).Encode(http.StatusBadRequest)
		return
	}
	json.NewEncoder(response).Encode(http.StatusAccepted)
}

// DeleteUser remove one user
func DeleteUser(response http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)

	response.Header().Set("Content-type", "application/json")
	_, err := connection.Exec("DELETE FROM users WHERE user_id=$1", params["id"])
	if err != nil {
		fmt.Println("Unable to exec DeleteUser query: ", err)
		json.NewEncoder(response).Encode(http.StatusBadRequest)
		return
	}
	json.NewEncoder(response).Encode(http.StatusAccepted)
}

// UpdateUser update one user
func UpdateUser(response http.ResponseWriter, request *http.Request) {
	var user User
	params := mux.Vars(request)

	response.Header().Set("Content-type", "application/json")
	json.NewDecoder(request.Body).Decode(&user)
	if user.Email == "" && user.Password == "" {
		fmt.Println("Unable to UpdateUser: Email and Password cannot be empties")
		json.NewEncoder(response).Encode(http.StatusBadRequest)
		return
	}
	if user.Email != "" {
		_, err := connection.Exec("UPDATE users SET email=$1 WHERE user_id=$2", user.Email, params["id"])
		if err != nil {
			fmt.Println("Unable to exec UpdateUser query: ", err)
			json.NewEncoder(response).Encode(http.StatusBadRequest)
			return
		}
	}
	if user.Password != "" {
		_, err := connection.Exec("UPDATE users SET password=$1 WHERE user_id=$2", user.Password, params["id"])
		if err != nil {
			fmt.Println("Unable to exec UpdateUser query: ", err)
			json.NewEncoder(response).Encode(http.StatusBadRequest)
			return
		}
	}
	json.NewEncoder(response).Encode(http.StatusAccepted)
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
		fmt.Println("Unable to parse Postgre environment: ", err)
		os.Exit(1)
	}
	connection, err = pgx.Connect(config)
	if err != nil {
		fmt.Println("Connection to database failed: ", err)
		os.Exit(1)
	}
}

func main() {
	router := mux.NewRouter()

	initDatabaseConnection()
	initRoutes(router)
	log.Fatal(http.ListenAndServe(":8000", router))
}
