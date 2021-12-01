package routing

import (
	"Dp218Go/models"
	"fmt"
	"net/http"
	"strconv"

	"Dp218Go/services"
	"github.com/gorilla/mux"
)

var stationService *services.StationService
var stationIDKey = "stationID"

var keyRoutesStation = []Route{
	{
		Uri:     `/stations`,
		Method:  http.MethodGet,
		Handler: getAllStations,
	},
	{
		Uri:     `/station/{` + stationIDKey + `}`,
		Method:  http.MethodGet,
		Handler: getStation,
	},
	{
		Uri:     `/station`,
		Method:  http.MethodPost,
		Handler: createStation,
	},
	{
		Uri:     `/station/{` + stationIDKey + `}`,
		Method:  http.MethodDelete,
		Handler: deleteStation,
	},
}

func AddStationHandler(router *mux.Router, service *services.StationService) {
	stationService = service
	for _, rt := range keyRoutesStation {
		router.Path(rt.Uri).HandlerFunc(rt.Handler).Methods(rt.Method)
		router.Path(APIprefix + rt.Uri).HandlerFunc(rt.Handler).Methods(rt.Method)
	}
}

func createStation(w http.ResponseWriter, r *http.Request) {
	format := GetFormatFromRequest(r)

	station := &models.Station{}
	DecodeRequest(format, w, r, station, nil)

	if err := stationService.AddStation(station); err != nil {
		EncodeError(format, w, ErrorRendererDefault(err))
		return
	}

	EncodeAnswer(format, w, station)
}

func getAllStations(w http.ResponseWriter, r *http.Request) {
	var station = &models.StationList{}
	var err error
	format := GetFormatFromRequest(r)

		station, err = stationService.GetAllStations()
	if err != nil {
		ServerErrorRender(format, w)
		return
	}

	EncodeAnswer(format, w, station, HTMLPath+"user-list.html")
}

func getStation(w http.ResponseWriter, r *http.Request) {
	format := GetFormatFromRequest(r)

	stationId, err := strconv.Atoi(mux.Vars(r)[stationIDKey])
	if err != nil {
		EncodeError(format, w, ErrorRendererDefault(err))
		return
	}
	station, err := stationService.GetStationById(stationId)
	if err != nil {
		EncodeError(format, w, ErrorRendererDefault(err))
		return
	}

	EncodeAnswer(format, w, station, HTMLPath+"station-edit.html")
}

func deleteStation(w http.ResponseWriter, r *http.Request) {
	format := GetFormatFromRequest(r)

	stationId, err := strconv.Atoi(mux.Vars(r)[stationIDKey])
	if err != nil {
		EncodeError(format, w, ErrorRendererDefault(err))
		return
	}
	err = stationService.DeleteStation(stationId)
	if err != nil {
		ServerErrorRender(format, w)
		return
	}
	EncodeAnswer(format, w, ErrorRenderer(fmt.Errorf(""), "success", http.StatusOK))
}
