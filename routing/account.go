package routing

import (
	"Dp218Go/models"
	"Dp218Go/services"
	"Dp218Go/utils"
	"errors"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

var accountService *services.AccountService
var accountIDKey = "accID"

var keyAccountRoutes = []Route{
	{
		Uri:     `/accounts`,
		Method:  http.MethodGet,
		Handler: getAllAccounts,
	},
	{
		Uri:     `/account/{` + accountIDKey + `}`,
		Method:  http.MethodGet,
		Handler: getAccountInfo,
	},
	{
		Uri:     `/account/{` + accountIDKey + `}`,
		Method:  http.MethodPost,
		Handler: updateAccountInfo,
	},
}

func AddAccountHandler(router *mux.Router, service *services.AccountService) {
	accountService = service
	accountRouter := router.NewRoute().Subrouter()
	accountRouter.Use(FilterAuth(authenticationService))

	for _, rt := range keyAccountRoutes {
		accountRouter.Path(rt.Uri).HandlerFunc(rt.Handler).Methods(rt.Method)
		accountRouter.Path(APIprefix + rt.Uri).HandlerFunc(rt.Handler).Methods(rt.Method)
	}
}

func getAllAccounts(w http.ResponseWriter, r *http.Request) {

	var accounts *models.AccountList
	var err error
	format := GetFormatFromRequest(r)

	user := GetUserFromContext(r)
	if user == nil {
		EncodeError(format, w, ErrorRendererDefault(errors.New("not authorized")))
		return
	}

	accounts, err = accountService.GetAccountsByOwner(*user)
	if err != nil {
		ServerErrorRender(format, w)
		return
	}

	EncodeAnswer(format, w, accounts, HTMLPath+"accounts.html")
}

func getAccountInfo(w http.ResponseWriter, r *http.Request) {
	format := GetFormatFromRequest(r)

	accID, err := strconv.Atoi(mux.Vars(r)[accountIDKey])
	if err != nil {
		EncodeError(format, w, ErrorRendererDefault(err))
		return
	}

	accData, err := accountService.GetAccountOutputStructByID(accID)
	if err != nil {
		EncodeError(format, w, ErrorRendererDefault(err))
		return
	}

	EncodeAnswer(format, w, accData, HTMLPath+"account.html")
}

func updateAccountInfo(w http.ResponseWriter, r *http.Request) {
	format := GetFormatFromRequest(r)

	accID, err := strconv.Atoi(mux.Vars(r)[accountIDKey])
	if err != nil {
		EncodeError(format, w, ErrorRendererDefault(err))
		return
	}
	account, err := accountService.GetAccountByID(accID)
	if err != nil {
		EncodeError(format, w, ErrorRendererDefault(err))
		return
	}
	actionType, err := GetParameterFromRequest(r, "ActionType", utils.ConvertStringToString())
	if err != nil {
		EncodeError(format, w, ErrorRendererDefault(err))
		return
	}

	switch actionType {
	case "AddMoneyToAccount":
		moneyAmount, err := GetParameterFromRequest(r, "MoneyAmount", utils.ConvertStringToFloat())
		if err != nil {
			EncodeError(format, w, ErrorRendererDefault(err))
			return
		}
		err = accountService.AddMoneyToAccount(account, int(moneyAmount.(float64)*100))
		if err != nil {
			EncodeError(format, w, ErrorRendererDefault(err))
			return
		}
	case "TakeMoneyFromAccount":
		moneyAmount, err := GetParameterFromRequest(r, "MoneyAmount", utils.ConvertStringToFloat())
		if err != nil {
			EncodeError(format, w, ErrorRendererDefault(err))
			return
		}
		err = accountService.TakeMoneyFromAccount(account, int(moneyAmount.(float64)*100))
		if err != nil {
			EncodeError(format, w, ErrorRendererDefault(err))
			return
		}
	default:
		return
	}

	getAccountInfo(w, r)
}
