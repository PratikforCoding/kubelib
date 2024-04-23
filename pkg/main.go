package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-chi/chi/v5"
)

const (
	API_PATH = "/app/v1/books"
)

type library struct {
	dbHost, dbPass, dbName string
}

type Book struct {
	Id, Name, Isbn string
}
func main() {
	dbHost := os.Getenv("DB_HOST")
	if dbHost == "" {
		dbHost = "localhost:3306"
	}

	dbPass := os.Getenv("DB_PASS")
	if dbPass == "" {
		dbPass = "pratikkotal"
	}

	apiPath := os.Getenv("API_PATH")
	if apiPath == "" {
		apiPath = API_PATH
	}

	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		dbName = "library"
	}

	l := library {
		dbHost: dbHost,
		dbPass: dbPass,
		dbName: dbName,
	}

	r := chi.NewRouter()
	r.Get("/app/v1/books", l.getBooks)
	r.Post("/app/v1/books", l.postBooks)

	server := &http.Server{
		Addr: ":8080",
		Handler: r,
	}

	server.ListenAndServe()
}

func (l *library)getBooks(w http.ResponseWriter, r *http.Request) {
	db := l.openConnection()
	rows, err := db.Query("select * from books")
	if err != nil {
		log.Fatal("Error getting books", err.Error())
	}

	books := []Book{}
	for rows.Next() {
		var id, name, isbn string
		err := rows.Scan(&id, &name, &isbn)
		if err != nil {
			log.Fatalf("While scanning the rows: %s\n", err.Error())
		}
		aBook := Book {
			Id: id,
			Name: name,
			Isbn: isbn,
		}

		books = append(books, aBook)
	}

	json.NewEncoder(w).Encode(books)
	l.closeConnection(db)
}

func (l *library) postBooks(w http.ResponseWriter, r *http.Request) {
	book := Book{}
	json.NewDecoder(r.Body).Decode(&book)

	db := l.openConnection()

	insertQuery, err := db.Prepare("insert into books values (?, ?, ?)")
	if err != nil {
		log.Fatalf("Preparing the db query %s\n", err.Error())
	}
	tx, err := db.Begin()
	if err != nil {
		log.Fatalf("While beginning the transaction %s\n", err.Error())
	}

	_, err = tx.Stmt(insertQuery).Exec(book.Id, book.Name, book.Isbn)
	if err != nil {
		log.Fatalf("Execing the insert command %s\n", err.Error())
	}
	err = tx.Commit()
	if err != nil {
		log.Fatalf("While commit the transaction %s\n", err.Error())
	}
	l.closeConnection(db)
}

func(l *library) openConnection() *sql.DB {
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s", "root", l.dbPass, l.dbHost, l.dbName)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("impossible to create the connection: %s", err)
	}
	return db
}

func (l *library) closeConnection(db *sql.DB) {
	err := db.Close()
	if err != nil {
		log.Fatalf("clossing connection: %s\n", err.Error())
	}
}

