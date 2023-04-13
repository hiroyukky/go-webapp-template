package main

import(
        "fmt"
        "log"
        "net/http"
        "os"

        "github.com/gorilla/mux"
        "go-template/internal/handler"
)

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
     r.Handle("/health", handler.ServerStatusHandler).Methods("GET")
     return r
}