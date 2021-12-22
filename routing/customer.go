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

var custTmpl = template.Must(template.ParseFiles(HTMLPath + "customer-home.html"))

type customerHandler struct {
	custService *services.CustomerService
}

func newCustomerHandler(service *services.CustomerService) *customerHandler {
	return &customerHandler{
		custService: service,
	}
}

//AddCustomerHandler registeres endpoints for customer
func AddCustomerHandler(router *mux.Router, service *services.CustomerService) {

	custHandler := newCustomerHandler(service)

	custRouter := router.PathPrefix("/customer").Subrouter()
	custRouter.Use(FilterAuth(authenticationService), FilterCustomer)

	custRouter.Path("/home").HandlerFunc(custHandler.HomeHandler).Methods(http.MethodGet)
	custRouter.Path("/station").HandlerFunc(custHandler.StationListHandler).Methods(http.MethodGet)
	custRouter.Path("/station/nearest").
		HandlerFunc(custHandler.StationNearestHandler).Queries("x", "{x}", "y", "{y}").Methods(http.MethodGet)
	custRouter.Path("/station/{id:[0-9]+}").HandlerFunc(custHandler.StationInfoHandler).Methods(http.MethodGet)

}

// HomeHandler is handler for rendering home page of customer
func (h *customerHandler) HomeHandler(w http.ResponseWriter, r *http.Request) {
	user := GetUserFromContext(r)
	// no need if wrapped with filteruser
	if user == nil {
		http.Error(w, "not authenticated", http.StatusForbidden)
		return
	}

	custTmpl.ExecuteTemplate(w, "customer-home.html", user)
}

// StationListHandler is handler that users customer service
// shows list of available stations on map
// returns json station list in response shows error if failed
func (h *customerHandler) StationListHandler(w http.ResponseWriter, r *http.Request) {
	// TODO show only not blocked
	sts, err := h.custService.ListStations()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(sts.Station)
}

//StationNearestHandler is handler that user customer service
// takes user location and returns nearest
// station in json format shows error if failed
func (h *customerHandler) StationNearestHandler(w http.ResponseWriter, r *http.Request) {

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

	nearest, err := h.custService.ShowNearestStation(x, y)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotAcceptable)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode([]*models.Station{nearest})
}

// StationInfoHandler is handler that shows general station info of station
// by station id received in reguest url var
func (h *customerHandler) StationInfoHandler(w http.ResponseWriter, r *http.Request) {

	idStr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotAcceptable)
		return
	}
	station, err := h.custService.ShowStation(id)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(station)

}

// FilterCustomer is middleware that restricts access to customer page
// checks if user role is customer or admin
// shows error if not allowed
func FilterCustomer(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := GetUserFromContext(r)
		if user == nil || !(user.Role.IsUser || user.Role.IsAdmin) {
			http.Error(w, "only customers allowed", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}
