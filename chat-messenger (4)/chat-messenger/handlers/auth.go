package handlers

import (
	"chat-messenger/database"
	"fmt"
	"html/template"
	"net/http"
	"sync"

	"github.com/gorilla/sessions"
	"gorm.io/gorm"
)

var (
	users     = make(map[string]string)
	userMutex = &sync.Mutex{}
	store     = sessions.NewCookieStore([]byte("secret-key"))

	// Use absolute or relative paths correctly
	loginTmpl = template.Must(template.ParseFiles("static/login/login.html"))
)

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		email := r.FormValue("email")
		password := r.FormValue("password")

		// Store user credentials in the database
		user := database.User{Email: email, Password: password}
		if result := database.DB.Create(&user); result.Error != nil {
			fmt.Println("Error creating user:", result.Error)
			http.Redirect(w, r, "/register", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	loginTmpl.Execute(w, nil)
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		email := r.FormValue("email")
		password := r.FormValue("password")

		// Check user credentials from the database
		var user database.User
		if result := database.DB.Where("email = ?", email).First(&user); result.Error != nil {
			if result.Error == gorm.ErrRecordNotFound {
				fmt.Println("Incorrect email or password")
				http.Redirect(w, r, "/", http.StatusSeeOther)
				return
			}
			fmt.Println("Error querying user:", result.Error)
			http.Redirect(w, r, "/", http.StatusInternalServerError)
			return
		}

		if user.Password == password {
			// Set up session
			session, _ := store.Get(r, "chat-session")
			session.Values["email"] = email
			err := session.Save(r, w)
			if err != nil {
				fmt.Println("Error saving session:", err)
			}
			http.Redirect(w, r, "/chat", http.StatusSeeOther)
			return
		}

		// Incorrect login attempt
		fmt.Println("Incorrect email or password")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	loginTmpl.Execute(w, nil)
}
