package main

import (
	"log"

	"github.com/ecetinerdem/forseerv2/internal/db"
	"github.com/ecetinerdem/forseerv2/internal/env"
	"github.com/ecetinerdem/forseerv2/internal/store"
)

func main() {

	cfg := config{
		addr: env.GetString("ADDR", ":8080"),
		db: dbConfig{
			addr:        env.GetString("DB_ADDR", "postgres://postgres:postgres@localhost:5432/forseer?sslmode=disable"),
			maxOpenConn: env.GetInt("DB_MAX_OPEN_CONN", 30),
			maxIdleConn: env.GetInt("DB_MAX_IDLE_CONN", 30),
			maxIdleTime: env.GetString("DB_MAX_IDLE_TIME", "15m"),
		},
		env:     env.GetString("ENV", "development"),
		version: env.GetString("VERSION", "1.0.0"),
	}

	db, err := db.New(cfg.db.addr, cfg.db.maxOpenConn, cfg.db.maxIdleConn, cfg.db.maxIdleTime)

	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()
	log.Println("database connection established")

	store := store.NewStorage(db)

	app := &application{
		config: cfg,
		store:  store,
	}

	mux := app.mount()

	log.Fatal(app.run(mux))
}
