package main

import (
	db2 "github.com/satyamkale27/Go-social.git/internal/db"
	"github.com/satyamkale27/Go-social.git/internal/env"
	store2 "github.com/satyamkale27/Go-social.git/internal/store"
	"log"
	"os"
)

func main() {
	cfg := config{
		addr: env.GetString("ADDR", ":8080"),
		db: dbConfig{
			addr:         env.GetString("DB_ADDR", ":postgress://user:adminpassword@localhost/social?sslmode=disable"),
			maxOpenConns: env.GetInt("MAX_OPEN_CONNS", 30),
			maxIdleConns: env.GetInt("MAX_IDLE_CONNS", 30),
			maxIdleTime:  env.GetString("MAX_IDLE_TIME", "15min"),
		},
	}

	db, err := db2.New(
		cfg.db.addr,
		cfg.db.maxOpenConns,
		cfg.db.maxIdleConns,
		cfg.db.maxIdleTime,
	)
	if err != nil {
		log.Panic(err)
	}

	store := store2.NewStorage(db)

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
