package main

import (
	"RESTCryptoServer/internal/auth"
	"RESTCryptoServer/internal/db"
	"RESTCryptoServer/internal/redis"
	"RESTCryptoServer/internal/crypto"
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

	cryptodb, err := db.NewCryptoDB()
	if err != nil {
		log.Println("error during opening/creation crypto postgres db: ", err)
		return
	}
	defer cryptodb.Close()

	cache, err := redis.NewRedisClient()
	if err != nil {
		log.Println("error during cache redis db: ", err)
		return
	}
	defer cache.Close()


	log.Println("Connected to PostgreSQL")
	log.Println("All database migrations completed")

	authService := auth.NewAuthService(userdb)
	cryptoService := crypto.NewCryptoService(cryptodb, cache)

	router := chi.NewRouter()
	router.Post("/auth/register", auth.RegisterHandler(authService))
	router.Post("/auth/login", auth.LoginHandler(authService))
	router.Get("/crypto", auth.AuthMiddleware(crypto.GETCryptosHandler(cryptoService)))
	router.Post("/crypto", auth.AuthMiddleware(crypto.POSTCryptoHandler(cryptoService)))
	router.Get("/crypto/{symbol}", auth.AuthMiddleware(crypto.GETCryptoSymbolHandler(cryptoService)))
	router.Put("/crypto/{symbol}/refresh", auth.AuthMiddleware(crypto.PUTCryptoSymbolRefreshHandler(cryptoService)))
	router.Get("/crypto/{symbol}/history", auth.AuthMiddleware(crypto.GETCryptoHistoryHandler(cryptoService)))
	router.Get("/crypto/{symbol}/stats", auth.AuthMiddleware(crypto.GETCryptoStatsHandler(cryptoService)))
	router.Delete("/crypto/{symbol}", auth.AuthMiddleware(crypto.DELETECryptoSymbolHandler(cryptoService)))

	log.Println("Starting server on http://:8080")
	err = http.ListenAndServe(":8080", router)
	if err != nil {
		log.Fatalln(err)
	}
}