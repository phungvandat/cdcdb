package main

import (
	"context"
	"os"
	"strconv"

	"github.com/jackc/pgx"
)

// Default config
const (
	dfHost   = "localhost"
	dfPort   = "5432"
	dfDBName = "repdb"
	dfUser   = "postgres"
	dfPass   = "repdbpass"
)

func initConn(ctx context.Context) {
	once.Do(func() {
		var (
			host   = os.Getenv("PG_HOST")
			port   = os.Getenv("PG_PORT")
			dbName = os.Getenv("PG_DB_NAME")
			user   = os.Getenv("PG_USER")
			pass   = os.Getenv("PG_PASS")
			err    error
		)
		if host == "" {
			host = dfHost
		}
		if port == "" {
			port = dfPort
		}
		if dbName == "" {
			dbName = dfDBName
		}
		if user == "" {
			user = dfUser
		}
		if pass == "" {
			pass = dfPass
		}

		var (
			portNum, _ = strconv.Atoi(port)
			cfg        = pgx.ConnConfig{
				Host:     host,
				Port:     uint16(portNum),
				Database: dbName,
				User:     user,
				Password: pass,
			}
		)

		repConn, err = pgx.ReplicationConnect(cfg)
		if err != nil {
			panic(err)
		}

		dbConn, err = pgx.Connect(cfg)
		if err != nil {
			panic(err)
		}
		err = dbConn.Ping(ctx)
		if err != nil {
			panic(err)
		}

		log.Infoln("db connected")
	})
}
