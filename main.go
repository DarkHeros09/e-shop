package main

import (
	"database/sql"
	"log"

	"github.com/DarkHeros09/e-shop/v2/api"
	db "github.com/DarkHeros09/e-shop/v2/db/sqlc"
	"github.com/DarkHeros09/e-shop/v2/util"

	_ "github.com/lib/pq"
)

func main() {
	config, err := util.LoadConfig(".") // we use . because app.env is on the same level with main.go
	if err != nil {
		log.Fatal("cannot load config:", err)
	}
	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db", err)
	}

	store := db.NewStore(conn)
	server := api.NewServer(store)

	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal("cannot start server:", err)
	}
}