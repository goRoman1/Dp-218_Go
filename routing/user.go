package routing

import (
	"Dp218Go/models"
	"net/http"
	"strconv"

	"Dp218Go/services"
	"github.com/gorilla/mux"
)

var userService *services.UserService
var userIDKey = "userID"

var keyRoutes = []Route{
	{
		Uri:     `/users`,
		Method:  http.MethodGet,
		Handler: getAllUsers,
	},
	{
		Uri:     `/user/{` + userIDKey + `}`,
		Method:  http.MethodGet,
		Handler: getUser,
	},
	{
		Uri:     `/user`,
		Method:  http.MethodPost,
		Handler: createUser,
	},
	{
		Uri:     `/user/{` + userIDKey + `}`,
		Method:  http.MethodPut,
		Handler: updateUser,
	},
	{
		Uri:     `/user/{` + userIDKey + `}`,
		Method:  http.MethodDelete,
		Handler: deleteUser,
	},
}

func AddUserHandler(router *mux.Router, service *services.UserService) {
	userService = service
	for _, rt := range keyRoutes {
		router.Path(rt.Uri).HandlerFunc(rt.Handler).Methods(rt.Method)
	}
}

func createUser(w http.ResponseWriter, r *http.Request) {
	user := &models.User{}
	DecodeRequest(FormatJSON, w, r, user)

	if err := userService.AddUser(user); err != nil {
		EncodeError(FormatJSON, w, ErrorRendererDefault(err))
		return
	}

	EncodeAnswer(FormatJSON, w, user)
}

func getAllUsers(w http.ResponseWriter, r *http.Request) {
	users, err := userService.GetAllUsers()
	if err != nil {
		ServerErrorRender(FormatJSON, w)
		return
	}

	EncodeAnswer(FormatJSON, w, users)
}

func getUser(w http.ResponseWriter, r *http.Request) {
	userId, err := strconv.Atoi(mux.Vars(r)[userIDKey])
	if err != nil {
		EncodeError(FormatJSON, w, ErrorRendererDefault(err))
	}
	user, err := userService.GetUserById(userId)
	if err != nil {
		EncodeError(FormatJSON, w, ErrorRendererDefault(err))
		return
	}

	EncodeAnswer(FormatJSON, w, &user)
}

func deleteUser(w http.ResponseWriter, r *http.Request) {
	userId, err := strconv.Atoi(mux.Vars(r)[userIDKey])
	if err != nil {
		EncodeError(FormatJSON, w, ErrorRendererDefault(err))
	}
	err = userService.DeleteUser(userId)
	if err != nil {
		ServerErrorRender(FormatJSON, w)
		return
	}
	EncodeError(FormatJSON, w, ErrorRenderer(err, "success", http.StatusOK))
}

func updateUser(w http.ResponseWriter, r *http.Request) {
	userId, err := strconv.Atoi(mux.Vars(r)[userIDKey])
	if err != nil {
		EncodeError(FormatJSON, w, ErrorRendererDefault(err))
	}
	userData := models.User{}
	DecodeRequest(FormatJSON, w, r, &userData)
	userData, err = userService.UpdateUser(userId, userData)
	if err != nil {
		ServerErrorRender(FormatJSON, w)
		return
	}

	EncodeAnswer(FormatJSON, w, &userData)
}
