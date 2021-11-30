package usecases

type UserUsecases interface {
	ChangeUsersBlockStatus(userId int) error
}

type RoleUsecases interface {
}