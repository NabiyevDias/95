package handlers

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net"
	"net/http"
	"sync"
)

type Message struct {
	Username string
	Content  string
}

var tmpl = template.Must(template.ParseFiles("static/index/index.html"))

var (
	messages []Message
	mu       sync.Mutex
)

func ChatHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "chat-session")
	email, ok := session.Values["email"].(string)
	if !ok || email == "" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	mu.Lock()
	defer mu.Unlock()
	err := tmpl.Execute(w, messages)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func SendHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		session, _ := store.Get(r, "chat-session")
		email, ok := session.Values["email"].(string)
		if !ok || email == "" {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		username := email
		content := r.FormValue("content")
		if username != "" && content != "" {
			mu.Lock()
			messages = append(messages, Message{Username: username, Content: content})
			mu.Unlock()

			go func() {
				conn, err := net.Dial("tcp", "localhost:8081")
				if err != nil {
					log.Println("TCP client error:", err)
					return
				}
				defer conn.Close()

				fmt.Fprint(conn, Encrypt(username+": "+content))
			}()
		}
	}
	http.Redirect(w, r, "/chat", http.StatusSeeOther)
}

func MessagesHandler(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(messages)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
