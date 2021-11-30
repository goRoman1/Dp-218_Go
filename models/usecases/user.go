package usecases

import (
	"Dp218Go/models"
	"fmt"
)

type UserUsecasesRepo interface {
	GetAllUsers() (*models.UserList, error)
	GetUserById(userId int) (models.User, error)
	AddUser(user *models.User) error
	UpdateUser(userId int, userData models.User) (models.User, error)
	DeleteUser(userId int) error
	FindUsersByLoginNameSurname(whatToFind string) (*models.UserList, error)
}

type UserUsecasesGeneral interface {
	ChangeUsersBlockStatus(userId int) error
}

type RoleUsecasesRepo interface {
	GetAllRoles() (*models.RoleList, error)
	GetRoleById(roleId int) (models.Role, error)
}

func FindRoleById(roles *models.RoleList, roleId int) (models.Role, error) {
	for _, v := range roles.Roles {
		if v.ID == roleId{
			return v, nil
		}
	}
	return models.Role{}, fmt.Errorf("not found role id=%d", roleId)
}