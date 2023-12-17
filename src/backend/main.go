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
	var dbConn *sql.DB = db.ConnectToPostgres()
	defer dbConn.Close()
	router := httprouter.New()
	controller.SetupRoutes(router, dbConn)
	startHttpServer(router)
}

func startHttpServer(router *httprouter.Router) {
	fmt.Println("Server listening on port 8080...")
	if err := http.ListenAndServe(":8080", router); err != nil {
		fmt.Println("Error starting server: ", err)
	}
}
