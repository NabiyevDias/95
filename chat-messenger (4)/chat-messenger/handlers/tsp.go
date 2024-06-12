package handlers

import (
	"bufio"
	"chat-messenger/database"
	"fmt"
	"log"
	"net"
	"strings"

	"gorm.io/gorm"
)

type client struct {
	conn net.Conn
	name string
}

var clients []client

func HandleTCPConnection(conn net.Conn) {
	// Чтение email пользователя
	email, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		log.Println("Error reading email:", err)
		return
	}
	email = strings.TrimSpace(email)
	fmt.Println(email)

	// Проверка пользователя в базе данных
	var user database.User
	result := database.DB.Where("email = ?", Decrypt(email)).First(&user)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			// Создаем нового пользователя
			newUser := database.User{Email: email}
			database.DB.Create(&newUser)
			log.Println("New user registered: " + Decrypt(email))
			log.Println("AES cipher: " + Encrypt(email))
		} else {
			log.Println("Error querying user:", result.Error)
			conn.Close()
			return
		}
	}

	newClient := client{conn: conn, name: email}
	clients = append(clients, newClient)

	for {
		message, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			break
		}
		message = strings.TrimSpace(message)
		encryptedMessage := Encrypt("------" + message)
		log.Printf("%s: %s", email, encryptedMessage)
		broadcast(fmt.Sprintf("%s: %s", email, encryptedMessage), &newClient)
	}

	removeClient(&newClient)
	conn.Close()
}

func broadcast(message string, sender *client) {
	for _, c := range clients {
		if c.conn == sender.conn {
			continue
		}
		_, err := fmt.Fprintln(c.conn, message)
		if err != nil {
			log.Println("Error broadcasting to", c.name, ":", err)
		}
	}
}

func removeClient(clientToRemove *client) {
	for i, c := range clients {
		if c.conn == clientToRemove.conn {
			clients = append(clients[:i], clients[i+1:]...)
			break
		}
	}
}

// package handlers

// import (
// 	"bufio"
// 	"fmt"
// 	"log"
// 	"net"
// 	"strings"
// 	"chat-messenger/database"

// 	"gorm.io/gorm"
// )

// type client struct {
// 	conn net.Conn
// 	name string
// }

// var clients []client

// func HandleTCPConnection(conn net.Conn) {

// 	defer conn.Close()

// 	email, err := bufio.NewReader(conn).ReadString('\n')
// 	if err != nil {
// 		log.Println("Error reading email:", err)
// 		return
// 	}
// 	email = strings.TrimSpace(email)
// 	fmt.Println(email)

// 	// Проверка пользователя в базе данных
// 	var user database.User
// 	result := database.DB.Where("email = ?", Decrypt(email)).First(&user)
// 	if result.Error != nil {
// 		if result.Error == gorm.ErrRecordNotFound {
// 			// Создаем нового пользователя
// 			newUser := database.User{Email: email}
// 			database.DB.Create(&newUser)
// 			log.Println("New user registered: " + email)
// 			log.Println("AES cipher: " + Encrypt(email))
// 		} else {
// 			log.Println("Error querying user:", result.Error)
// 			conn.Close()
// 			return
// 		}
// 	}

// 	newClient := client{conn: conn, name: email}
// 	clients = append(clients, newClient)

// 	log.Println(email, "connected")
// 	broadcast(Encrypt(email+" joined the chat"), &newClient)

// 	for {
// 		message, err := bufio.NewReader(conn).ReadString('\n')
// 		if err != nil {
// 			break
// 		}
// 		message = strings.TrimSpace(message)
// 		encryptedMessage := Encrypt("------" + message)
// 		log.Printf("%s: %s", email, encryptedMessage)
// 		broadcast(fmt.Sprintf("%s: %s", email, encryptedMessage), &newClient)
// 	}

// 	removeClient(&newClient)
// 	conn.Close()
// }

// func broadcast(message string, sender *client) {
// 	for _, c := range clients {
// 		if c.conn == sender.conn {
// 			continue
// 		}
// 		_, err := fmt.Fprintln(c.conn, message)
// 		if err != nil {
// 			log.Println("Error broadcasting to", c.name, ":", err)
// 		}
// 	}
// }

// func removeClient(clientToRemove *client) {
// 	for i, c := range clients {
// 		if c.conn == clientToRemove.conn {
// 			clients = append(clients[:i], clients[i+1:]...)
// 			break
// 		}
// 	}
// }
