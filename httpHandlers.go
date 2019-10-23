package main

import (
	"fmt"
	"html/template"
	"net/http"
)

type formResponse struct {
	Success bool
	Message string
}

func index(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("register.html"))
	if r.Method != http.MethodPost {
		tmpl.Execute(w, nil)
		return
	}

	user := user{
		Username:        r.FormValue("Username"),
		Password:        r.FormValue("Password"),
		ConfirmPassword: r.FormValue("ConfirmPassword"),
	}
	ip := getUserIP(r)
	_ = user
	connectToDB()
	if res, err := userExists(user.Username); !res && err == nil {
		if user.Password == user.ConfirmPassword {
			createUser(user.Username, getMD5Hash(user.Password), ip)
			formResponse := formResponse{
				Success: true,
				Message: fmt.Sprintf("User %s created successfully!", user.Username),
			}
			tmpl.Execute(w, formResponse)
		} else {
			formResponse := formResponse{
				Success: false,
				Message: "Passwords don't match!",
			}
			tmpl.Execute(w, formResponse)
		}
	} else {
		formResponse := formResponse{
			Success: false,
			Message: "Username is already taken!",
		}
		tmpl.Execute(w, formResponse)
	}
	defer db.Close()
}
