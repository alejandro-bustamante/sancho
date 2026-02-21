package service

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/alejandro-bustamante/sancho/server/internal/model"
	"github.com/alejandro-bustamante/sancho/server/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	db repository.Database
}

func NewUserService(db repository.Database) *UserService {
	return &UserService{db: db}
}

func (s *UserService) RegisterUser(ctx context.Context, username, password, email string) (*model.User, error) {
	hashedPassword, err := s.hashPassword(password)
	if err != nil {
		return nil, err
	}

	userParams := repository.InsertUserParams{
		Username:     username,
		PasswordHash: hashedPassword,
		Email:        sql.NullString{String: email, Valid: email != ""},
	}

	userDB, err := s.db.InsertUser(ctx, userParams)
	if err != nil {
		return nil, err
	}

	user := model.UserFromDB(userDB)
	return &user, nil
}

func (s *UserService) AuthenticateUser(ctx context.Context, username, password string) (bool, error) {
	userDB, err := s.db.GetUserByUsername(ctx, username)
	if err != nil {
		return false, err
	}

	return s.checkPassword(userDB.PasswordHash, password), nil
}

func (s *UserService) hashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}
	return string(hash), nil
}

func (s *UserService) checkPassword(hash, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
