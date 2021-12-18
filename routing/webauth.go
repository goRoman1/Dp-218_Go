package routing

import (
	"Dp218Go/models"
	"Dp218Go/services"
	"context"
	"errors"
	"net/http"

	"github.com/gorilla/mux"
)

type userKey string

var (
	ukey                  userKey = "user"
	authenticationService *services.AuthService
	ErrSignUp             = errors.New("signup error")
	ErrSignIn             = errors.New("signin error")
)

func AddAuthHandler(router *mux.Router, service *services.AuthService) {
	authenticationService = service
	router.Path("/signup").HandlerFunc(SignUp(authenticationService)).Methods(http.MethodPost)
	router.Path("/signin").HandlerFunc(SignIn(authenticationService)).Methods(http.MethodPost)
	router.Path("/signout").HandlerFunc(SignOut(authenticationService)).Methods(http.MethodGet)
}

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

		if err := sv.SignUp(user); err != nil {

			http.Error(w, ErrSignUp.Error(), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/login", http.StatusFound)
	}
}

func SignIn(sv *services.AuthService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// // TODO implement validation

		req := &services.AuthRequest{
			Email:    r.FormValue("email"),
			Password: r.FormValue("password"),
		}

		if err := sv.SignIn(w, r, req); err != nil {

			EncodeError(FormatHTML, w, &ResponseStatus{
				Err:        ErrSignIn,
				StatusCode: http.StatusForbidden,
				StatusText: ErrSignIn.Error(),
				Message:    "cant get user" + err.Error(),
			})
			return
		}

		http.Redirect(w, r, "/home", http.StatusFound)
	}
}

func SignOut(sv *services.AuthService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		err := sv.SignOut(w, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/login", http.StatusFound)
	}
}

func FilterAuth(sv *services.AuthService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user, err := sv.GetUserFromRequest(r)
			if err != nil {
				http.Error(w, err.Error(), http.StatusForbidden)
				return
			}
			newReq := r.WithContext(context.WithValue(r.Context(), ukey, user))

			next.ServeHTTP(w, newReq)
		})
	}
}

func GetUserFromContext(r *http.Request) *models.User {
	val := r.Context().Value(ukey)
	user, ok := val.(*models.User)

	if ok {
		return user
	}
	return nil
}
