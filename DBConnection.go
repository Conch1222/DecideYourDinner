package main

import (
	"database/sql"
	"fmt"
	"sync"
)

const (
	USERNAME = "admin"
	PASSWORD = "admin"
	NETWORK  = "tcp"
	SERVER   = "127.0.0.1"
	PORT     = 3306
	DATABASE = "web"
)

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
		fmt.Println("Connecting to " + DATABASE)
		conn := fmt.Sprintf("%s:%s@%s(%s:%d)/%s", USERNAME, PASSWORD, NETWORK, SERVER, PORT, DATABASE)
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
