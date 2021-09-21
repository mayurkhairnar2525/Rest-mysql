package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
)

type BookManagement struct {
	ID   int `json:"id"`
	Name string `json:"name"`
}

var db *sql.DB
var err error

func getBooks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var books []BookManagement
	result, err := db.Query("SELECT id, name from bookmanagement")
	if err != nil {
		panic(err.Error())
	}
	defer result.Close()
	for result.Next() {
		var book BookManagement
		err := result.Scan(&book.ID, &book.Name)
		if err != nil {
			panic(err.Error())
		}
		books = append(books, book)
	}
	json.NewEncoder(w).Encode(books)
}

func createBook(w http.ResponseWriter, r *http.Request) {
	stmt, err := db.Prepare("INSERT INTO bookmanagement(name) VALUES(?)")
	if err != nil {
		panic(err.Error())
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err.Error())
	}
	keyVal := make(map[string]string)
	json.Unmarshal(body, &keyVal)
	name := keyVal["name"]
	_, err = stmt.Exec(name)
	if err != nil {
		panic(err.Error())
	}
	fmt.Fprintf(w, "New book was created")

}

func getBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	result, err := db.Query("SELECT id, name FROM bookmanagement WHERE id = ?", params["id"])
	if err != nil {
		panic(err.Error())
	}
	defer result.Close()
	var book BookManagement
	for result.Next() {
		err := result.Scan(&book.ID, &book.Name)
		if err != nil {
			panic(err.Error())
		}
	}
	json.NewEncoder(w).Encode(book)
}

func updateBook(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	stmt, err := db.Prepare("UPDATE bookmanagement SET name = ? WHERE id = ?")
	if err != nil {
		panic(err.Error())
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err.Error())
	}
	keyVal := make(map[string]string)
	json.Unmarshal(body, &keyVal)
	newName := keyVal["name"]
	_, err = stmt.Exec(newName, params["id"])
	if err != nil {
		panic(err.Error())
	}
	fmt.Fprintf(w, "Book with ID = %s was updated", params["id"])
}

func deleteBook(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	stmt, err := db.Prepare("DELETE FROM bookmanagement WHERE id = ?")
	if err != nil {
		panic(err.Error())
	}
	_, err = stmt.Exec(params["id"])
	if err != nil {
		panic(err.Error())
	}
	fmt.Fprintf(w, "Book with ID = %s was deleted", params["id"])
}

func handler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello Everyone!"))
}

func main() {
	db, err = sql.Open("mysql", "root:12345678@tcp(0.0.0.0:9090)/bookstore")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	router := mux.NewRouter()

	router.HandleFunc("/books", getBooks).Methods("GET")

	router.HandleFunc("/books", createBook).Methods("POST")
	router.HandleFunc("/books/{id}", getBook).Methods("GET")
	router.HandleFunc("/books/{id}", updateBook).Methods("PUT")
	router.HandleFunc("/books/{id}", deleteBook).Methods("DELETE")
	router.HandleFunc("/", handler)
	http.ListenAndServe(":8090", router)
}
