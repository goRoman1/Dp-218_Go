package models

import (
	"fmt"
	"net/http"
)

type User struct {
	ID          int    `json:"id"`
	LoginEmail  string `json:"login_email"`
	IsBlocked   bool   `json:"is_blocked"`
	UserName    string `json:"user_name"`
	UserSurname string `json:"user_surname"`
	CreatedAt   string `json:"created_at"`
	RoleID      int    `json:"role_id"`
}

type UserList struct {
	Users []User `json:"users"`
}

func (u *User) Bind(r *http.Request) error {
	if u.LoginEmail == "" {
		return fmt.Errorf("login_email is a required field")
	}

	return nil
}

func (*UserList) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (*User) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}
