package handlers

import (
	"fmt"
	"net/http"
	"text/template"

	clib "github.com/Astenna/Nubes/movie_review/client_lib"
)

var currentlyLoggedIn any

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		var account clib.AccountStub
		account.Password = r.PostFormValue("Password")
		account.Email = r.PostFormValue("Email")

		initializedAccount, _ := clib.LoadAccount(account.Email)
		if initializedAccount == nil {
			fmt.Fprintf(w, "Account with %s not found", account.Email)
			return
		}
		loggedIn, err := initializedAccount.VerifyPassword(account.Password)

		if !loggedIn && err != nil {
			fmt.Fprintf(w, "Incorrect password")
		} else if initializedAccount != nil {
			fmt.Fprintf(w, "Successfully logged in")
			currentlyLoggedIn = initializedAccount
		}
		return
	}
	t, _ := template.ParseFiles("templates//login.html")
	t.Execute(w, nil)
}
