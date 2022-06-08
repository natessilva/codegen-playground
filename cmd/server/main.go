package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"codegen/app/db/model"
	"codegen/app/pkg/apimux"
	"codegen/app/pkg/app"
	"codegen/app/pkg/app/ticket"
	"codegen/app/pkg/app/user"
	"codegen/app/pkg/app/workspace"
	"codegen/app/pkg/authn"

	_ "github.com/lib/pq"
)

//go:generate ./generate.sh

// TODO pass this in from the env
const jwtKey = "secret key"

func main() {
	db, err := sql.Open("postgres", "sslmode=disable")
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer db.Close()

	mux := http.NewServeMux()
	model := model.New(db)

	authnServer := apimux.NewServer()
	authn.RegisterAuthNService(authnServer, authn.NewService(model, db, jwtKey))
	mux.Handle("/authn/", http.StripPrefix("/authn", authnServer))

	appServer := apimux.NewServer()
	app.RegisterUserService(appServer, user.NewService(model, db))
	app.RegisterTicketService(appServer, ticket.NewService(model, db))
	app.RegisterWorkspaceService(appServer, workspace.NewService(model, db, jwtKey))
	mux.Handle("/app/", authn.Handle(model, jwtKey, http.StripPrefix("/app", appServer)))

	srv := &http.Server{
		Addr:    ":" + os.Getenv("SERVER_PORT"),
		Handler: apimux.HandleCors(mux),
	}
	go func() {
		log.Fatal(srv.ListenAndServe())
	}()
	log.Printf("Listening on port %s", os.Getenv(("SERVER_PORT")))

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	srv.Shutdown(ctx)
}
