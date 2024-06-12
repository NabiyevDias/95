package main

import (
	"chat-messenger/database"
	"chat-messenger/handlers"
	"log"
	"net"
	"net/http"

	"github.com/gorilla/sessions"
)

var store = sessions.NewCookieStore([]byte("super-secret-key"))

func main() {
	database.InitDB()
	handlers.SetStore(store)

	http.HandleFunc("/", handlers.LoginHandler)
	http.HandleFunc("/register", handlers.RegisterHandler)
	http.HandleFunc("/chat", handlers.ChatHandler)
	http.HandleFunc("/send", handlers.SendHandler)
	http.HandleFunc("/messages", handlers.MessagesHandler)

	// Static files
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// Start HTTP server
	go func() {
		log.Println("Starting server on : http://localhost:8080")
		err := http.ListenAndServe(":8080", nil)
		if err != nil {
			log.Printf("Server start error %s", err)
		}
	}()

	// Start TCP server
	go startTCPServer()

	select {} // Keep the main function running
}

func startTCPServer() {
	listener, err := net.Listen("tcp", ":8081")
	if err != nil {
		log.Println("TCP server error:", err)
		return
	}
	defer listener.Close()
	log.Println("TCP server started on :8081")

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Accept error:", err)
			continue
		}
		go handlers.HandleTCPConnection(conn)
	}
}
