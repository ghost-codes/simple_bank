package main

import (
	"database/sql"
	"log"

	"github.com/ghost-codes/simplebank/api"
	db "github.com/ghost-codes/simplebank/db/sqlc"
	"github.com/ghost-codes/simplebank/util"
	_ "github.com/lib/pq"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("unable to load config:", err)
	}
	conn, err := sql.Open(config.DBDriver, config.DBSource())

	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}
	store := db.NewStore(conn)
	server, err := api.NewServer(*store, config)

	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	err = server.Start(config.ServerAddr)
	if err != nil {
		log.Fatal("error occured, Server could not start; Error:", err)
	}
}
