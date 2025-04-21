package main

import (
	"github.com/satyamkale27/Go-social.git/internal/db"
	"github.com/satyamkale27/Go-social.git/internal/env"
	"github.com/satyamkale27/Go-social.git/internal/store"
	"log"
)

// this main is for running the sed script that inserts dummy data

func main() {

	addr := env.GetString("DB_ADDR", "postgres://admin:adminpassword@localhost/socialnetwork?sslmode=disable")
	conn, err := db.New(addr, 3, 3, "15m")
	if err != nil {
		log.Fatal(err)
	}
	storee := store.NewStorage(conn)
	db.Seed(storee)
}
