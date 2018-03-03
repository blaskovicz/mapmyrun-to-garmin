package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	swarmed "github.com/blaskovicz/go-swarmed"
	"github.com/blaskovicz/mapmyrun-to-garmin/web"
	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
)

func main() {
	err := swarmed.LoadSecrets()
	if err != nil {
		panic(fmt.Errorf("swarmed.LoadSecrets: %s", err))
	}
	env := os.Getenv("ENVIRONMENT")
	var secureCSRF bool
	if env == "production" {
		secureCSRF = true
	}
	sessionStorage := sessions.NewCookieStore([]byte(os.Getenv("COOKIE_KEY")))
	router := mux.NewRouter()
	router.HandleFunc("/", web.Index).Methods("GET")
	router.HandleFunc("/routes/new", web.NewRouteForm(sessionStorage)).Methods("GET")
	router.HandleFunc("/routes/new", web.PostRouteForm(sessionStorage)).Methods("POST")
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	log.Fatal(http.ListenAndServe(":"+port, csrf.Protect([]byte(os.Getenv("CSRF_KEY")), csrf.Secure(secureCSRF))(router)))
}
