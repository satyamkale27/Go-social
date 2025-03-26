package main

import (
	"log"
)

func main() {
	cfg := config{
		addr: ":8080",
	}

	app := &application{
		config: cfg,
	}

	mux := app.mount()
	log.Fatal(app.run(mux))

	/*
		When you call app.run(),
		Go automatically passes the app pointer to the run method.
		This is why you don't need to call run(app) explicitly.
	*/
}
