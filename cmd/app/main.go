package main

import (
	//"context"
	"fmt"
	"log"
	"net/http"
	"os"

	jwt "github.com/form3tech-oss/jwt-go"
	"github.com/gorilla/mux"
	"go-template/internal/handler"
	"go-template/internal/middleware"
)

var secretKey []byte = []byte(os.Getenv("JWT_SECRET"))

func main() {
	r := Router()
	http.Handle("/", r)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Listening on port %s", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}

func Router() *mux.Router {
	r := mux.NewRouter()
	m := middleware.NewJWtMiddleware(middleware.Options{
		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
			return secretKey, nil
		},
		SigningMethod: jwt.SigningMethodHS256,
	})
	r.Handle("/health", handler.ServerStatusHandler).Methods("GET")
	r.Handle("/works", m.Handler(handler.WorksListHandler)).Methods("GET")
	return r
}
