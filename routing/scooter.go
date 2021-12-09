package routing

import (
	"Dp218Go/models"
	"Dp218Go/services"
	"fmt"
	"github.com/gorilla/mux"
	"html/template"
	"net/http"
	"strconv"
)

var scooterService *services.ScooterService
var scooterGrpcService *services.GrpcScooterService
var scooterIDKey = "scooterId"

// TODO make dynamic access to this variables from UI
var choosenWay = models.Coordinate{Latitude: 48.4423, Longitude: 35.0434}
var choosenScooter = 1

var scooterRoutes = []Route{
	{
		Uri:     `/scooters`,
		Method:  http.MethodGet,
		Handler: getAllScooters,
	},
	{
		Uri:     `/scooter/{` + scooterIDKey + `}`,
		Method:  http.MethodGet,
		Handler: getScooterById,
	}, {
		Uri:     `/run`,
		Method:  http.MethodGet,
		Handler: StartScooterTrip,
	},
}

func AddScooterHandler(router *mux.Router, service *services.ScooterService) {
	scooterService = service
	for _, rt := range scooterRoutes {
		router.Path(rt.Uri).HandlerFunc(rt.Handler).Methods(rt.Method)
		router.Path(APIprefix + rt.Uri).HandlerFunc(rt.Handler).Methods(rt.Method)
	}
}

func AddGrpcScooterHandler(router *mux.Router, service *services.GrpcScooterService) {
	scooterGrpcService = service
	for _, rt := range scooterRoutes {
		router.Path(rt.Uri).HandlerFunc(rt.Handler).Methods(rt.Method)
		router.Path(APIprefix + rt.Uri).HandlerFunc(rt.Handler).Methods(rt.Method)
	}
}

func getAllScooters(w http.ResponseWriter, r *http.Request) {

	scooters, err := scooterService.GetAllScooters()

	if err != nil {
		ServerErrorRender(FormatJSON, w)
		fmt.Println(err)
		return
	}

	EncodeAnswer(FormatJSON, w, scooters)
}

func getScooterById(w http.ResponseWriter, r *http.Request) {

	scooterID, err := strconv.Atoi(mux.Vars(r)[scooterIDKey])
	if err != nil {
		EncodeError(FormatJSON, w, ErrorRendererDefault(err))
		return
	}

	scooter, err := scooterService.GetScooterById(scooterID)
	if err != nil {
		EncodeError(FormatJSON, w, ErrorRendererDefault(err))
		return
	}

	EncodeAnswer(FormatJSON, w, scooter)
}

func StartScooterTrip(w http.ResponseWriter, r *http.Request) {
	err := scooterGrpcService.InitAndRun(choosenScooter, choosenWay)
	if err != nil {
		EncodeError(FormatJSON, w, ErrorRendererDefault(err))
		return
	}
}

func ShowTripPage(w http.ResponseWriter, r *http.Request) {

	scooterList, err := scooterService.GetAllScooters()
	if err!= nil {
		fmt.Println(err)
	}

	tmpl, err := template.ParseFiles("./templates/html/scooter-run.html")
	if err != nil {
		fmt.Println(err)
	}
	err = tmpl.Execute(w, scooterList)
	if err != nil {
		fmt.Println()
	}
}
