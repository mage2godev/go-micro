package main

import (
	"authentication/data"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
)

const webPort = "1226"

var counts int64

type Config struct {
	DB     *sql.DB
	Models data.Models
}

func main() {
	log.Println("Starting authentication service")

	// connect to db
	conn := connectDB()
	if conn == nil {
		log.Panic("Could not connect to database")
	}

	app := Config{
		DB:     conn,
		Models: data.New(conn),
	}

	srv := http.Server{
		Addr:    fmt.Sprintf(":%s" + webPort),
		Handler: app.routes(),
	}
	err := srv.ListenAndServe()

	if err != nil {
		log.Panic(err)
	}
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pdx", dsn)
	if err != nil {
		return nil, err
	}
	
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return db, nil
}

func connectDB() *sql.DB {
	dsn := os.Getenv("DSN")

	for {
		connection, err := openDB(dsn)
		if err != nil {
			log.Println("Postgress not yet ready Error connecting to database")
			counts++
		} else {
			fmt.Println("Postgress successfully connected to database")
			return connection
		}

		if counts > 10 {
			log.Println("Too many retries")
			log.Println(err)
		}

		log.Println("Retrying in 2 seconds...")
		time.Sleep(2 * time.Second)
		continue
	}
}
