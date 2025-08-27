package auth

import (
	"RESTCryptoServer/internal/db"
	"errors"
)

var ErrInvalidUserData = errors.New("incorrect user data")

type LoginPasswordJSON struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type TokenResponse struct {
    Token string `json:"token"`
}

type AuthService struct {
	UsersDB *db.UserDB
}

func NewAuthService(udb *db.UserDB) *AuthService {
	return &AuthService{UsersDB: udb}
}

func (authService *AuthService) Insert(login string, password string) error {
	return authService.UsersDB.Insert(login, HashPassword(password))
}

func (authService *AuthService) Exist(login string) bool {
	_, err := authService.UsersDB.Get(login)
	
	return err == nil
}

func (authService *AuthService) UserValidation(login string, password string) error {
	realPassword, err := authService.UsersDB.Get(login)
	if err != nil {
		return err
	}

	if !PasswordMatches(realPassword, password) {
		return ErrInvalidUserData
	}

	return nil
}

func (authService *AuthService) Delete(login string) error {
	return authService.UsersDB.Delete(login)
}