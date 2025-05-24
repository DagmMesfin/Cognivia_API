package usecase

import (
	"errors"
	"log"
	"os"
	"time"

	domain "cognivia-api/Domain"
	"cognivia-api/infrastructure"

	"github.com/golang-jwt/jwt/v5"
)

type userUseCase struct {
	userRepo        domain.UserRepository
	passwordService infrastructure.PasswordService
}

func NewUserUseCase(userRepo domain.UserRepository) domain.UserUseCase {
	return &userUseCase{
		userRepo:        userRepo,
		passwordService: infrastructure.NewPasswordService(),
	}
}

func (u *userUseCase) Register(user *domain.User) (*domain.User, error) {
	// Check if user already exists
	existingUser, err := u.userRepo.FindByEmail(user.Email)
	if err != nil {
		log.Printf("Error checking existing user: %v", err)
		return nil, err
	}
	if existingUser != nil {
		log.Printf("User already exists with email: %s", user.Email)
		return nil, errors.New("user already exists")
	}

	// Log original password
	log.Printf("Original password: %s", user.Password)
	log.Printf("Original password length: %d", len(user.Password))

	// Hash password using password service
	hashedPassword, err := u.passwordService.PasswordHasher(user.Password)
	if err != nil {
		log.Printf("Error hashing password: %v", err)
		return nil, err
	}
	user.Password = hashedPassword

	// Log hashed password details
	log.Printf("Hashed password: %s", user.Password)
	log.Printf("Hashed password length: %d", len(user.Password))

	//defaults for bio, profile pic and settings
	user.Bio = ""
	user.ProfilePic = "https://avatar.iran.liara.run/public/1"
	user.Settings = domain.UserSettings{
		Theme:                "light",
		Language:             "en",
		EmailNotifications:   false,
		BrowserNotifications: false,
		MobileNotifications:  false,
	}

	// Create user
	err = u.userRepo.Create(user)
	if err != nil {
		log.Printf("Error creating user: %v", err)
		return nil, err
	}

	log.Printf("User created successfully: %s", user.Email)
	return user, nil
}

func (u *userUseCase) Login(email, password string) (string, error) {
	user, err := u.GetUserByEmail(email)
	if err != nil {
		return "", err
	}

	if u.VerifyPassword(user.Password, password) {
		return "", errors.New("invalid credentials")
	}

	return u.GenerateToken(user)
}

func (u *userUseCase) GetUserByEmail(email string) (*domain.User, error) {
	user, err := u.userRepo.FindByEmail(email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}
	return user, nil
}

func (u *userUseCase) VerifyPassword(hash, password string) bool {
	return !u.passwordService.PasswordComparator(hash, password)
}

func (u *userUseCase) GenerateToken(user *domain.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID.Hex(),
		"email":   user.Email,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	})

	return token.SignedString([]byte(os.Getenv("JWT_SECRET")))
}

func (u *userUseCase) GetUserByID(id string) (*domain.User, error) {
	return u.userRepo.FindByID(id)
}

func (u *userUseCase) UpdateUser(user *domain.User) error {
	return u.userRepo.Update(user)
}

func (u *userUseCase) DeleteUser(id string) error {
	return u.userRepo.Delete(id)
}
