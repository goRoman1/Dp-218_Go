package entities

import (
	"net/http"
	"time"
)

type User struct {
	ID          int    `json:"id"`
	LoginEmail  string `json:"login_email"`
	IsBlocked   bool   `json:"is_blocked"`
	UserName    string `json:"user_name"`
	UserSurname string `json:"user_surname"`
	CreatedAt   time.Time `json:"created_at"`
	RoleID      int    `json:"role_id"`
}

type UserList struct {
	Users []User `json:"users"`
}

func (u *User) Bind(r *http.Request) error {
	return nil
}

func (*UserList) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (*User) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}
