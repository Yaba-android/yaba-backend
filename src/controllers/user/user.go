package user

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx"
)

// User user
type User struct {
	ID       string `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// DbHandler database connection handler
var DbHandler *pgx.Conn

// userExist check if user with this ID exist in database
func userExist(ID string) bool {
	rows, err := DbHandler.Query("SELECT * FROM users WHERE user_id=$1", ID)
	if err != nil {
		fmt.Println("Unable to check if userExist: ", err)
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
			fmt.Println("Unable to convertUserDBToUser: ", err)
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

// GetUsers find all users
func GetUsers(response http.ResponseWriter, request *http.Request) {
	rows, err := DbHandler.Query("SELECT * FROM users")

	if err != nil {
		fmt.Println("Unable to exec GetUsers query: ", err)
	}
	encodeUserResponse(rows, &response)
}

// GetUser find one user
func GetUser(response http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	rows, err := DbHandler.Query("SELECT * FROM users WHERE user_id=$1", params["id"])

	if err != nil {
		fmt.Println("Unable to exec GetUser query: ", err)
	}
	encodeUserResponse(rows, &response)
}

// CreateUser create a new user
func CreateUser(response http.ResponseWriter, request *http.Request) {
	var newUser User

	response.Header().Set("Content-type", "application/json")
	json.NewDecoder(request.Body).Decode(&newUser)
	if newUser.Email == "" || newUser.Password == "" {
		fmt.Println("Unable to CreateUser: Email or Password cannot be empty")
		json.NewEncoder(response).Encode(http.StatusBadRequest)
		return
	}
	_, err := DbHandler.Exec("INSERT INTO users(user_id, email, password) values($1, $2, $3)", newUser.ID, newUser.Email, newUser.Password)
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
	if !userExist(params["id"]) {
		fmt.Println("Unable to DeleteUser: User with this ID doesn't exist")
		json.NewEncoder(response).Encode(http.StatusBadRequest)
		return
	}
	_, err := DbHandler.Exec("DELETE FROM users WHERE user_id=$1", params["id"])
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
		_, err := DbHandler.Exec("UPDATE users SET email=$1 WHERE user_id=$2", user.Email, params["id"])
		if err != nil {
			fmt.Println("Unable to exec UpdateUser query: ", err)
			json.NewEncoder(response).Encode(http.StatusBadRequest)
			return
		}
	}
	if user.Password != "" {
		_, err := DbHandler.Exec("UPDATE users SET password=$1 WHERE user_id=$2", user.Password, params["id"])
		if err != nil {
			fmt.Println("Unable to exec UpdateUser query: ", err)
			json.NewEncoder(response).Encode(http.StatusBadRequest)
			return
		}
	}
	json.NewEncoder(response).Encode(http.StatusAccepted)
}
