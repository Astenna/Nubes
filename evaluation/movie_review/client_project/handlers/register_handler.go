package handlers

import (
	"fmt"
	"net/http"
	"text/template"

	clib "github.com/Astenna/Nubes/movie_review/client_lib"
)

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		var account clib.AccountStub
		account.Password = r.PostFormValue("Password")
		account.Email = r.PostFormValue("Email")
		account.Nickname = r.PostFormValue("Name")

		_, err := clib.ExportAccount(account)
		if err != nil {
			fmt.Fprintf(w, "Error occurred when creating the user")
		}

		fmt.Fprintf(w, "Account created successfully")
		return
	}
	t, _ := template.ParseFiles("templates//register.html")
	t.Execute(w, nil)
}
