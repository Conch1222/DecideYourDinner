package main

import (
	"database/sql"
	"fmt"
	"os"
	"sync"
)

const NETWORK = "tcp"

type DBConnection struct {
	db *sql.DB
}

var DBConn *DBConnection
var onceDBConn sync.Once

func connectDB() *DBConnection {
	onceDBConn.Do(func() {
		db := initDB()
		DBConn = db
	})
	return DBConn
}

func initDB() *DBConnection {
	if DBConn == nil {
		username := os.Getenv("DB_USER")
		password := os.Getenv("DB_PASSWORD")
		server := os.Getenv("DB_HOST")
		port := os.Getenv("DB_PORT")
		database := os.Getenv("DB_NAME")

		fmt.Println("Connecting to " + database)
		conn := fmt.Sprintf("%s:%s@%s(%s:%s)/%s", username, password, NETWORK, server, port, database)
		db, err := sql.Open("mysql", conn)
		if err != nil {
			fmt.Println("Open Mysql error: ", err)
			return nil
		}
		if err := db.Ping(); err != nil {
			fmt.Println("database connect error: ", err.Error())
			return nil
		}

		return &DBConnection{db: db}
	}
	return DBConn
}
