package routing

import (
	"Dp218Go/auth"
	"Dp218Go/models"
	"Dp218Go/services"
	"context"
	"errors"
	"fmt"
	"net/http"
)

type Uid string

var (
	uid       Uid = "user"
	ErrSignUp     = errors.New("signup error")
	ErrSignIn     = errors.New("signin error")
)

func SignUp(sv *services.AuthService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// TODO implement validation
		user := &models.User{
			LoginEmail:  r.FormValue("email"),
			IsBlocked:   true,
			UserName:    r.FormValue("name"),
			UserSurname: r.FormValue("surname"),
			Role:        models.Role{ID: 2},
			Password:    r.FormValue("password"),
		}

		pass, err := auth.HashPassword(user.Password)
		if err != nil {
			fmt.Println(err)
			http.Error(w, ErrSignUp.Error(), http.StatusInternalServerError)
			return
		}
		user.Password = pass

		err = sv.DB.AddUser(user)

		if err != nil {
			fmt.Println(err)
			http.Error(w, ErrSignUp.Error(), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/login", http.StatusFound)
	}
}

func SignIn(sv *services.AuthService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// // TODO implement validation

		type authRequest struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}

		req := authRequest{
			Email:    r.FormValue("email"),
			Password: r.FormValue("password"),
		}
		fmt.Println("got req: ", req)

		user, err := sv.DB.GetUserByEmail(req.Email)
		if err != nil {
			//fmt.Println("cant get user", err)
			//http.Error(w, ErrSignIn.Error(), http.StatusForbidden)
			EncodeError(FormatHTML, w, &ResponseStatus{
				Err:        ErrSignIn,
				StatusCode: http.StatusForbidden,
				StatusText: ErrSignIn.Error(),
				Message:    "cant get user" + err.Error(),
			})
			return
		}

		if err := auth.CheckPassword(user.Password, req.Password); err != nil {
			fmt.Println(err)
			http.Error(w, ErrSignIn.Error(), http.StatusForbidden)
			return
		}

		session, err := sv.GetSessionStore().Get(r, services.SessionName)
		if err != nil {
			fmt.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		session.Values[services.SessionVal] = user
		err = session.Save(r, w)
		if err != nil {
			fmt.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/home", http.StatusFound)
	}
}

func SignOut(sv *services.AuthService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, err := sv.GetSessionStore().Get(r, services.SessionName)
		if err != nil {
			fmt.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		session.Values[services.SessionVal] = nil
		session.Options.MaxAge = -1
		err = session.Save(r, w)
		if err != nil {
			fmt.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/login", http.StatusFound)
	}
}

func WrapUserContext(sv *services.AuthService, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, err := sv.GetUserFromRequest(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusForbidden)
			return
		}
		newReq := r.WithContext(context.WithValue(r.Context(), uid, user))

		next(w, newReq)
	}
}

func GetUserFromContext(r *http.Request) *models.User {
	val := r.Context().Value(uid)
	user, ok := val.(models.User)
	if ok {
		return &user
	}
	return nil
}

// usage
// WrapUserContext(authService, endPointUser)
// -> GetUserFromContext(r) = user
