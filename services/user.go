package services

import (
	"Dp218Go/models"
	"Dp218Go/repositories"
)

type UserService struct {
	repoUser repositories.UserRepo
	repoRole repositories.RoleRepo
}

func NewUserService(repoUser repositories.UserRepo, repoRole repositories.RoleRepo) *UserService {
	return &UserService{repoUser: repoUser, repoRole: repoRole}
}

func (ser *UserService) GetAllUsers() (*models.UserList, error) {
	return ser.repoUser.GetAllUsers()
}

func (ser *UserService) AddUser(user *models.User) error {
	return ser.repoUser.AddUser(user)
}

func (ser *UserService) GetUserByID(userID int) (models.User, error) {
	return ser.repoUser.GetUserByID(userID)
}

func (ser *UserService) DeleteUser(userID int) error {
	return ser.repoUser.DeleteUser(userID)
}

func (ser *UserService) UpdateUser(userID int, userData models.User) (models.User, error) {
	return ser.repoUser.UpdateUser(userID, userData)
}

func (ser *UserService) FindUsersByLoginNameSurname(whatToFind string) (*models.UserList, error) {
	return ser.repoUser.FindUsersByLoginNameSurname(whatToFind)
}

func (ser *UserService) GetAllRoles() (*models.RoleList, error) {
	return ser.repoRole.GetAllRoles()
}

func (ser *UserService) GetRoleByID(roleID int) (models.Role, error) {
	return ser.repoRole.GetRoleByID(roleID)
}

func (ser *UserService) ChangeUsersBlockStatus(userID int) error {
	user, err := ser.repoUser.GetUserByID(userID)
	if err != nil {
		return err
	}
	user.IsBlocked = !user.IsBlocked
	_, err = ser.repoUser.UpdateUser(userID, user)
	return err
}
