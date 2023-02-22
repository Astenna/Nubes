package handlers

import (
	"fmt"
	"net/http"
	"text/template"

	"github.com/Astenna/Nubes/evaluation/movie_review_baseline/client_project/models"
)

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		var account models.Account
		account.Password = r.PostFormValue("Password")
		account.Email = r.PostFormValue("Email")
		account.Nickname = r.PostFormValue("Name")

		_, err := invokeLambdaToGetSingleItem[string](account, "register")

		if err != nil {
			fmt.Fprintf(w, "Error occurred when creating the user")
			return
		}

		fmt.Fprintf(w, "Account created successfully")
		return
	}
	t, _ := template.ParseFiles("templates//register.html")
	t.Execute(w, nil)
}
