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

	querySQL := `SELECT * FROM users ORDER BY id DESC;`
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
	querySQL := `INSERT INTO users(login_email, is_blocked, user_name, user_surname, role_id) 
		VALUES($1, $2, $3, $4, $5)
		RETURNING id, created_at;`
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

	querySQL := `SELECT * FROM users WHERE id = $1;`
	row := pg.QueryResultRow(context.Background(), querySQL, userId)
	var roleId int
	err := row.Scan(&user.ID, &user.LoginEmail, &user.IsBlocked,
		&user.UserName, &user.UserSurname, &user.CreatedAt, &roleId)
	user.Role, err = pg.GetRoleById(roleId)

	return user, err
}

func (pg *Postgres) DeleteUser(userId int) error {
	querySQL := `DELETE FROM users WHERE id = $1;`
	_, err := pg.QueryExec(context.Background(), querySQL, userId)
	return err
}

func (pg *Postgres) UpdateUser(userId int, userData models.User) (models.User, error) {
	user := models.User{}
	querySQL := `UPDATE users 
		SET login_email=$1, is_blocked=$2, user_name=$3, user_surname=$4, role_id=$5 
		WHERE id=$6 
		RETURNING id, created_at, login_email, is_blocked, user_name, user_surname, role_id;`
	var roleId int
	err := pg.QueryResultRow(context.Background(), querySQL, userData.LoginEmail, userData.IsBlocked, userData.UserName,
		userData.UserSurname, userData.Role.ID, userId).Scan(&user.ID, &user.CreatedAt, &user.LoginEmail, &user.IsBlocked, &user.UserName, &user.UserSurname, &roleId)
	if err != nil {
		return user, err
	}
	user.Role, err = pg.GetRoleById(roleId)
	if err != nil {
		return user, err
	}
	return user, nil
}

func (pg *Postgres) FindUsersByLoginNameSurname(whatToFind string) (*models.UserList, error) {
	list := &models.UserList{}

	roles, err := pg.GetAllRoles()
	if err != nil {
		return list, err
	}

	querySQL := `SELECT * FROM users 
		WHERE LOWER(login_email) LIKE LOWER($1) 
			OR LOWER(user_name) LIKE LOWER($1) 
			OR LOWER(user_surname) LIKE LOWER($1) 
		ORDER BY id DESC;`
	rows, err := pg.QueryResult(context.Background(), querySQL, whatToFind+"%")
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
