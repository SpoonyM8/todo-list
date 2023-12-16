package main

import (
	"database/sql"
	"fmt"
	"main/controller"
	"main/db"
	"net/http"

	"github.com/julienschmidt/httprouter"
	_ "github.com/lib/pq"
)

func main() {
	var dbConnection *sql.DB = db.ConnectToPostgres()
	router := httprouter.New()
	controller.SetupRoutes(router, dbConnection)
	startHttpServer(router)
}

func startHttpServer(router *httprouter.Router) {
	fmt.Println("Server listening on port 8080...")
	if err := http.ListenAndServe(":8080", router); err!= nil {
		fmt.Println("Error starting server: ", err)
	}
}
