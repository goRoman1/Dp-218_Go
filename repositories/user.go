package repositories

import (
	"Dp218Go/models"
	"context"
)

type UserRepo interface {
	GetAllUsers() (*models.UserList, error)
	GetUserByID(userID int) (models.User, error)
	GetUserByEmail(email string) (models.User, error)
	AddUser(user *models.User) error
	UpdateUser(userID int, userData models.User) (models.User, error)
	DeleteUser(userID int) error
	FindUsersByLoginNameSurname(whatToFind string) (*models.UserList, error)
}

type RoleRepo interface {
	GetAllRoles() (*models.RoleList, error)
	GetRoleByID(roleID int) (models.Role, error)
}

type AuthRepo interface {
	GetUserByEmail(context.Context, string) (models.User, error)
}
