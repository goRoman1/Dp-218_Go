package db

import (
	"database/sql"

	"github.com/ITA-Dnipro/Dp-218_Go/models"
)

func (db Database) GetAllUsers() (*models.UserList, error) {
	list := &models.UserList{}
	rows, err := db.Conn.Query("SELECT * FROM Users ORDER BY ID DESC;")
	if err != nil {
		return list, err
	}
	for rows.Next() {
		var user models.User
		err := rows.Scan(&user.ID, &user.LoginEmail, &user.IsBlocked,
			&user.UserName, &user.UserSurname, &user.CreatedAt, &user.RoleID)
		if err != nil {
			return list, err
		}
		list.Users = append(list.Users, user)
	}
	return list, nil
}

func (db Database) AddUser(user *models.User) error {
	var id int
	var created_at string
	query := `INSERT INTO Users(LoginEmail, IsBlocked, UserName, UserSurname, RoleID) VALUES($1, $2, $3, $4, $5) RETURNING ID, CreatedAt;`
	err := db.Conn.QueryRow(query, user.LoginEmail, user.IsBlocked, user.UserName, user.UserSurname, user.RoleID).Scan(&id, &created_at)
	if err != nil {
		return err
	}
	user.ID = id
	user.CreatedAt = created_at
	return nil
}

func (db Database) GetUserById(userId int) (models.User, error) {
	user := models.User{}
	query := `SELECT * FROM Users WHERE ID = $1;`
	row := db.Conn.QueryRow(query, userId)
	switch err := row.Scan(&user.ID, &user.LoginEmail, &user.IsBlocked,
		&user.UserName, &user.UserSurname, &user.CreatedAt, &user.RoleID); err {
	case sql.ErrNoRows:
		return user, ErrNoMatch
	default:
		return user, err
	}
}

func (db Database) DeleteUser(userId int) error {
	query := `DELETE FROM Users WHERE ID = $1;`
	_, err := db.Conn.Exec(query, userId)
	switch err {
	case sql.ErrNoRows:
		return ErrNoMatch
	default:
		return err
	}
}

func (db Database) UpdateUser(userId int, userData models.User) (models.User, error) {
	user := models.User{}
	query := `UPDATE Users SET LoginEmail=$1, IsBlocked=$2, UserName=$3, UserSurname=$4, RoleID=$5 WHERE ID=$6 RETURNING ID, LoginEmail;`
	err := db.Conn.QueryRow(query, userData.LoginEmail, userData.IsBlocked, userData.UserName,
		userData.UserSurname, userData.RoleID, userId).Scan(&user.ID, &user.LoginEmail)
	if err != nil {
		if err == sql.ErrNoRows {
			return user, ErrNoMatch
		}
		return user, err
	}
	return user, nil
}
