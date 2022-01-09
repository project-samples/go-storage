package main

import (
	"awesomeProject/internal/app"
	"context"
	"fmt"
	"github.com/core-go/config"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

func main() () {
	var conf app.Root
	er1 := config.Load(&conf, "configs/config")
	if er1 != nil {
		panic(er1)
	}

	r := mux.NewRouter()

	er2 := app.Route(r, context.Background(), conf)
	if er2 != nil {
		panic(er2)
	}

	fmt.Println("Start server")
	server := ""
	if conf.Server.Port != nil {
		server = ":" + strconv.FormatInt(*conf.Server.Port, 10)
	}
	if er3 := http.ListenAndServe(server, r); er3 != nil {
		panic(er3)
	}
}