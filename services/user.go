package services

import (
	"Dp218Go/models"
	"Dp218Go/repositories"
)

type UserUsecases interface {
	GetAllUsers() (*models.UserList, error)
	GetUserById(userId int) (models.User, error)
	AddUser(user *models.User) error
	UpdateUser(userId int, userData models.User) (models.User, error)
	DeleteUser(userId int) error
}

type UserService struct {
	repo repositories.AnyDatabase
}

func NewUserService(db repositories.AnyDatabase) *UserService {
	return &UserService{db}
}

func (db *UserService) GetAllUsers() (*models.UserList, error) {
	return db.repo.GetAllUsers()
}

func (db *UserService) AddUser(user *models.User) error {
	return db.repo.AddUser(user)
}

func (db *UserService) GetUserById(userId int) (models.User, error) {
	return db.repo.GetUserById(userId)
}

func (db *UserService) DeleteUser(userId int) error {
	return db.repo.DeleteUser(userId)
}

func (db *UserService) UpdateUser(userId int, userData models.User) (models.User, error) {
	return db.repo.UpdateUser(userId, userData)
}
