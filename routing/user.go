package routing

import (
	"Dp218Go/models"
	"fmt"
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
		Method:  http.MethodPost,
		Handler: updateUser,
	},
	{
		Uri:     `/user/{` + userIDKey + `}`,
		Method:  http.MethodDelete,
		Handler: deleteUser,
	},
}

type userRole struct {
	models.User
}

func (ur *userRole) ListOfRoles() []models.Role {
	roles, _ := userService.GetAllRoles()
	return roles.Roles
}


func AddUserHandler(router *mux.Router, service *services.UserService) {
	userService = service
	for _, rt := range keyRoutes {
		router.Path(rt.Uri).HandlerFunc(rt.Handler).Methods(rt.Method)
		router.Path(APIprefix + rt.Uri).HandlerFunc(rt.Handler).Methods(rt.Method)
	}
}

func createUser(w http.ResponseWriter, r *http.Request) {
	format := GetFormatFromRequest(r)

	user := &models.User{}
	DecodeRequest(FormatJSON, w, r, user, nil)

	if err := userService.AddUser(user); err != nil {
		EncodeError(format, w, ErrorRendererDefault(err))
		return
	}

	EncodeAnswer(FormatJSON, w, user)
}

func getAllUsers(w http.ResponseWriter, r *http.Request) {
	format := GetFormatFromRequest(r)

	users, err := userService.GetAllUsers()
	if err != nil {
		ServerErrorRender(format, w)
		return
	}

	EncodeAnswer(format, w, users, HTMLPath+"user-list.html")
}

func getUser(w http.ResponseWriter, r *http.Request) {
	format := GetFormatFromRequest(r)

	userId, err := strconv.Atoi(mux.Vars(r)[userIDKey])
	if err != nil {
		EncodeError(format, w, ErrorRendererDefault(err))
	}
	user, err := userService.GetUserById(userId)
	if err != nil {
		EncodeError(format, w, ErrorRendererDefault(err))
		return
	}

	EncodeAnswer(format, w, &userRole{user}, HTMLPath+"user-edit.html")
}

func deleteUser(w http.ResponseWriter, r *http.Request) {
	format := GetFormatFromRequest(r)

	userId, err := strconv.Atoi(mux.Vars(r)[userIDKey])
	if err != nil {
		EncodeError(format, w, ErrorRendererDefault(err))
	}
	err = userService.DeleteUser(userId)
	if err != nil {
		ServerErrorRender(format, w)
		return
	}
	EncodeAnswer(format, w, ErrorRenderer(fmt.Errorf(""), "success", http.StatusOK))
}

func updateUser(w http.ResponseWriter, r *http.Request) {
	format := GetFormatFromRequest(r)

	userId, err := strconv.Atoi(mux.Vars(r)[userIDKey])
	if err != nil {
		EncodeError(format, w, ErrorRendererDefault(err))
	}
	userData := models.User{}
	DecodeRequest(format, w, r, &userData, DecodeUserUpdateRequest)
	userData, err = userService.UpdateUser(userId, userData)
	if err != nil {
		ServerErrorRender(format, w)
		return
	}

	EncodeAnswer(format, w, &userRole{userData}, HTMLPath+"user-edit.html")
}

func DecodeUserUpdateRequest(r *http.Request, data interface{}) error  {
	r.ParseForm()
	userData := models.User{}
	userData.LoginEmail = r.FormValue("LoginEmail")
	userData.UserName = r.FormValue("UserName")
	userData.UserSurname = r.FormValue("UserSurname")
	roleId, err := strconv.Atoi(r.FormValue("RoleID"))
	if err!=nil{
		return err
	}
	userData.Role, err = userService.GetRoleById(roleId)
	if err!=nil{
		return err
	}
	data = &userData
	return nil
}
