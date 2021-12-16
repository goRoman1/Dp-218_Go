package routing

import (
	"Dp218Go/models"
	"Dp218Go/services"
	"encoding/json"
	"net/http"
	"strconv"
	"text/template"

	"github.com/gorilla/mux"
)

var customerService *services.CustomerService
var custTmpl = template.Must(template.ParseFiles("templates/html/customer-home.html"))

func AddCustomerHandler(router *mux.Router, service *services.CustomerService) {
	customerService = service

	custRouter := router.PathPrefix("/customer").Subrouter()
	custRouter.Use(authenticationService.FilterAuth, FilterCustomer)

	custRouter.Path("/home").HandlerFunc(CustomerHomeHandler).Methods(http.MethodGet)
	custRouter.Path("/station").HandlerFunc(CustomerStationListHandler).Methods(http.MethodGet)
	custRouter.Path("/station/nearest").
		HandlerFunc(CustomerStationNearestHandler).Queries("x", "{x}", "y", "{y}").Methods(http.MethodGet)
	custRouter.Path("/station/{id:[0-9]+}").HandlerFunc(CustomerStationInfoHandler).Methods(http.MethodGet)

}

func CustomerHomeHandler(w http.ResponseWriter, r *http.Request) {
	user := services.GetUserFromContext(r)
	// no need if wrapped with filteruser
	if user == nil {
		http.Error(w, "not authenticated", http.StatusForbidden)
		return
	}

	custTmpl.ExecuteTemplate(w, "customer-home.html", user)
}

func CustomerStationListHandler(w http.ResponseWriter, r *http.Request) {
	// TODO show only not blocked
	sts, err := customerService.ListStations()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(sts.Station)
}

func CustomerStationNearestHandler(w http.ResponseWriter, r *http.Request) {

	xStr := r.FormValue("x")
	yStr := r.FormValue("y")

	x, err := strconv.ParseFloat(xStr, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotAcceptable)
		return
	}

	y, err := strconv.ParseFloat(yStr, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotAcceptable)
		return
	}

	nearest, err := customerService.ShowNearestStation(x, y)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotAcceptable)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode([]*models.Station{nearest})
}

func CustomerStationInfoHandler(w http.ResponseWriter, r *http.Request) {

	idStr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotAcceptable)
		return
	}
	station, err := customerService.ShowStation(id)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(station)

}

func FilterCustomer(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := services.GetUserFromContext(r)
		if user == nil || !(user.Role.IsUser || user.Role.IsAdmin) {
			http.Error(w, "only customers allowed", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}
