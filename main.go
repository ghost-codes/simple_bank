package main

import (
	"database/sql"
	"log"

	"github.com/ghost-codes/simplebank/api"
	db "github.com/ghost-codes/simplebank/db/sqlc"
	_ "github.com/lib/pq"
)

const (
	dbDriver = "postgres"
	dbSource = "postgresql://postgres:password@localhost:5432/bank?sslmode=disable"
	addr     = "0.0.0.0:3000"
)

func main() {
	conn, err := sql.Open(dbDriver, dbSource)

	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}
	store := db.NewStore(conn)
	server := api.NewServer(store)

	err = server.Start(addr)
	if err != nil {
		log.Fatal("error occured, Server could not start; Error:", err)
	}
}
