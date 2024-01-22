package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"codegen/app/db/model"
	"codegen/app/pkg/apimux"
	"codegen/app/pkg/app"
	"codegen/app/pkg/app/space"
	"codegen/app/pkg/app/ticket"
	"codegen/app/pkg/app/user"
	"codegen/app/pkg/authn"

	"github.com/jackc/pgx/v5/pgxpool"
)

//go:generate ./generate.sh

// TODO pass this in from the env
const jwtKey = "secret key"

func main() {
	db, err := pgxpool.New(context.Background(), "")
	// db, err := sql.Open("postgres", "sslmode=disable")
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
	app.RegisterSpaceService(appServer, space.NewService(model, db, jwtKey))
	app.RegisterTicketService(appServer, ticket.NewService(model, db))
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
