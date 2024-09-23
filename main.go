package main

import (
	_ "github.com/go-sql-driver/mysql"
)

func main() {

	err := setKeyAndSetCookieStore()
	if err != nil {
		panic(err)
	}

	setHTTPRouter()

}
