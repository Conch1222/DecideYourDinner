package main

import (
	"log"
	"net/http"
)

func setHTTPRouter() {
	http.HandleFunc("/", hello)
	http.HandleFunc("/signUp", signUp)
	http.HandleFunc("/login", login)
	http.HandleFunc("/logout", logout)

	// protect page
	http.HandleFunc("/mainPage", authMiddleware(mainPage))
	http.HandleFunc("/mainPage/result", authMiddleware(resultPage))
	http.HandleFunc("/mainPage/historicalRecord", authMiddleware(historicalRecords))

	err := http.ListenAndServe(":9090", nil) //set port
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
