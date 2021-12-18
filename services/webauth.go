package services

import (
	"Dp218Go/models"
	"Dp218Go/repositories"
	"Dp218Go/utils"
	"encoding/gob"
	"fmt"
	"net/http"

	"github.com/gorilla/sessions"
)

type userKey string

type AuthService struct {
	DB        repositories.UserRepo
	sessStore sessions.Store
}

const (
	ukey        userKey = "user"
	sessionName         = "login"
	sessionVal          = "user"
)

func NewAuthService(db repositories.UserRepo, store sessions.Store) *AuthService {

	gob.Register(&models.User{})
	return &AuthService{
		DB:        db,
		sessStore: store,
	}
}

type AuthRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (sv *AuthService) SignUp(user *models.User) error {
	pass, err := utils.HashPassword(user.Password)
	if err != nil {
		return err
	}
	user.Password = pass

	err = sv.DB.AddUser(user)

	if err != nil {
		return err
	}
	return nil
}

func (sv *AuthService) SignIn(w http.ResponseWriter, r *http.Request, authreq *AuthRequest) error {
	user, err := sv.DB.GetUserByEmail(authreq.Email)

	if err != nil {
		return err
	}

	if err := utils.CheckPassword(user.Password, authreq.Password); err != nil {
		return err
	}

	session, err := sv.GetSessionStore().Get(r, sessionName)
	if err != nil {
		return err
	}

	session.Values[sessionVal] = user
	err = session.Save(r, w)
	if err != nil {
		return err
	}

	return nil
}

func (sv *AuthService) SignOut(w http.ResponseWriter, r *http.Request) error {
	session, err := sv.GetSessionStore().Get(r, sessionName)
	if err != nil {
		return err
	}

	session.Values[sessionVal] = nil
	session.Options.MaxAge = -1

	err = session.Save(r, w)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func (sv *AuthService) GetUserFromRequest(r *http.Request) (*models.User, error) {
	sess, err := sv.sessStore.Get(r, sessionName)
	if err != nil {
		return nil, err
	}

	val, ok := sess.Values[sessionVal]
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
