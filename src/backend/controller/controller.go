package controller

import (
	"fmt"
	"net/http"
	"main/service"
	"github.com/julienschmidt/httprouter"
)

func SetupRoutes(router *httprouter.Router) {
	router.GET("/helloworld", func (writer http.ResponseWriter, req *http.Request, params httprouter.Params) {
		fmt.Fprintf(writer, "Hello World!")
	})

	router.POST("/register", service.HandleRegister)
}