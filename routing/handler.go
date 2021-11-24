package routing

import (
	"github.com/go-chi/render"
	"github.com/gorilla/mux"
	"net/http"
)

type Route struct {
	Uri        string
	Method     string
	Handler 	func(http.ResponseWriter, *http.Request)
}

func NewRouter() *mux.Router {
	router := mux.NewRouter()
	router.MethodNotAllowedHandler = http.HandlerFunc(methodNotAllowedHandler)
	router.NotFoundHandler = http.HandlerFunc(notFoundHandler)
	return router
}

func methodNotAllowedHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(405)
	render.Render(w, r, ErrMethodNotAllowed)
}

func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(400)
	render.Render(w, r, ErrNotFound)
}
