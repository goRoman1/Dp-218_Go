package repositories

import (
	"context"
	"database/sql"
	"time"

	model "Dp218Go/domain/entities"
	"Dp218Go/pkg/postgres"
)

type UserRepoDb struct {
	*postgres.Postgres
}

func New(pg *postgres.Postgres) *UserRepoDb {
	return &UserRepoDb{pg}
}

func (r *UserRepoDb) GetAllUsers() (*model.UserList, error) {
	list := &model.UserList{}
	r.QuerySQL = `SELECT * FROM Users ORDER BY ID DESC;`
	rows, err := r.QueryResult(context.Background())
	if err != nil {
		return list, err
	}
	for rows.Next() {
		var user model.User
		err := rows.Scan(&user.ID, &user.LoginEmail, &user.IsBlocked,
			&user.UserName, &user.UserSurname, &user.CreatedAt, &user.RoleID)
		if err != nil {
			return list, err
		}
		list.Users = append(list.Users, user)
	}
	return list, nil
}

func (r *UserRepoDb) AddUser(user *model.User) error {
	var id int
	var createdAt time.Time
	r.QuerySQL = `INSERT INTO Users(LoginEmail, IsBlocked, UserName, UserSurname, RoleID) VALUES($1, $2, $3, $4, $5) RETURNING ID, CreatedAt;`
	err := r.QueryResultRow(context.Background(), user.LoginEmail, user.IsBlocked, user.UserName, user.UserSurname, user.RoleID).Scan(&id, &createdAt)
	if err != nil {
		return err
	}
	user.ID = id
	user.CreatedAt = createdAt
	return nil
}

func (r *UserRepoDb) GetUserById(userId int) (model.User, error) {
	user := model.User{}
	r.QuerySQL = `SELECT * FROM Users WHERE ID = $1;`
	row := r.QueryResultRow(context.Background(), userId)
	switch err := row.Scan(&user.ID, &user.LoginEmail, &user.IsBlocked,
		&user.UserName, &user.UserSurname, &user.CreatedAt, &user.RoleID); err {
	case sql.ErrNoRows:
		return user, ErrNoMatch
	default:
		return user, err
	}
}

func (r *UserRepoDb) DeleteUser(userId int) error {
	r.QuerySQL = `DELETE FROM Users WHERE ID = $1;`
	_, err := r.QueryExec(context.Background(), userId)
	switch err {
	case sql.ErrNoRows:
		return ErrNoMatch
	default:
		return err
	}
}

func (r *UserRepoDb) UpdateUser(userId int, userData model.User) (model.User, error) {
	user := model.User{}
	r.QuerySQL = `UPDATE Users SET LoginEmail=$1, IsBlocked=$2, UserName=$3, UserSurname=$4, RoleID=$5 WHERE ID=$6 RETURNING ID, LoginEmail;`
	err := r.QueryResultRow(context.Background(), userData.LoginEmail, userData.IsBlocked, userData.UserName,
		userData.UserSurname, userData.RoleID, userId).Scan(&user.ID, &user.LoginEmail)
	if err != nil {
		if err == sql.ErrNoRows {
			return user, ErrNoMatch
		}
		return user, err
	}
	return user, nil
}
