package main

import (
	"log"

	"github.com/ecetinerdem/forseerv2/internal/db"
	"github.com/ecetinerdem/forseerv2/internal/env"
	"github.com/ecetinerdem/forseerv2/internal/store"
)

func main() {

	addr := env.GetString("DB_ADDR", "postgres://postgres:postgres@localhost:5432/forseer?sslmode=disable")
	log.Println("Connecting to:", addr)

	conn, err := db.New(addr, 3, 3, "15m")
	if err != nil {
		log.Fatal(err)
	}

	defer conn.Close()
	store := store.NewStorage(conn)

	db.Seed(store)
}
