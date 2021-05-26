package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/codegangsta/negroni"
	"github.com/davisbento/go-api/core/articles"
	"github.com/davisbento/go-api/core/jwtManager"
	"github.com/davisbento/go-api/core/users"
	"github.com/davisbento/go-api/database"
	"github.com/davisbento/go-api/web/handlers"
	"github.com/gorilla/mux"
)

func main() {
	db := database.Connect()
	aService := articles.NewService(db)
	uService := users.NewService(db)
	r := mux.NewRouter()

	n := negroni.New(
		negroni.NewLogger(),
		jwtManager.JwtAuth(),
	)

	//handlers
	handlers.MakeArticlesHandler(r, n, aService)
	handlers.MakeUsersHandlers(r, n, uService)

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
