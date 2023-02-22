package handlers

import (
	"fmt"
	"net/http"
	"text/template"

	"github.com/Astenna/Nubes/evaluation/movie_review_baseline/client_project/models"
)

var currentlyLoggedIn string

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		var account models.LoginParams
		account.Password = r.PostFormValue("Password")
		account.Email = r.PostFormValue("Email")

		correctPassword, err := invokeLambdaToGetSingleItem[bool](account, "login")
		if err != nil {
			fmt.Fprintf(w, "Error occurred while logging in %s", err.Error())
			return
		}

		if !*correctPassword && err != nil {
			fmt.Fprintf(w, "Incorrect password")
			return
		}

		fmt.Fprintf(w, "Successfully logged in")
		currentlyLoggedIn = account.Email
		return
	}
	t, _ := template.ParseFiles("templates//login.html")
	t.Execute(w, nil)
}
