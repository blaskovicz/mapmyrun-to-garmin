package main

import (
	"fmt"
	"net/http"
	"os"

	swarmed "github.com/blaskovicz/go-swarmed"
	"github.com/blaskovicz/mapmyrun-to-garmin/web"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

func cors(router http.Handler) http.Handler {
	if os.Getenv("ENVIRONMENT") != "development" {
		return router
	}
	return handlers.CORS(
		handlers.AllowedHeaders([]string{"Authorization", "Content-Type", "X-Requested-With"}),
		handlers.AllowedOrigins([]string{"*"}),
		handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PATCH", "DELETE", "PUT", "OPTIONS"}),
	)(router)
}

func main() {
	err := swarmed.LoadSecrets()
	if err != nil {
		panic(fmt.Errorf("swarmed.LoadSecrets: %s", err))
	}

	router := mux.NewRouter()
	router.HandleFunc("/api/garmin/import", web.ApiPostGarminImport).Methods("POST")
	router.NotFoundHandler = http.FileServer(http.Dir("./dist"))
	listen := os.Getenv("PORT")
	if listen == "" {
		listen = "3000"
	}
	listen = ":" + listen
	logrus.Infof("Starting on %s", listen)
	logrus.Fatal(http.ListenAndServe(listen, cors(router)))
}
