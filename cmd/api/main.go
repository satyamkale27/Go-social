package main

import (
	db2 "github.com/satyamkale27/Go-social.git/internal/db"
	"github.com/satyamkale27/Go-social.git/internal/env"
	mailer2 "github.com/satyamkale27/Go-social.git/internal/mailer"
	store2 "github.com/satyamkale27/Go-social.git/internal/store"
	"go.uber.org/zap"
	"os"
	"time"
)

const version = "0.0.1"

func main() {
	cfg := config{
		addr:        env.GetString("ADDR", ":8080"),
		apiURL:      env.GetString("API_URL", "http://localhost:8080"),
		frontendUrl: env.GetString("FRONTEND_URL", "http://localhost:4000"),
		db: dbConfig{
			addr:         env.GetString("DB_ADDR", "postgres://admin:adminpassword@localhost/socialnetwork?sslmode=disable"),
			maxOpenConns: env.GetInt("MAX_OPEN_CONNS", 30),
			maxIdleConns: env.GetInt("MAX_IDLE_CONNS", 30),
			maxIdleTime:  env.GetString("MAX_IDLE_TIME", "15m"),
		},
		env: env.GetString("ENV", "development"),
		mail: mailConfig{
			fromEmail: env.GetString("FROM_EMAIL", ""),
			exp:       time.Hour * 24 * 3, // 3 days
			sendGrid: sendgridConfig{
				apiKey: env.GetString("SENDGRID_API_KEY", ""),
			},
		},
	}

	logger := zap.Must(zap.NewDevelopment()).Sugar()
	defer logger.Sync()

	db, err := db2.New(
		cfg.db.addr,
		cfg.db.maxOpenConns,
		cfg.db.maxIdleConns,
		cfg.db.maxIdleTime,
	)
	if err != nil {
		logger.Fatal(err)
	}

	defer db.Close()
	logger.Info("Database connection pool established")

	store := store2.NewStorage(db)

	mailer := mailer2.NewSendgrid(cfg.mail.sendGrid.apiKey, cfg.mail.fromEmail)

	app := &application{
		config: cfg,
		store:  store,
		logger: logger,
		mailer: mailer,
	}
	os.LookupEnv("PATH")

	mux := app.mount()
	logger.Fatal(app.run(mux))

	/*
		When you call app.run(),
		Go automatically passes the app pointer to the run method.
		This is why you don't need to call run(app) explicitly.
	*/
}
