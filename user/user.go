package user

import (
	"encoding/json"
	"net/http"

	"github.com/peksinsara/e-voting-RDBMS/database"
)

// User represents the user registration data
type User struct {
	UserID        int    `json:"user_id"`
	FullName      string `json:"full_name"`
	MothersName   string `json:"mothers_name"`
	Email         string `json:"email"`
	PhoneNumber   string `json:"phone_number"`
	JMBG          string `json:"jmbg"`
	Password      string `json:"password"`
	IsAdmin       bool   `json:"is_admin"`
	AdminUsername string `json:"admin_username"`
	AdminPassword string `json:"admin_password"`
}

// LoginRequest represents the user login request data
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LoginUser handles the user login endpoint
func LoginUser(w http.ResponseWriter, r *http.Request) {
	var loginReq LoginRequest
	err := json.NewDecoder(r.Body).Decode(&loginReq)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	db := database.GetDB()

	var user User
	err = db.QueryRow("SELECT user_id, full_name FROM User WHERE email = ? AND password = ?", loginReq.Email, loginReq.Password).Scan(&user.UserID, &user.FullName)
	if err != nil {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	loginRes := map[string]interface{}{
		"user_id":   user.UserID,
		"full_name": user.FullName,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(loginRes)
}

// RegisterUser handles the user registration endpoint
func RegisterUser(w http.ResponseWriter, r *http.Request) {
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	db := database.GetDB()

	stmt, err := db.Prepare("INSERT INTO User (full_name, mothers_name, email, phone_number, jmbg, password) VALUES (?, ?, ?, ?, ?, ?)")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(user.FullName, user.MothersName, user.Email, user.PhoneNumber, user.JMBG, user.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("User registered successfully"))
}
