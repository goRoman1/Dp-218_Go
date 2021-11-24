package routing

import (
	"net/http"
	"strconv"

	model "Dp218Go/domain/entities"
	iface "Dp218Go/domain/interfaces"
	repo "Dp218Go/repositories"

	"github.com/go-chi/render"
	"github.com/gorilla/mux"
)

var	userRepo iface.UserRepo
var userIDKey = "userID"

var keyRoutes = []Route{
	{
		Uri:         `/users`,
		Method:    http.MethodGet,
		Handler:	getAllUsers,
	},
	{
		Uri:         `/user/{`+userIDKey+`}`,
		Method:     http.MethodGet,
		Handler:	getUser,
	},
	{
		Uri:         `/user`,
		Method:     http.MethodPost,
		Handler:	createUser,
	},
	{
		Uri:         `/user/{`+userIDKey+`}`,
		Method:     http.MethodPut,
		Handler:	updateUser,
	},
	{
		Uri:         `/user/{`+userIDKey+`}`,
		Method:     http.MethodDelete,
		Handler:	deleteUser,
	},
}

func AddUserHandler(router *mux.Router, repo iface.UserRepo) {
	userRepo = repo
	for _, rt := range keyRoutes{
		router.Path(rt.Uri).HandlerFunc(rt.Handler).Methods(rt.Method)
	}
}

func createUser(w http.ResponseWriter, r *http.Request) {
	user := &model.User{}
	if err := render.Bind(r, user); err != nil {
		render.Render(w, r, ErrBadRequest)
		return
	}
	if err := userRepo.AddUser(user); err != nil {
		render.Render(w, r, ErrorRenderer(err))
		return
	}
	if err := render.Render(w, r, user); err != nil {
		render.Render(w, r, ServerErrorRenderer(err))
		return
	}
}

func getAllUsers(w http.ResponseWriter, r *http.Request) {
	users, err := userRepo.GetAllUsers()
	if err != nil {
		render.Render(w, r, ServerErrorRenderer(err))
		return
	}
	if err := render.Render(w, r, users); err != nil {
		render.Render(w, r, ErrorRenderer(err))
		return
	}
}

func getUser(w http.ResponseWriter, r *http.Request) {
	userId, err := strconv.Atoi(mux.Vars(r)[userIDKey])
	if err != nil {
		render.Render(w, r, ErrorRenderer(err))
	}
	user, err := userRepo.GetUserById(userId)
	if err != nil {
		if err == repo.ErrNoMatch {
			render.Render(w, r, ErrNotFound)
		} else {
			render.Render(w, r, ErrorRenderer(err))
		}
		return
	}
	if err := render.Render(w, r, &user); err != nil {
		render.Render(w, r, ServerErrorRenderer(err))
		return
	}
}

func deleteUser(w http.ResponseWriter, r *http.Request) {
	userId, err := strconv.Atoi(mux.Vars(r)[userIDKey])
	if err != nil {
		render.Render(w, r, ErrorRenderer(err))
	}
	err = userRepo.DeleteUser(userId)
	if err != nil {
		if err == repo.ErrNoMatch {
			render.Render(w, r, ErrNotFound)
		} else {
			render.Render(w, r, ServerErrorRenderer(err))
		}
		return
	}
	render.Render(w, r, StatusOK)
}

func updateUser(w http.ResponseWriter, r *http.Request) {
	userId, err := strconv.Atoi(mux.Vars(r)[userIDKey])
	if err != nil {
		render.Render(w, r, ErrorRenderer(err))
	}
	userData := model.User{}
	if err := render.Bind(r, &userData); err != nil {
		render.Render(w, r, ErrBadRequest)
		return
	}
	userData, err = userRepo.UpdateUser(userId, userData)
	if err != nil {
		if err == repo.ErrNoMatch {
			render.Render(w, r, ErrNotFound)
		} else {
			render.Render(w, r, ServerErrorRenderer(err))
		}
		return
	}
	if err := render.Render(w, r, &userData); err != nil {
		render.Render(w, r, ServerErrorRenderer(err))
		return
	}
}
