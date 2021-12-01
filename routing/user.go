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
		Uri:     `/users`,
		Method:  http.MethodPost,
		Handler: allUsersOperation,
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

type userWithRoleList struct {
	models.User
}

func (ur *userWithRoleList) ListOfRoles() []models.Role {
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

	var users = &models.UserList{}
	var err error
	format := GetFormatFromRequest(r)

	r.ParseForm()
	searchData := r.FormValue("SearchData")

	if len(searchData) == 0 {
		users, err = userService.GetAllUsers()
	} else {
		users, err = userService.FindUsersByLoginNameSurname(searchData)
	}
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
		return
	}
	user, err := userService.GetUserById(userId)
	if err != nil {
		EncodeError(format, w, ErrorRendererDefault(err))
		return
	}

	EncodeAnswer(format, w, &userWithRoleList{user}, HTMLPath+"user-edit.html")
}

func deleteUser(w http.ResponseWriter, r *http.Request) {
	format := GetFormatFromRequest(r)

	userId, err := strconv.Atoi(mux.Vars(r)[userIDKey])
	if err != nil {
		EncodeError(format, w, ErrorRendererDefault(err))
		return
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
		return
	}
	userData, err := userService.GetUserById(userId)
	if err != nil {
		EncodeError(format, w, ErrorRendererDefault(err))
		return
	}
	DecodeRequest(format, w, r, &userData, DecodeUserUpdateRequest)
	userData, err = userService.UpdateUser(userId, userData)
	if err != nil {
		ServerErrorRender(format, w)
		return
	}

	EncodeAnswer(format, w, &userWithRoleList{userData}, HTMLPath+"user-edit.html")
}

func allUsersOperation(w http.ResponseWriter, r *http.Request) {
	format := GetFormatFromRequest(r)

	r.ParseForm()
	if _, ok := r.Form["ActionType"]; !ok {

		return
	}
	actionType := r.FormValue("ActionType")
	switch actionType {
	case "BlockUser":
		userId, err := strconv.Atoi(r.FormValue("UserID"))
		if err != nil {
			EncodeError(format, w, ErrorRendererDefault(err))
			return
		}
		err = userService.ChangeUsersBlockStatus(userId)
		if err != nil {
			EncodeError(format, w, ErrorRendererDefault(err))
			return
		}
	default:
		EncodeError(format, w, ErrorRendererDefault(fmt.Errorf("unknown users operation")))
	}
	getAllUsers(w, r)
}

func DecodeUserUpdateRequest(r *http.Request, data interface{}) error {

	var err error
	r.ParseForm()
	//userData := models.User{}
	userData := data.(*models.User)

	if _, ok := r.Form["LoginEmail"]; ok {
		userData.LoginEmail = r.FormValue("LoginEmail")

	}
	if _, ok := r.Form["UserName"]; ok {
		userData.UserName = r.FormValue("UserName")
	}
	if _, ok := r.Form["UserSurname"]; ok {
		userData.UserSurname = r.FormValue("UserSurname")
	}
	if _, ok := r.Form["RoleID"]; ok {

		var roleId int
		roleId, err = strconv.Atoi(r.FormValue("RoleID"))
		if err != nil {
			return err
		}
		userData.Role, err = userService.GetRoleById(roleId)
		if err != nil {
			return err
		}
	}
	if _, ok := r.Form["IsBlocked"]; ok {
		userData.IsBlocked, _ = strconv.ParseBool(r.FormValue("IsBlocked"))
	}

	//reflect.ValueOf(data).Elem().Set(reflect.ValueOf(userData))
	return nil
}
