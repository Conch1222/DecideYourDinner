package main

import (
	"GoWeb/File"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

var resultsNearByMap []*File.NearBy
var rankStoresList []File.NearByResult
var storeCount = 0

func hello(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		return
	}
	fmt.Println(r.Form)
	fmt.Fprintf(w, "Hello, %s!", r.Form.Get("name"))
}

func signUp(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		return
	}
	if r.Method == "GET" {
		t, _ := template.ParseFiles("File/SignUp.html")
		log.Println(t.Execute(w, nil))
	} else if r.Method == "POST" {
		userName := r.Form["username"][0]
		password := r.Form["password"][0]
		confirmPassword := r.Form["confirmPassword"][0]

		if err := validateSignUp(userName, password, confirmPassword); err != nil {
			fmt.Println(err)
			fmt.Fprintf(w, "Sign up fail: %s", err)
			return
		}

		db := connectDB()
		_, err := saveUserInfo(db, userName, password)
		if err != nil {
			fmt.Println(err)
			fmt.Fprintf(w, "Sign up fail: %s", err)
			return
		}
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}
}

func login(w http.ResponseWriter, r *http.Request) {
	fmt.Println("method:", r.Method)
	if err := r.ParseForm(); err != nil {
		return
	}
	if r.Method == "GET" {
		t, _ := template.ParseFiles("File/Login.html")
		log.Println(t.Execute(w, nil))
	} else if r.Method == "POST" {

		db := connectDB()
		userName := r.Form["username"][0]
		password := r.Form["password"][0]
		user := getUserByUserName(db, userName)

		if err := validateLogin(user, password); err != nil {
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
		for k, v := range optionMap {
			resultNearBy, err := client.getUserNearBy(w, *loc, k)
			if err != nil {
				fmt.Println(err)
				fmt.Fprintf(w, err.Error())
				return
			}
			resultNearBy.Option = k
			resultNearBy.Weight = v
			resultsNearBy = append(resultsNearBy, resultNearBy)
		}
		resultsNearByMap = resultsNearBy

		if err := validateResultStatus(resultsNearByMap); err != nil {
			fmt.Println(err)
			fmt.Fprintf(w, err.Error())
			return
		}

		eliminateNotOpenResult(resultsNearByMap)
		rankStoresList = rankAllResults(resultsNearByMap)
		storeCount = 0

		http.Redirect(w, r, "/mainPage/result", http.StatusSeeOther)
	}
}

func resultPage(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" && storeCount < len(rankStoresList) {
		store := rankStoresList[storeCount]

		type data struct {
			StoreName    string
			StoreAddress string
			StoreRating  string
			StoreMapLink string
		}

		t := template.Must(template.ParseFiles("File/ResultPage.html"))
		if err := t.Execute(w, data{
			StoreName:    store.Name,
			StoreAddress: convertAddress(store),
			StoreRating:  strconv.FormatFloat(store.Rating, 'f', -1, 64),
			StoreMapLink: File.Url_GoogleSearch + store.Name + File.Url_GoogleSearch_PlaceIdParm + store.PlaceId,
		}); err != nil {
			fmt.Println(err)
			return
		}
		storeCount++
	} else {
		fmt.Fprintf(w, "Recommendation has ended!")
	}
}
