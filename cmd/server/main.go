package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/sawadashota/di-sample/driver"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stdout, err)
		os.Exit(1)
	}
}

func run() error {
	d := driver.NewDefaultDriver()
	d.CallRegistry()

	router := mux.NewRouter()
	d.Registry().RegisterRoutes(router)

	d.Registry().Logger().Infoln("starting server 127.0.0.1:8080")
	return http.ListenAndServe(":8080", router)
}
