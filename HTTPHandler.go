package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
)

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

		noError, err := validateLogin(user, password)
		if !noError {
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
		noError, err, optionMap := validateMainPageInput(r)
		if !noError {
			fmt.Println(err)
			fmt.Fprintf(w, "Invalid input: %s", err)
			return
		}

		fmt.Println(optionMap)
		fmt.Fprintln(w, "success!")

		client := getClient()
		loc, err := client.getUserLocation(w)
		if err != nil {
			fmt.Println(err)
			fmt.Fprintf(w, err.Error())
			return
		}

		for k, _ := range optionMap {
			client.getUserNearBy(w, loc, k)
		}

		//http.Redirect(w, r, "/mainPage?username="+userName, http.StatusSeeOther)
	}
}

func result(w http.ResponseWriter, r *http.Request) {

}
