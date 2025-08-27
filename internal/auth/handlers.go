package auth

import (
	"net/http"
	"encoding/json"
	"log"
	"strings"
)

func RegisterHandler(authService *AuthService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var LoginPasswordJSON LoginPasswordJSON

		err := json.NewDecoder(r.Body).Decode(&LoginPasswordJSON)
		if err != nil {
			http.Error(w, `Bad Request`, http.StatusBadRequest)
			return
		}

		login := LoginPasswordJSON.Username
		password := LoginPasswordJSON.Password

		if authService.Exist(login) {
			http.Error(w, `User already exists`, http.StatusConflict)
			return
		}

		err = authService.Insert(login, password)
		if err != nil {
			log.Println("DB update error:", err)
			http.Error(w, `Server error`, http.StatusInternalServerError)
			return
		}

		tokenString, err := GenerateToken(login)
		if err != nil {
			log.Println("Token creation error:", err)
			http.Error(w, `Server error`, http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)

		json.NewEncoder(w).Encode(TokenResponse{Token: tokenString})
	}	
}

func LoginHandler(authService *AuthService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var loginJSON LoginPasswordJSON

		err := json.NewDecoder(r.Body).Decode(&loginJSON)
		if err != nil {
			log.Println("Error during JSON parsing: ", err)
			http.Error(w, `Bad Request`, http.StatusBadRequest)
			return
		}

		login := loginJSON.Username
		password := loginJSON.Password

		err = authService.UserValidation(login, password)
		if err != nil {
			log.Println("Error during user validation: ", err)
			http.Error(w, `Incorrect login or password`, http.StatusUnauthorized)
			return 
		}

		tokenString, err := GenerateToken(login)
		if err != nil {
			log.Println("Token creation error:", err)
			http.Error(w, `Server error`, http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		json.NewEncoder(w).Encode(TokenResponse{Token: tokenString})
	}	
}

func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Missing Authorization header", http.StatusUnauthorized)
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			http.Error(w, "Invalid Authorization header format", http.StatusUnauthorized)
			return
		}

		tokenString := parts[1]

		_, err := ValidateToken(tokenString)
		if err != nil {
			http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
			return
		}

		next(w, r)
	}
}