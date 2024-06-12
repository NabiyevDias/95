package database

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
)

var DB *gorm.DB

type User struct {
	Email    string `gorm:"primaryKey"`
	Password string
}

type Message struct {
	gorm.Model
	UserEmail string
	Content   string
}

func InitDB() {
	var err error
	DB, err = gorm.Open(sqlite.Open("chat.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect database:", err)
	}

	DB.AutoMigrate(&User{}, &Message{})
}


func CreateChat(email string) {
	// Create user record if not exists
	user := User{Email: email}
	if result := DB.FirstOrCreate(&user, User{Email: email}); result.Error != nil {
		log.Println("Error creating user:", result.Error)
		return
	}

	// Create chat record
	message := Message{UserEmail: email, Content: "Chat created"}
	if result := DB.Create(&message); result.Error != nil {
		log.Println("Error creating chat message:", result.Error)
		return
	}

	log.Println("Chat created for user:", email)
}

func DeleteChat(email string) {
	// Delete all messages associated with the user
	if result := DB.Where("user_email = ?", email).Delete(&Message{}); result.Error != nil {
		log.Println("Error deleting messages:", result.Error)
		return
	}

	// Delete user record
	if result := DB.Where("email = ?", email).Delete(&User{}); result.Error != nil {
		log.Println("Error deleting user:", result.Error)
		return
	}

	log.Println("Chat deleted for user:", email)
}
