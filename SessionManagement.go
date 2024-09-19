package main

import (
	"GoWeb/Error"
	"GoWeb/Type"
	"fmt"
	"github.com/gorilla/sessions"
	"net/http"
)

var CookieStore *sessions.CookieStore

func setKeyAndSetCookieStore() {
	key, err := ReadKey("File/SessionKey.txt", Error.InvalidSessionKey)
	if err != nil {
		panic(err)
	}
	CookieStore = sessions.NewCookieStore([]byte(key))
	CookieStore.MaxAge(600)
}

func authMiddleware(nextPage http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, _ := CookieStore.Get(r, Type.SESSION_NAME)
		if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}

		nextPage.ServeHTTP(w, r)
	}
}

func isAuth(w http.ResponseWriter, r *http.Request, userID int64) bool {
	session, _ := CookieStore.Get(r, Type.SESSION_NAME)
	return userID > 0 && session.Values["authenticated"] == true
}

func processAuth(w http.ResponseWriter, r *http.Request, userID int64) {
	session, _ := CookieStore.Get(r, Type.SESSION_NAME)
	session.Values["authenticated"] = true
	session.Values["user_id"] = userID
	session.Options.MaxAge = 600

	if err := session.Save(r, w); err != nil {
		fmt.Fprintf(w, "Internal error: %s", err)
		return
	}
}

func clearSession(w http.ResponseWriter, r *http.Request) {
	session, _ := CookieStore.Get(r, Type.SESSION_NAME)
	session.Values["authenticated"] = false
	session.Values["user_id"] = 0
	session.Options.MaxAge = -1

	if err := session.Save(r, w); err != nil {
		fmt.Fprintf(w, "Internal error: %s", err)
		return
	}
}
