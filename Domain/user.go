package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Email      string             `bson:"email" json:"email"`
	Password   string             `bson:"password" json:"password"`
	Name       string             `bson:"name" json:"name"`
	Bio        string             `bson:"bio" json:"bio"`
	ProfilePic string             `bson:"profile_pic" json:"profile_pic"`
	CreatedAt  time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt  time.Time          `bson:"updated_at" json:"updated_at"`
	Settings   UserSettings       `bson:"settings" json:"settings"`
}

type UserSettings struct {
	Theme                string `bson:"theme" json:"theme"`
	Language             string `bson:"language" json:"language"`
	EmailNotifications   bool   `bson:"email_notifications" json:"email_notifications"`
	BrowserNotifications bool   `bson:"browser_notifications" json:"browser_notifications"`
	MobileNotifications  bool   `bson:"mobile_notifications" json:"mobile_notifications"`
}

type UserRepository interface {
	Create(user *User) error
	FindByEmail(email string) (*User, error)
	FindByID(id string) (*User, error)
	Update(user *User) error
	Delete(id string) error
}

type UserUseCase interface {
	Register(user *User) (*User, error) // Changed to return the user to match the return type of Create
	Login(email, password string) (string, error)
	GetUserByID(id string) (*User, error)
	GetUserByEmail(email string) (*User, error)
	VerifyPassword(hash, password string) bool // Changed to bool to match the return type of PasswordComparator
	GenerateToken(user *User) (string, error)
	UpdateUser(user *User) error
	DeleteUser(id string) error
}
