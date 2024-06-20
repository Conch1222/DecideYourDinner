package main

import (
	"database/sql"
	"fmt"
)

const (
	USERNAME = "admin"
	PASSWORD = "admin"
	NETWORK  = "tcp"
	SERVER   = "127.0.0.1"
	PORT     = 3306
	DATABASE = "web"
)

var DB *sql.DB

func connectDB() *sql.DB {
	if DB == nil {
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
		DB = db
	}
	return DB
}
