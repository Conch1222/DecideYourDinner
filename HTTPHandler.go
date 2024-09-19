package main

import (
	"GoWeb/Error"
	"GoWeb/Type"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"sync"
)

var resultsNearByMap []*Type.NearBy
var rankStoresList []Type.NearByResult
var storeCount = 0
var currentUser = new(Type.User)
var pageMutex sync.Mutex

func hello(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		return
	}
	fmt.Println(r.Form)
	fmt.Fprintf(w, "Hello, %s! It's for test.", r.Form.Get("name"))
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

		dbConn := connectDB()
		_, err := dbConn.saveUserInfo(userName, password)
		if err != nil {
			fmt.Println(err)
			fmt.Fprintf(w, "Sign up fail: %s", err)
			return
		}
		http.Redirect(w, r, "/login", http.StatusFound)
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

		dbConn := connectDB()
		userName := r.Form["username"][0]
		password := r.Form["password"][0]
		user := dbConn.getUserByUserName(userName)

		if err := validateLogin(user, password); err != nil {
			fmt.Println(err)
			fmt.Fprintf(w, "Login fail: %s", err)
			return
		}
		currentUser = user
		processAuth(w, r, currentUser.ID)

		fmt.Printf("login success, user id : %d ", currentUser.ID)
		http.Redirect(w, r, "/mainPage", http.StatusFound)
	}
}

func logout(w http.ResponseWriter, r *http.Request) {
	clearSession(w, r)
	http.Redirect(w, r, "/login", http.StatusFound)
}

func mainPage(w http.ResponseWriter, r *http.Request) {
	fmt.Println("method: ", r.Method)
	if err := r.ParseForm(); err != nil {
		return
	}

	if r.Method == "GET" {
		if !isAuth(w, r, currentUser.ID) {
			fmt.Fprintf(w, Error.InvalidUser)
		}

		type data struct {
			UserName string
		}

		t := template.Must(template.ParseFiles("File/MainPage.html"))
		if err := t.Execute(w, data{
			UserName: currentUser.UserName,
		}); err != nil {
			fmt.Println(err)
			return
		}
	} else if r.Method == "POST" {
		pageMutex.Lock()
		defer pageMutex.Unlock()

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

		resultsNearBy := make([]*Type.NearBy, 0)
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

		http.Redirect(w, r, "/mainPage/result", http.StatusFound)
	}
}

func resultPage(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" && storeCount < len(rankStoresList) {
		if !isAuth(w, r, currentUser.ID) {
			fmt.Fprintf(w, Error.InvalidUser)
		}
		store := rankStoresList[storeCount]
		log.Println(store)

		type data struct {
			StoreName    string
			StoreAddress string
			StoreRating  string
			StoreMapLink string
		}

		resultData := new(Type.QueryData)
		resultData.Init(store.Name, convertAddress(store), store.Rating,
			Type.Url_GoogleSearch+store.Name+Type.Url_GoogleSearch_PlaceIdParm+store.PlaceId)

		t := template.Must(template.ParseFiles("File/ResultPage.html"))
		if err := t.Execute(w, data{
			StoreName:    resultData.StoreName,
			StoreAddress: resultData.StoreAddress,
			StoreRating:  strconv.FormatFloat(resultData.StoreRating, 'f', -1, 64),
			StoreMapLink: resultData.StoreMapLink,
		}); err != nil {
			fmt.Println(err)
			return
		}

		dbConn := connectDB()
		_, err := dbConn.saveQueryRecord(currentUser.ID, resultData.StoreName, resultData.StoreAddress,
			resultData.StoreRating, resultData.StoreMapLink)
		if err != nil {
			fmt.Fprintf(w, err.Error())
			return
		}

		storeCount++
	} else {
		fmt.Fprintf(w, "Recommendation has ended!")
	}
}

func historicalRecords(w http.ResponseWriter, r *http.Request) {
	fmt.Println("method:", r.Method)
	if err := r.ParseForm(); err != nil {
		return
	}
	if r.Method == "GET" {
		if !isAuth(w, r, currentUser.ID) {
			fmt.Fprintf(w, Error.InvalidUser)
		}

		dbConn := connectDB()
		records, err := dbConn.getQueryRecord(currentUser.ID, 20)
		if err != nil {
			fmt.Fprintf(w, err.Error())
			return
		}

		if len(records) == 0 {
			fmt.Fprintf(w, Error.HistoricalRecordEmpty)
			return
		}

		log.Println(records)
		t, _ := template.ParseFiles("File/HistoricalRecords.html")
		log.Println(t.Execute(w, records))
	}
}
