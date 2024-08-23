package main

import (
	"log"
	"net/http"
)

func setHTTPRouter() {
	http.HandleFunc("/", hello)
	http.HandleFunc("/login", login)
	http.HandleFunc("/mainPage", mainPage)
	http.HandleFunc("/mainPage/result", resultPage)
	err := http.ListenAndServe(":9090", nil) //set port
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
