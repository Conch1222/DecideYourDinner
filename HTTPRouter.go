package main

import (
	"log"
	"net/http"
)

func setHTTPRouter() {
	http.HandleFunc("/", hello)
	http.HandleFunc("/login", login)
	http.HandleFunc("/mainPage", mainPage)
	http.HandleFunc("/mainPage/result", result)
	err := http.ListenAndServe(":9090", nil) //設定監聽的埠
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
