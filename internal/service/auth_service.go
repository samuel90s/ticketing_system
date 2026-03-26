package service

import (
	"errors"

	"ticketing-system/internal/config"
	"ticketing-system/internal/model"

	"golang.org/x/crypto/bcrypt"
)

// ======================
// REGISTER
// ======================
func Register(name, email, password, role string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return err
	}

	user := model.User{
		Name:     name,
		Email:    email,
		Password: string(hashedPassword),
		Role:     role, // 🔥 TAMBAH INI
	}

	if err := config.DB.Create(&user).Error; err != nil {
		return err
	}

	return nil
}

// ======================
// LOGIN
// ======================
func Login(email, password string) (*model.User, error) {
	var user model.User

	if err := config.DB.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, errors.New("user not found")
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, errors.New("wrong password")
	}

	return &user, nil
}

// ======================
// GET USER
// ======================
func GetUserByID(id uint) (model.User, error) {
	var user model.User
	err := config.DB.First(&user, id).Error
	return user, err
}
