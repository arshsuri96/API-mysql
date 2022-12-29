package main

import (
	"database/sql"
	"fmt"
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

	l := library{
		dbHost: dbHost,
		dbPass: dbPass,
		dbName: dbName,
	}
	r := mux.NewRouter()
	r.HandleFunc(apiPath, l.getBooks).Methods("GET")
	http.ListenAndServe(":8080", r)
}

func (l library) getBooks(http.ResponseWriter, *http.Request) {
	//open connection
	db := l.openConnection()
	//read the books

	books := []Book{}
	rows, err := db.Query("select * from library")
	if err != nil {
		log.Fatalf("querying the book table %s\n", err.Error())
	}
	for rows.Next() {
		var id, name, isbn string
		err := rows.Scan(&id, &name, &isbn)
		if err != nil {
			log.Fatalf("error while scanning the row %s\n", err.Error())
		}
		aBook := Book{
			Id:   id,
			Name: name,
			Isbn: isbn,
		}
		books = append(books, aBook)
	}

	//close connection
}

func (l library) openConnection() *sql.DB {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@(%s)/%s", "root", l.dbHost, l.dbPass, l.dbName))
	if err != nil {
		log.Fatalf("opening the connection to DB %s\n", err.Error())
	}
	return db
}
