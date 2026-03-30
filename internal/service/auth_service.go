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

func Register(name, email, password string) error {
	// Cek email sudah terdaftar
	var existing model.User
	if err := config.DB.Where("email = ?", email).First(&existing).Error; err == nil {
		return errors.New("email already registered")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}

	user := model.User{
		Name:     name,
		Email:    email,
		Password: string(hashedPassword),
		Role:     "user", // selalu user, tidak bisa diubah dari luar
	}

	return config.DB.Create(&user).Error
}

// ======================
// LOGIN
// ======================

func Login(email, password string) (*model.User, error) {
	var user model.User

	if err := config.DB.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, errors.New("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, errors.New("invalid credentials")
	}

	return &user, nil
}

// ======================
// GET USER BY ID
// ======================

func GetUserByID(id uint) (model.User, error) {
	var user model.User
	err := config.DB.First(&user, id).Error
	return user, err
}
