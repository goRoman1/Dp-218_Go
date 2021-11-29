package postgres

import (
	"Dp218Go/models"
	"Dp218Go/models/usecases"
	"context"
	"time"
)

func (pg *Postgres) GetAllUsers() (*models.UserList, error) {
	list := &models.UserList{}

	roles, err := pg.GetAllRoles()
	if err != nil {
		return list, err
	}

	querySQL := `SELECT * FROM Users ORDER BY ID DESC;`
	rows, err := pg.QueryResult(context.Background(), querySQL)
	if err != nil {
		return list, err
	}

	for rows.Next() {
		var user models.User
		var roleId int
		err := rows.Scan(&user.ID, &user.LoginEmail, &user.IsBlocked,
			&user.UserName, &user.UserSurname, &user.CreatedAt, &roleId)
		if err != nil {
			return list, err
		}

		user.Role, err = usecases.FindRoleById(roles, roleId)
		if err != nil {
			return list, err
		}

		list.Users = append(list.Users, user)
	}
	return list, nil
}

func (pg *Postgres) AddUser(user *models.User) error {
	var id int
	var createdAt time.Time
	querySQL := `INSERT INTO Users(LoginEmail, IsBlocked, UserName, UserSurname, RoleID) VALUES($1, $2, $3, $4, $5) RETURNING ID, CreatedAt;`
	err := pg.QueryResultRow(context.Background(), querySQL, user.LoginEmail, user.IsBlocked, user.UserName, user.UserSurname, user.Role.ID).Scan(&id, &createdAt)
	if err != nil {
		return err
	}
	user.ID = id
	user.CreatedAt = createdAt
	return nil
}

func (pg *Postgres) GetUserById(userId int) (models.User, error) {
	user := models.User{}

	querySQL := `SELECT * FROM Users WHERE ID = $1;`
	row := pg.QueryResultRow(context.Background(), querySQL, userId)
	var roleId int
	err := row.Scan(&user.ID, &user.LoginEmail, &user.IsBlocked,
		&user.UserName, &user.UserSurname, &user.CreatedAt, &roleId)
	user.Role, err = pg.GetRoleById(roleId)

	return user, err
}

func (pg *Postgres) DeleteUser(userId int) error {
	querySQL := `DELETE FROM Users WHERE ID = $1;`
	_, err := pg.QueryExec(context.Background(), querySQL, userId)
	return err
}

func (pg *Postgres) UpdateUser(userId int, userData models.User) (models.User, error) {
	user := models.User{}
	querySQL := `UPDATE Users SET LoginEmail=$1, IsBlocked=$2, UserName=$3, UserSurname=$4, RoleID=$5 WHERE ID=$6 RETURNING ID, LoginEmail;`
	err := pg.QueryResultRow(context.Background(), querySQL, userData.LoginEmail, userData.IsBlocked, userData.UserName,
		userData.UserSurname, userData.Role.ID, userId).Scan(&user.ID, &user.LoginEmail)
	if err != nil {
		return user, err
	}
	return user, nil
}

func (pg *Postgres) GetAllRoles() (*models.RoleList, error) {
	list := &models.RoleList{}
	querySQL := `SELECT * FROM roles ORDER BY id DESC;`
	rows, err := pg.QueryResult(context.Background(), querySQL)
	if err != nil {
		return list, err
	}
	for rows.Next() {
		var role models.Role
		err := rows.Scan(&role.ID, &role.Name, &role.IsAdmin, &role.IsUser, &role.IsSupplier)
		if err != nil {
			return list, err
		}
		list.Roles = append(list.Roles, role)
	}
	return list, nil
}

func (pg *Postgres) GetRoleById(roleId int) (models.Role, error) {
	role := models.Role{}
	querySQL := `SELECT * FROM roles WHERE id = $1;`
	row := pg.QueryResultRow(context.Background(), querySQL, roleId)
	err := row.Scan(&role.ID, &role.Name, &role.IsAdmin, &role.IsUser, &role.IsSupplier)
	return role, err
}
