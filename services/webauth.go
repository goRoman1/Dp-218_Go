package services

import (
	"Dp218Go/models"
	"Dp218Go/repositories"
	"encoding/gob"
	"fmt"
	"net/http"

	"github.com/gorilla/sessions"
)

type AuthService struct {
	DB        repositories.UserRepo
	sessStore sessions.Store
}

const (
	SessionName = "login"
	SessionVal  = "user"
)

func NewAuthService(db repositories.UserRepo, store sessions.Store) *AuthService {

	gob.Register(&models.User{})

	return &AuthService{
		DB:        db,
		sessStore: store,
	}
}

func (sv *AuthService) GetUserFromRequest(r *http.Request) (*models.User, error) {
	sess, err := sv.sessStore.Get(r, SessionName)
	if err != nil {
		return nil, err
	}

	val, ok := sess.Values[SessionVal]
	if !ok {
		return nil, fmt.Errorf("%s", "no user in session")
	}

	var user = &models.User{}
	if user, ok = val.(*models.User); !ok {
		return nil, fmt.Errorf("%s", "no user in session")

	}

	return user, nil

}

func (sv *AuthService) GetSessionStore() sessions.Store {
	return sv.sessStore
}
