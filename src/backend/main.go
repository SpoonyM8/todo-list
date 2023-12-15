package main

import (
	"database/sql"
	//"encoding/json"
	"fmt"
	"log"
	"main/controller"
	"net/http"

	"github.com/julienschmidt/httprouter"
	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "dev"
	password = "root"
	dbname   = "todo_list_db"
)

func main() {
	var db *sql.DB
	db = connectToPostgres()
	fmt.Println(db)
	router := httprouter.New()
	controller.SetupRoutes(router)
	startHttpServer(router)
}

func connectToPostgres() *sql.DB {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to db! Let's go!")
	return db
}

func handler(writer http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(writer, "Hello, Go Web Server!")
}

func startHttpServer(router *httprouter.Router) {

	fmt.Println("Server listening on port 8080...")
	err := http.ListenAndServe(":8080", router)
	if err != nil {
		fmt.Println("Error starting server: ", err)
	}
}
