package main

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"net"
	"net/http"
	"rest-api-tutorial/internal/user"
	"rest-api-tutorial/pkg/logging"
	"time"
)

func main() {
	logger := logging.GetLogger()
	logger.Info("create route")
	//log.Println("create router")
	router := httprouter.New()
	//router.GET("/:name", IndexHandler)

	logger.Info("register user handler")
	//log.Println("register user handler")
	handler := user.NewHandler(logger)
	handler.Register(router)

	start(router)
}

func start(router *httprouter.Router) {
	logger := logging.GetLogger()
	logger.Info("start application")
	//log.Println("start application")

	//we use package net becouse we will run app on IP and on socket
	listener, err := net.Listen("tcp", ":1234")

	if err != nil {
		panic(err)
	}

	server := &http.Server{
		Handler:      router,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	logger.Info("server is listening to port 1234")
	//log.Println("server is listening to port 1234")

	logger.Fatal(server.Serve(listener))
	//log.Fatalln(server.Serve(listener))
}
func IndexHandler(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	name := params.ByName("name")
	w.Write([]byte(fmt.Sprintf("Hello %s", name)))
}
