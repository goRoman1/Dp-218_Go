package interfaces

import (
	models "Dp218Go/domain/entities"
)

type UserRepo interface {
	GetAllUsers() (*models.UserList, error)
	GetUserById(userId int) (models.User, error)
	AddUser(user *models.User) error
	UpdateUser(userId int, userData models.User) (models.User, error)
	DeleteUser(userId int) error
}
