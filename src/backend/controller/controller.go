package controller

import (
	"database/sql"
	"main/service"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"golang.org/x/net/context"
)

func SetupRoutes(router *httprouter.Router, dbConn *sql.DB) {
	router.POST("/register", applyDbMiddleware(dbConn, service.HandleRegister))
	router.POST("/login", applyDbMiddleware(dbConn, service.HandleLogin))
	router.POST("/todo", applyDbMiddleware(dbConn, service.HandleCreateSimpleTodo))
	router.GET("/todo", applyDbMiddleware(dbConn, service.HandleGetTodos))
	router.DELETE("/todo", applyDbMiddleware(dbConn, service.HandleDeleteSimpleTodo))
}

func applyDbMiddleware(dbConn *sql.DB, next httprouter.Handle) httprouter.Handle {
	return func(writer http.ResponseWriter, req *http.Request, params httprouter.Params) {
		ctx := context.WithValue(req.Context(), "db", dbConn)
		next(writer, req.WithContext(ctx), params)
	}
}
