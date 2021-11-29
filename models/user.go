package models

import (
	"time"
)

type Role struct {
	ID int `json:"id"`
	Name       string `json:"name"`
	IsAdmin    bool	`json:"is_admin"`
	IsUser     bool	`json:"is_user"`
	IsSupplier bool `json:"is_supplier"`
}

type RoleList struct {
	Roles []Role `json:"roles"`
}

type User struct {
	ID          int    `json:"id"`
	LoginEmail  string `json:"login_email"`
	IsBlocked   bool   `json:"is_blocked"`
	UserName    string `json:"user_name"`
	UserSurname string `json:"user_surname"`
	CreatedAt   time.Time `json:"created_at"`
	Role      Role    `json:"role"`
}

type UserList struct {
	Users []User `json:"users"`
}
