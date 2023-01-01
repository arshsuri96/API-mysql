package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

const (
	API_PATH = "/apis/v1/books"
)

type Book struct {
	Id, Name, Isbn string
}

type library struct {
	dbHost, dbPass, dbName string
}

func main() {
	// DB_HOST is of form host:port
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
	r.HandleFunc(apiPath, l.getBooks).Methods(http.MethodGet)
	r.HandleFunc(apiPath, l.postBooks).Methods(http.MethodPost)
	http.ListenAndServe(":8080", r)
}

func (l library) postBooks(w http.ResponseWriter, r *http.Request) {
	//read a request into an instance of book
	book := Book{}
	json.NewDecoder(r.Body).Decode(&book)
	//open connection
	db := l.openConnection()

	insertQuery, err := db.Prepare("insert into books values(?,?,?)")
	if err != nil {
		log.Fatalf("preparing the db query %s\n", err.Error())
	}
	tx, err := db.Begin()
	if err != nil {
		log.Fatalf("")
	}
	_, err = tx.Stmt(insertQuery).Exec(book.Id, book.Name, book.Isbn)
	if err != nil {
		log.Fatalf("executing the insert command %s\n", err.Error())
	}

	err = tx.Commit()
	if err != nil {
		log.Fatalf("Error while commiting the tx %s\n", err.Error())
	}
	//close the request
	l.closeConnection(db)

}

func (l library) getBooks(w http.ResponseWriter, r *http.Request) {
	books := []Book{}
	db := l.openConnection()
	rows, err := db.Query("select * from books")
	if err != nil {
		log.Fatalf("Querying from the table %s", err.Error())
	}
	for rows.Next() {
		var id, name, isbn string
		err := rows.Scan(&id, &name, &isbn)
		if err != nil {
			log.Fatalf("while scanning rows %s", err.Error())
		}
		abook := Book{
			Id:   id,
			Name: name,
			Isbn: isbn,
		}
		books = append(books, abook)
	}
	json.NewEncoder(w).Encode(books)
	l.closeConnection(db)
}

func (l library) openConnection() *sql.DB {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@(%s)/%s", "root", l.dbPass, l.dbHost, l.dbName))
	if err != nil {
		log.Fatalf("opening the connection to the database %s\n", err.Error())
	}
	return db
}

func (l library) closeConnection(db *sql.DB) {
	err := db.Close()
	if err != nil {
		log.Fatalf("closing connection %s\n", err.Error())
	}
}
