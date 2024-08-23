package main

import (
	"GoWeb/File"
	"fmt"
	"html/template"
	"log"
	"net/http"
)

var resultsNearByMap []*File.NearBy
var optionWeightMap *map[string]float64

func hello(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		return
	}
	fmt.Println(r.Form)
	fmt.Fprintf(w, "Hello, %s!", r.Form.Get("name"))
}

func login(w http.ResponseWriter, r *http.Request) {
	fmt.Println("method:", r.Method)
	if err := r.ParseForm(); err != nil {
		return
	}
	if r.Method == "GET" {
		t, _ := template.ParseFiles("File/Login.gtpl")
		log.Println(t.Execute(w, nil))
	} else if r.Method == "POST" {

		fmt.Println("username:", r.Form["username"])
		fmt.Println("password:", r.Form["password"])

		db := connectDB()
		userName := r.Form["username"][0]
		password := r.Form["password"][0]
		user := getUserByUserName(db, userName)

		err := validateLogin(user, password)
		if err != nil {
			fmt.Println(err)
			fmt.Fprintf(w, "Login fail: %s", err)
			return
		}
		fmt.Println("login success!")
		http.Redirect(w, r, "/mainPage?username="+userName, http.StatusSeeOther)
	}
}

func mainPage(w http.ResponseWriter, r *http.Request) {
	fmt.Println("method: ", r.Method)
	if err := r.ParseForm(); err != nil {
		return
	}

	if r.Method == "GET" {
		type data struct {
			UserName string
		}
		t := template.Must(template.ParseFiles("File/MainPage.html"))
		if err := t.Execute(w, data{
			UserName: r.URL.Query().Get("username"),
		}); err != nil {
			fmt.Println(err)
			return
		}
	} else if r.Method == "POST" {
		fmt.Println("user input: ", r.Form)
		optionMap, err := validateMainPageInput(r)
		if err != nil {
			fmt.Println(err)
			fmt.Fprintf(w, "Invalid input: %s", err)
			return
		}

		fmt.Println(optionMap)

		client := getClient()
		loc, err := client.getUserLocation(w)
		if err != nil {
			fmt.Println(err)
			fmt.Fprintf(w, err.Error())
			return
		}

		resultsNearBy := make([]*File.NearBy, 0)
		for k, _ := range optionMap {
			resultNearBy, err := client.getUserNearBy(w, *loc, k)
			if err != nil {
				fmt.Println(err)
				fmt.Fprintf(w, err.Error())
				return
			}
			resultsNearBy = append(resultsNearBy, resultNearBy)
		}
		resultsNearByMap = resultsNearBy
		optionWeightMap = &optionMap
		http.Redirect(w, r, "/mainPage/result", http.StatusSeeOther)
	}
}

func resultPage(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		err := validateResultStatus(resultsNearByMap)
		if err != nil {
			fmt.Println(err)
			fmt.Fprintf(w, err.Error())
			return
		}

		eliminateNotOpenResult(resultsNearByMap)
		fmt.Fprintln(w, "success!")
	}
}
