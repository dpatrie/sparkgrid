package services

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/mademediacorp/onboard/gae/api/storage"
)

type route struct {
	method  string
	pattern string
	handler http.HandlerFunc
}

func newRouter(routes []route) http.Handler {
	router := httprouter.New()
	router.RedirectTrailingSlash = false
	router.RedirectFixedPath = false

	for _, r := range routes {
		if r.method == "POST" || r.method == "PUT" {
			r.handler = withRequestBody(r.handler)
		}

		router.HandlerFunc(r.method, r.pattern, r.handler)
	}
	return router
}

func withRequestBody(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Body == nil {
			http.Error(w, "Missing request body", 400)
			return
		}
		defer r.Body.Close()
		next(w, r)
	}
}

var (
	ErrUnauthorized       = errors.New("unauthorized")
	ErrPreconditionFailed = errors.New("precondition failed")
	ErrBadRequest         = errors.New("bad request")
	ErrNotFound           = errors.New("not found")
	ErrNotModified        = errors.New("not modified")
)

func Error(w http.ResponseWriter, err error) {
	switch err {
	case ErrNotFound:
		http.Error(w, err.Error(), http.StatusNotFound)
	case ErrUnauthorized:
		http.Error(w, err.Error(), http.StatusUnauthorized)
	case ErrBadRequest:
		http.Error(w, err.Error(), http.StatusBadRequest)
	case ErrNotModified:
		http.Error(w, err.Error(), http.StatusNotModified)
	case ErrPreconditionFailed:
		http.Error(w, err.Error(), http.StatusPreconditionFailed)
	case storage.ErrAlreadyExists:
		http.Error(w, err.Error(), http.StatusConflict)
	default:
		log.Println(err.Error())
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}

func JSONResponse(w http.ResponseWriter, v interface{}) {
	w.Header().Set("Content-type", "application/json")
	json.NewEncoder(w).Encode(v)
}

func urlParam(r *http.Request, key string) string {
	return httprouter.ParamsFromContext(r.Context()).ByName(key)
}
