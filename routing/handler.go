package routing

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"html/template"
	"net/http"
)

type Route struct {
	Uri     string
	Method  string
	Handler func(http.ResponseWriter, *http.Request)
}

const (
	FormatJSON = iota
	FormatHTML
	HTMLforError = "file:///home/Dp218Go/templates/html/error.html"
)

func NewRouter() *mux.Router {
	router := mux.NewRouter()
	router.MethodNotAllowedHandler = http.HandlerFunc(methodNotAllowedHandler)
	router.NotFoundHandler = http.HandlerFunc(notFoundHandler)
	return router
}

func methodNotAllowedHandler(w http.ResponseWriter, r *http.Request) {
	EncodeError(FormatHTML, w, ErrorRenderer(fmt.Errorf("method not allowed"), "Not allowed", http.StatusMethodNotAllowed))
}

func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	EncodeError(FormatHTML, w, ErrorRenderer(fmt.Errorf("resource not found"), "Not found", http.StatusNotFound))
}

func ServerErrorRender(format int, w http.ResponseWriter) {
	EncodeError(format, w, ErrorRenderer(fmt.Errorf("server error"), "Internal server error", http.StatusInternalServerError))
}

func EncodeError(format int, w http.ResponseWriter, respErr *ResponseStatus) {
	var err error
	switch format {
	case FormatJSON:
		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(respErr)
	case FormatHTML:
		w.Header().Set("Content-Type", "text/html")
		var tmpl *template.Template
		if tmpl, err = template.ParseFiles(HTMLforError); err == nil {
			err = tmpl.Execute(w, respErr)
		}
	default:
		err = fmt.Errorf("format error")
	}

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(respErr.StatusCode)
}

func EncodeAnswer(format int, w http.ResponseWriter, answer interface{}, htmlTemplates ...string) {
	var err error

	switch format {
	case FormatJSON:
		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(answer)
	case FormatHTML:
		w.Header().Set("Content-Type", "text/html")
		var tmpl *template.Template
		if tmpl, err = template.ParseFiles(htmlTemplates[0]); err == nil {
			err = tmpl.Execute(w, answer)
		}
	default:
		err = fmt.Errorf("format error")
	}

	if err != nil {
		ServerErrorRender(format, w)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func DecodeRequest(format int, w http.ResponseWriter, r *http.Request, requestData interface{}) {
	var err error
	switch format {
	case FormatJSON:
		w.Header().Set("Content-Type", "application/json")
		err = json.NewDecoder(r.Body).Decode(requestData)
	case FormatHTML:
		w.Header().Set("Content-Type", "text/html")
		//TODO: make decode from html forms work
	default:
		err = fmt.Errorf("format error")
	}

	if err != nil {
		EncodeError(FormatJSON, w, ErrorRenderer(err, "Bad request", http.StatusBadRequest))
		return
	}

	w.WriteHeader(http.StatusOK)
}
