package main

import (
	"github.com/satyamkale27/Go-social.git/internal/env"
	store2 "github.com/satyamkale27/Go-social.git/internal/store"
	"log"
	"os"
)

func main() {
	cfg := config{
		addr: env.GetString("ADDR", ":8080"),
	}

	store := store2.NewStorage(nil)

	app := &application{
		config: cfg,
		store:  store,
	}
	os.LookupEnv("PATH")

	mux := app.mount()
	log.Fatal(app.run(mux))

	/*
		When you call app.run(),
		Go automatically passes the app pointer to the run method.
		This is why you don't need to call run(app) explicitly.
	*/
}
