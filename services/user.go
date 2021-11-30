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

func (ser *UserService) GetUserById(userId int) (models.User, error) {
	return ser.repoUser.GetUserById(userId)
}

func (ser *UserService) DeleteUser(userId int) error {
	return ser.repoUser.DeleteUser(userId)
}

func (ser *UserService) UpdateUser(userId int, userData models.User) (models.User, error) {
	return ser.repoUser.UpdateUser(userId, userData)
}

func (ser *UserService) FindUsersByLoginNameSurname(whatToFind string) (*models.UserList, error) {
	return ser.repoUser.FindUsersByLoginNameSurname(whatToFind)
}

func (ser *UserService) GetAllRoles() (*models.RoleList, error) {
	return ser.repoRole.GetAllRoles()
}

func (ser *UserService) GetRoleById(roleId int) (models.Role, error) {
	return ser.repoRole.GetRoleById(roleId)
}

func (ser *UserService) ChangeUsersBlockStatus(userId int) error {
	user, err:= ser.repoUser.GetUserById(userId)
	if err!=nil{
		return err
	}
	user.IsBlocked = !user.IsBlocked
	_, err = ser.repoUser.UpdateUser(userId, user)
	return err
}