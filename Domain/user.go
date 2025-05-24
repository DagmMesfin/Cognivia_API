package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Email     string             `bson:"email" json:"email"`
	Password  string             `bson:"password" json:"-"`
	Name      string             `bson:"name" json:"name"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
}

type UserRepository interface {
	Create(user *User) error
	FindByEmail(email string) (*User, error)
	FindByID(id string) (*User, error)
	Update(user *User) error
	Delete(id string) error
}

type UserUseCase interface {
	Register(user *User) error
	Login(email, password string) (string, error)
	GetUserByID(id string) (*User, error)
	GetUserByEmail(email string) (*User, error)
	VerifyPassword(hash, password string) bool
	GenerateToken(user *User) (string, error)
	UpdateUser(user *User) error
	DeleteUser(id string) error
}
