package routing

import (
	"Dp218Go/configs"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"html/template"
	"net/http"
	"strings"
)

type Route struct {
	Uri     string
	Method  string
	Handler func(http.ResponseWriter, *http.Request)
}

const (
	FormatJSON = iota
	FormatHTML
)

var (
	HTMLPath      = configs.TEMPLATES_PATH + "html/"
	MainPageHTML  = HTMLPath + "main-page.html"
	ErrorPageHTML = HTMLPath + "error.html"
	APIprefix     = "/api/v1"
)

func NewRouter() *mux.Router {
	router := mux.NewRouter()
	router.MethodNotAllowedHandler = http.HandlerFunc(methodNotAllowedHandler)
	router.NotFoundHandler = http.HandlerFunc(notFoundHandler)
	router.PathPrefix("/templates/").Handler(http.StripPrefix("/templates/",
		http.FileServer(http.Dir(configs.TEMPLATES_PATH))))
	router.HandleFunc("/", showHomePage)
	router.HandleFunc("/login", showLoginPage)
	return router
}

func showHomePage(w http.ResponseWriter, r *http.Request) {
	EncodeAnswer(FormatHTML, w, nil, MainPageHTML)
}

func showLoginPage(w http.ResponseWriter, r *http.Request) {
	EncodeAnswer(FormatHTML, w, nil, HTMLPath + "login-registration.html")
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
		if tmpl, err = template.ParseFiles(ErrorPageHTML); err == nil {
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
		if len(htmlTemplates) == 0 {
			EncodeError(format, w, &ResponseStatus{StatusText: "Encode error", Message: "no html temlates", StatusCode: http.StatusInternalServerError})
			return
		}
		w.Header().Set("Content-Type", "text/html")
		var tmpl *template.Template
		if tmpl, err = template.ParseFiles(htmlTemplates[0]); err == nil {
			err = tmpl.Execute(w, answer)
		}
	default:
		err = fmt.Errorf("format error")
	}

	if err != nil {
		EncodeError(format, w, &ResponseStatus{Err: err, StatusText: "Encode error", Message: err.Error(), StatusCode: http.StatusInternalServerError})
		return
	}

	w.WriteHeader(http.StatusOK)
}

func DecodeRequest(format int, w http.ResponseWriter, r *http.Request, requestData interface{}, htmlDecoder func(r *http.Request, dataToDecode interface{}) error) {
	var err error
	switch format {
	case FormatJSON:
		w.Header().Set("Content-Type", "application/json")
		err = json.NewDecoder(r.Body).Decode(requestData)
	case FormatHTML:
		w.Header().Set("Content-Type", "text/html")
		err = htmlDecoder(r, requestData)
	default:
		err = fmt.Errorf("format error")
	}

	if err != nil {
		EncodeError(FormatJSON, w, ErrorRenderer(err, "Bad request", http.StatusBadRequest))
		return
	}

	w.WriteHeader(http.StatusOK)
}

func GetFormatFromRequest(r *http.Request) int {
	if strings.Contains(r.RequestURI, APIprefix) {
		return FormatJSON
	}
	return FormatHTML
}
