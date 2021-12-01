package repositories

import (
	"Dp218Go/models"
	"context"
)

type UserRepo interface {
	GetAllUsers() (*models.UserList, error)
	GetUserById(userId int) (models.User, error)
	GetUserByEmail(email string) (models.User, error)
	AddUser(user *models.User) error
	UpdateUser(userId int, userData models.User) (models.User, error)
	DeleteUser(userId int) error
	FindUsersByLoginNameSurname(whatToFind string) (*models.UserList, error)
}

type RoleRepo interface {
	GetAllRoles() (*models.RoleList, error)
	GetRoleById(roleId int) (models.Role, error)
}

type AuthRepo interface {
	GetUserByEmail(context.Context, string) (models.User, error)
}
