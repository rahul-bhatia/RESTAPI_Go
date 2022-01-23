package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"database/sql"
	_ "github.com/lib/pq"

	"github.com/gorilla/mux"
)
const(
	host = "localhost"
	port = 5432
	user = "postgres"
	password = "pass"
	dbname = "first_db"
)

var Db *sql.DB 
var err error

// Model for books
type Book struct {
	ID     string  `json:"id"`
	Isbn   string  `json:"isbn"`
	Title  string  `json:"title"`
	Author string  `json:"author"`
}
type JsonResponse struct {
    Type    string `json:"type"`
    Data    []Book `json:"data"`
    Message string `json:"message"`
}

type Author struct {
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
}

//Init books : slice
var books []Book

//middleware handlers
func getBooks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	//json.NewEncoder(w).Encode(books)
	row,err := Db.Query(`SELECT * from "Books"`)
	if err != nil{
		log.Fatal(fmt.Println("Error",err))
   }
	var book []Book
	for row.Next(){
		var id     string
		var isbn   string
		var title  string
		var author string

		err = row.Scan(&id,&isbn,&title,&author)
		book =append(book,Book{ID:id,Isbn:isbn,Title:title,Author:author})
	}
	var response = JsonResponse{Type: "success", Data: book}
	json.NewEncoder(w).Encode(response)
}

func getBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for _, item := range books {
		fmt.Println(item.ID, params)
		if item.ID == params["id"] {
			json.NewEncoder(w).Encode(item)
			return
		}
	}
	json.NewEncoder(w).Encode(&Book{})
}

func createBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var book Book
	//fmt.Println("Body:", r.Body)
	_ = json.NewDecoder(r.Body).Decode(&book)
	fmt.Println("Body:", book)
	book.ID = strconv.Itoa(rand.Intn(1000000))
	books = append(books, book)

	//Trying to add in db
	insertstmt := `insert into "Books" ("ID","Isbn","Title","Author") values($1,$2,$3,$4)`
	_, e := Db.Exec(insertstmt,book.ID,book.Isbn,book.Title,"carol")
	if e != nil{
	 	log.Fatal(fmt.Println("Error",e))
	}
	//json.NewEncoder(w).Encode(book)
}

func updateBooks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	var book Book
	_ = json.NewDecoder(r.Body).Decode(&book)
	book.ID = strconv.Itoa(rand.Intn(1000000))
	books = append(books, book)
	// for index, item := range books {
	// 	if params["id"] == item.ID {
	// 		books = append(books[:index], books[index+1:]...)
	// 		var book Book
	// 		_ = json.NewDecoder(r.Body).Decode(&book)
	// 		book.ID = strconv.Itoa(rand.Intn(10000))
	// 		books = append(books, book)
	// 		json.NewEncoder(w).Encode(book)
	// 		return
	// 	}
	// }
	insertstmt := `Update "Books" SET "ID"= $1,"Isbn"=$2 ,"Title"=$3 ,"Author"= $4 WHERE "ID"=$5`
	_, e := Db.Exec(insertstmt,book.ID,book.Isbn,book.Title,"carol",params["id"])
	if e != nil{
		log.Fatal(fmt.Println("Error",e))
   	}
	//json.NewEncoder(w).Encode(books)

}

func deleteBooks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for index, item := range books {
		if params["id"] == item.ID {
			books = append(books[:index], books[index+1:]...)
			json.NewEncoder(w).Encode(books)
		}
	}
	json.NewEncoder(w).Encode(books)
}

// Connecting Go project with LocalDB.
func DB(){
	psqlconn := fmt.Sprintf("host= %s port = %d user = %s password = %s dbname = %s sslmode=disable",host,port,user,password,dbname)

	Db, err = sql.Open("postgres",psqlconn)
	if err != nil{
		CheckError(err)
	}
	//defer Db.Close()
}

func CheckError(err error){
	log.Fatal(fmt.Println("Error:",err))
}

func main() {
	// Init Router
	r := mux.NewRouter()

	//Setup the Database connection
	DB()

	// Mock Data : No more needed.
	books = append(books, Book{ID: "1", Isbn: "3928", Title: "Book one", Author: "John"})
	books = append(books, Book{ID: "2", Isbn: "23128", Title: "Book two", Author:"Steve"})
	
	//Router handlers
	r.HandleFunc("/api/books", getBooks).Methods("GET")
	r.HandleFunc("/api/book/{id}", getBook).Methods("GET")
	r.HandleFunc("/api/books", createBook).Methods("POST")
	r.HandleFunc("/api/books/{id}", updateBooks).Methods("PUT")
	r.HandleFunc("/api/books/{id}", deleteBooks).Methods("DELETE")

	//Run server
	log.Fatal(http.ListenAndServe(":8080", r))
}
