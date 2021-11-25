package interfaces

import (
	model "Dp218Go/domain/dto"
)

type UserRepo interface {
	GetAllUsers() (*model.UserList, error)
	GetUserById(userId int) (model.User, error)
	AddUser(user *model.User) error
	UpdateUser(userId int, userData model.User) (model.User, error)
	DeleteUser(userId int) error
}
