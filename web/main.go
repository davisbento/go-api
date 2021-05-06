package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/codegangsta/negroni"
	"github.com/davisbento/go-api/core/articles"
	"github.com/davisbento/go-api/database"
	"github.com/davisbento/go-api/web/handlers"
	"github.com/gorilla/mux"
)

func main() {
	db := database.Connect()
	service := articles.NewService(db)
	r := mux.NewRouter()

	n := negroni.New(
		negroni.NewLogger(),
	)

	//handlers
	handlers.MakeArticlesHandler(r, n, service)

	http.Handle("/", r)

	srv := &http.Server{
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		Addr:         ":4000",
		Handler:      http.DefaultServeMux,
		ErrorLog:     log.New(os.Stderr, "logger: ", log.Lshortfile),
	}

	fmt.Println("server up")

	err := srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
