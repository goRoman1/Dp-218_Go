package routing

import (
	iface "Dp218Go/domain/interfaces"
	repo "Dp218Go/repositories"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"

	"github.com/go-chi/render"
)

var	scooterRepo iface.ScooterRepo
var scooterIdKey = "scooterId"

var scooterRoutes = []Route{
	{
		Uri:         `/scooters`,
		Method:    http.MethodGet,
		Handler:	getAllScooters,
	},
	{
		Uri:         `/scoot/{`+scooterIdKey+`}`,
		Method:     http.MethodGet,
		Handler:	getScooterById,
	},
}

func AddScooterHandler(router *mux.Router, repo iface.ScooterRepo) {
	scooterRepo = repo
	for _, rt := range scooterRoutes{
		router.Path(rt.Uri).HandlerFunc(rt.Handler).Methods(rt.Method)
	}
}

func getAllScooters(w http.ResponseWriter, r *http.Request) {
	scooters, err := scooterRepo.GetAllScooters()
	if err != nil {
		render.Render(w, r, ServerErrorRenderer(err))
		return
	}
	if err := render.Render(w, r, scooters); err != nil {
		render.Render(w, r, ErrorRenderer(err))
		return
	}
}

func getScooterById(w http.ResponseWriter, r *http.Request) {
	scooterId, err := strconv.Atoi(mux.Vars(r)[scooterIdKey])
	if err != nil {
		render.Render(w, r, ErrorRenderer(err))
	}
	scooter, err := scooterRepo.GetScooterById(scooterId)
	if err != nil {
		if err == repo.ErrNoMatch {
			render.Render(w, r, ErrNotFound)
		} else {
			render.Render(w, r, ErrorRenderer(err))
		}
		return
	}
	if err := render.Render(w, r, &scooter); err != nil {
		render.Render(w, r, ServerErrorRenderer(err))
		return
	}
}

