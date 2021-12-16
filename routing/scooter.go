package routing

import (
	"Dp218Go/models"
	"Dp218Go/services"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

var scooterService *services.ScooterService
var scooterGrpcService *services.GrpcScooterService
var orderService *services.OrderService
var scooterIDKey = "scooterId"

// TODO make dynamic access to this variables from UI
//var choosenWay = models.Coordinate{Latitude: 48.4221, Longitude: 35.0196}
var choosenWay = models.Coordinate{Latitude: 48.42543, Longitude: 35.02183} // dafi
//var choosenWay = models.Coordinate{48.42272,35.02280} // visokovoltnaya
//var choosenWay = models.Coordinate{Latitude: 48.42367 , Longitude: 35.04436} // ostapa vishni
var chosenScooterID int
var userFromRequest = models.User{ID: 1, LoginEmail: "guru_admin@guru.com", UserName: "Guru", UserSurname: "Sadh"}

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
	},
	{
		Uri:     `/start-trip`,
		Method:  http.MethodGet,
		Handler: showTripPage,
	},
	{
		Uri:     `/run`,
		Method:  http.MethodGet,
		Handler: startScooterTrip,
	},
	{
		Uri:     `/choose-scooter`,
		Method:  http.MethodPost,
		Handler: chooseScooter,
	},
	{
		Uri:     `/choose-station`,
		Method:  http.MethodPost,
		Handler: chooseStation,
	},
}

type combineForTemplate struct {
	models.ScooterListDTO
	models.StationList
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

func startScooterTrip(w http.ResponseWriter, r *http.Request) {
	statusStart, err := scooterService.CreateScooterStatusInRent(chosenScooterID)
	if err != nil {
		fmt.Println(err)
	}

	err = scooterGrpcService.InitAndRun(chosenScooterID, choosenWay)
	if err != nil {
		fmt.Println(err)
		EncodeError(FormatJSON, w, ErrorRendererDefault(err))
	}

	statusEnd, err := scooterService.CreateScooterStatusInRent(chosenScooterID)

	distance := statusEnd.Location.Distance(statusStart.Location)

	_, err = orderService.CreateOrder(userFromRequest, chosenScooterID, statusStart.ID, statusEnd.ID, distance)
	if err != nil {
		fmt.Println(err)
	}
}

func showTripPage(w http.ResponseWriter, r *http.Request) {
	scooterList, err := scooterService.GetAllScooters()
	if err != nil {
		fmt.Println(err)
	}

	stationList, err := stationService.GetAllStations()
	if err != nil {
		fmt.Println(err)
	}

	EncodeAnswer(FormatHTML, w, &combineForTemplate{*scooterList, *stationList}, HTMLPath+"scooter-run.html")
}

func chooseScooter(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		fmt.Println(err)
	}
	chosenScooterID, err = strconv.Atoi(r.Form.Get("id"))
	if err != nil {
		fmt.Println(err)
	}
}

func chooseStation(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		fmt.Println(err)
	}

	chosenStationID, err := strconv.Atoi(r.Form.Get("id"))
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(chosenStationID)
}