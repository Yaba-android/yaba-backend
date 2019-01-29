package user

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

// userExist check if user with this ID exist in database
func userExist(ID string) bool {
	rows, err := connection.Query("SELECT * FROM users WHERE user_id=$1", ID)
	if err != nil {
		fmt.Println("Unable to check if UserExist: ", err)
		return (false)
	}
	for rows.Next() {
		rows.Close()
		return (true)
	}
	rows.Close()
	return (false)
}

// convertUserDBToUser iterate throw query result and convert user in database to user in backend
func convertUserDBToUser(rows *pgx.Rows, users *[]User) bool {
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
	rows.Close()
	return (status)
}

// encodeUserResponse encode response with result or error
func encodeUserResponse(rows *pgx.Rows, response *http.ResponseWriter) {
	var users []User

	(*response).Header().Set("Content-type", "application/json")
	if convertUserDBToUser(rows, &users) {
		json.NewEncoder(*response).Encode(users)
	} else {
		json.NewEncoder(*response).Encode(http.StatusBadRequest)
	}
}

// getUsers find all users
func getUsers(response http.ResponseWriter, request *http.Request) {
	rows, err := connection.Query("SELECT * FROM users")

	if err != nil {
		fmt.Println("Unable to exec GetUsers query: ", err)
	}
	encodeUserResponse(rows, &response)
}

// getUser find one user
func getUser(response http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	rows, err := connection.Query("SELECT * FROM users WHERE user_id=$1", params["id"])

	if err != nil {
		fmt.Println("Unable to exec GetUser query: ", err)
	}
	encodeUserResponse(rows, &response)
}

// createUser create a new user
func createUser(response http.ResponseWriter, request *http.Request) {
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

// deleteUser remove one user
func deleteUser(response http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)

	response.Header().Set("Content-type", "application/json")
	if !userExist(params["id"]) {
		fmt.Println("Unable to DeleteUser: User with this ID doesn't exist")
		json.NewEncoder(response).Encode(http.StatusBadRequest)
		return
	}
	_, err := connection.Exec("DELETE FROM users WHERE user_id=$1", params["id"])
	if err != nil {
		fmt.Println("Unable to exec DeleteUser query: ", err)
		json.NewEncoder(response).Encode(http.StatusBadRequest)
		return
	}
	json.NewEncoder(response).Encode(http.StatusAccepted)
}

// updateUser update one user
func updateUser(response http.ResponseWriter, request *http.Request) {
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

// initRoutes initialize endpoints for the REST API
func initRoutes(router *mux.Router) {
	router.HandleFunc("/api/users", getUsers).Methods("GET")
	router.HandleFunc("/api/users/{id}", getUser).Methods("GET")
	router.HandleFunc("/api/users", createUser).Methods("POST")
	router.HandleFunc("/api/users/{id}", deleteUser).Methods("DELETE")
	router.HandleFunc("/api/users/{id}", updateUser).Methods("PUT")
}

// initDatabaseConnection initialize connection with the database
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

// InitUserRouter init all user routes to communicate with database
func InitUserRouter() {
	router := mux.NewRouter()

	initDatabaseConnection()
	initRoutes(router)
	log.Fatal(http.ListenAndServe(":8000", router))
}
