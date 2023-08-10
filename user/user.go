package user

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/peksinsara/e-voting-RDBMS/database"
)

// User represents the user registration data
type User struct {
	UserID      int    `json:"user_id"`
	FullName    string `json:"full_name"`
	MothersName string `json:"mothers_name"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phone_number"`
	JMBG        string `json:"jmbg"`
	Password    string `json:"password"`
	IsAdmin     bool   `json:"is_admin"`
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

	token, err := GenerateJWT(user.UserID)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Authorization", "Bearer "+token)

	response := map[string]interface{}{
		"token":     token,
		"user_id":   user.UserID,
		"full_name": user.FullName,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// RegisterUser handles the user registration endpoint
func RegisterUser(w http.ResponseWriter, r *http.Request) {
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Perform validation checks
	if !isValidPhoneNumber(user.PhoneNumber) {
		http.Error(w, "Invalid phone number format", http.StatusBadRequest)
		return
	}

	if !isValidJMBG(user.JMBG) {
		http.Error(w, "Invalid JMBG format", http.StatusBadRequest)
		return
	}

	db := database.GetDB()

	// Check if email, JMBG, or phone number already exist in the database
	var existingUser User
	err = db.QueryRow("SELECT user_id FROM User WHERE email = ? OR jmbg = ? OR phone_number = ?", user.Email, user.JMBG, user.PhoneNumber).Scan(&existingUser.UserID)
	if err == nil {
		http.Error(w, "User with the same email, JMBG, or phone number already exists", http.StatusConflict)
		return
	} else if err != sql.ErrNoRows {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	stmt, err := db.Prepare("INSERT INTO User (full_name, mothers_name, email, phone_number, jmbg, password, is_admin) VALUES (?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(user.FullName, user.MothersName, user.Email, user.PhoneNumber, user.JMBG, user.Password, user.IsAdmin)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("User registered successfully"))
}

// GetProfile handles the user profile endpoint
func GetProfile(w http.ResponseWriter, r *http.Request) {
	userID, err := GetUserIDFromToken(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	db := database.GetDB()

	var user User
	err = db.QueryRow("SELECT full_name, mothers_name, email, phone_number, jmbg, is_admin FROM User WHERE user_id=?", userID).
		Scan(&user.FullName, &user.MothersName, &user.Email, &user.PhoneNumber, &user.JMBG, &user.IsAdmin)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func GetUserIDFromToken(r *http.Request) (int, error) {
	tokenString := r.Header.Get("Authorization")
	if tokenString == "" {
		return 0, errors.New("Authorization token not found")
	}

	parts := strings.Split(tokenString, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return 0, errors.New("Invalid authorization header format")
	}

	token := parts[1]

	jwtToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte("720439bnskdjfad7skdfba0snlk"), nil
	})

	if err != nil {
		return 0, err
	}

	if claims, ok := jwtToken.Claims.(jwt.MapClaims); ok && jwtToken.Valid {
		userID, ok := claims["user_id"].(float64)
		if !ok {
			return 0, errors.New("Invalid user_id claim")
		}
		return int(userID), nil
	}

	return 0, errors.New("Invalid token or token expired")
}

func GenerateJWT(userID int) (string, error) {
	signingKey := []byte("720439bnskdjfad7skdfba0snlk")

	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(), // Token expiration time (24 hours)
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(signingKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// LogoutUser handles the user logout endpoint
func LogoutUser(w http.ResponseWriter, r *http.Request) {
	// To log out, we can simply invalidate the token by setting it to an empty string
	// This will cause the client's stored token to become invalid, effectively logging the user out
	w.Header().Set("Authorization", "")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("User logged out successfully"))
}

func isValidPhoneNumber(phoneNumber string) bool {
	return len(phoneNumber) >= 9 && len(phoneNumber) <= 12
}

func isValidJMBG(jmbg string) bool {
	return len(jmbg) == 13
}
