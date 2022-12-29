package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

type library struct {
	dbHost, dbPass, dbName string
}

type Book struct {
	Id, Name, Isbn string
}

const (
	API_PATH = "/apis/v1/books"
)

func main() {
	dbHost := os.Getenv("DB_HOST")
	if dbHost == "" {
		dbHost = "localhost:3306"
	}
	dbPass := os.Getenv("DB_PASS")
	if dbPass == "" {
		dbPass = "arshsuri"
	}
	apiPath := os.Getenv("API_PATH")
	if apiPath == "" {
		apiPath = API_PATH
	}

	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		dbName = "library"
	}

	r := mux.NewRouter()
	r.HandleFunc(apiPath, getBooks).Methods("GET")
	http.ListenAndServe(":8080", r)
}

func getBooks(http.ResponseWriter, *http.Request) {
	log.Printf("getBooks was called")
}
