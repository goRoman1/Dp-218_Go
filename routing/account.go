package routing

import (
	"Dp218Go/models"
	"Dp218Go/services"
	"github.com/gorilla/mux"
	"net/http"
)

var accountService *services.AccountService
var accountIDKey = "accID"

var keyAccountRoutes = []Route{
	{
		Uri:     `/accounts`,
		Method:  http.MethodGet,
		Handler: getAllAccounts,
	},
}

func AddAccountHandler(router *mux.Router, service *services.AccountService) {
	accountService = service
	for _, rt := range keyAccountRoutes {
		router.Path(rt.Uri).HandlerFunc(rt.Handler).Methods(rt.Method)
		router.Path(APIprefix + rt.Uri).HandlerFunc(rt.Handler).Methods(rt.Method)
	}
}

func getAllAccounts(w http.ResponseWriter, r *http.Request) {

	var accounts = &models.AccountList{}
	var err error
	format := GetFormatFromRequest(r)

	user, err := AuthService.GetUserFromRequest(r)
	if err != nil {
		EncodeError(FormatHTML, w, ErrorRendererDefault(err))
		return
	}

	accounts, err = accountService.GetAccountsByOwner(*user)
	if err != nil {
		ServerErrorRender(format, w)
		return
	}

	//TODO: fill accounting page
	EncodeAnswer(format, w, accounts, HTMLPath+"accounting.html")
}