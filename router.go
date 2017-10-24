package main

import (
	"encoding/json"
	//"fmt"
	//"github.rackspace.com/backbone/token-go.git"
	"log"
	"net/http"
	"strings"
)

type AppContext struct {
	Body       []byte
	StatusCode int
	//Web
	Url   string
	Err   error
	Mongo *Mongo
}

func NewAppContext() *AppContext {
	return &AppContext{}
}

type AppHandler struct {
	*AppContext
	H func(*AppContext, http.ResponseWriter, *http.Request) (*jsonData, *appError)
}

type appError struct {
	Code    int
	Message string
	Err     string
}

func (ah AppHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	allowedMethods := []string{
		"POST",
		"GET",
		"OPTIONS",
		"PUT",
		"DELETE",
		"PATCH",
	}

	allowHeaders := []string{
		"Accept",
		"Content-Type",
		"Content-Length",
		"Accept-Encoding",
		"Authorization",
		"X-CRSF-Token",
		"Authorization Bearer",
	}

	if origin := r.Header.Get("Origin"); origin != "" {
		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set(
			"Access-Control-Allow-Methods",
			strings.Join(allowedMethods, ", "))

		w.Header().Set(
			"Access-Control-Allow-Headers",
			strings.Join(allowHeaders, ", "))
	}
	if r.Method == "OPTIONS" {
		return
	}

	if j, err := ah.H(ah.AppContext, w, r); err != nil {
		log.Printf("%#v", err)
		serveError(w, r, err)
	} else {
		serveJSON(w, j)
	}
}

func use(app AppHandler, middleware ...func(http.Handler) http.Handler) http.Handler {
	var response http.Handler = app
	for _, m := range middleware {
		response = m(response)
	}

	return response
}

// func authToken(h http.Handler) http.Handler {
//         return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//                 jtoken := jwtToken.New()
//                 t, err := jtoken.ParseTokenFromRequest(r)
//                 if err == nil && t.Valid {
//                         h.ServeHTTP(w, r)
//                 } else {
//                         err = fmt.Errorf("Not Authorized: %v", err)
//                         log.Printf("%#v\n", err)
//                         serveError(w, r, &appError{401, "Not Authorized; need valid JWT token.", err.Error()})
//                 }
//         })
// }

func serveError(w http.ResponseWriter, r *http.Request, err *appError) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(err.Code)
	b, _ := json.MarshalIndent(err, "", "\t")
	w.Write(b)
}

func serveJSON(w http.ResponseWriter, j *jsonData) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(j.Code)
	w.Write(j.Byte)
}
