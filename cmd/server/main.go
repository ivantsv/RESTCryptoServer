package main

import (
	"RESTCryptoServer/internal/auth"
	"RESTCryptoServer/internal/db"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func main() {
	userdb, err := db.NewUserDB()
	if err != nil {
		log.Println("error during opening/creation users postgres db: ", err)
		return
	}
	defer userdb.Close()

	authService := auth.NewAuthService(userdb)

	router := chi.NewRouter()
	router.Post("/auth/register", auth.RegisterHandler(authService))
	router.Post("/auth/login", auth.LoginHandler(authService))

	log.Println("Starting server on http://:8080")
	err = http.ListenAndServe(":8080", router)
	if err != nil {
		log.Fatalln(err)
	}
}