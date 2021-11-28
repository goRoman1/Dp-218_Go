package usecases

import "Dp218Go/models"

type UserUsecases interface {
	GetAllUsers() (*models.UserList, error)
	GetUserById(userId int) (models.User, error)
	AddUser(user *models.User) error
	UpdateUser(userId int, userData models.User) (models.User, error)
	DeleteUser(userId int) error
}
