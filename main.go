package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/natalliakoita/weather_backend/apiclient"
	"github.com/natalliakoita/weather_backend/datasource"
	"github.com/natalliakoita/weather_backend/handlers"
	"github.com/natalliakoita/weather_backend/service"
)

func main() {
	conn, err := datasource.NewDB()
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	ds := datasource.NewDS(conn)

	c := http.Client{Timeout: time.Duration(40) * time.Second}
	a := apiclient.NewApiWeather(&c)

	apiSvc := service.NewApiService(a)
	dbSvc := service.NewDbService(&ds)

	h := handlers.NewApiHandler(&dbSvc, &apiSvc)

	router := mux.NewRouter()
	router.HandleFunc("/api/v0/{city}/weather", h.GetWeatherByCity).Methods(http.MethodGet)
	router.HandleFunc("/api/v0/weather", h.WeatherListRequest).Methods(http.MethodGet)

	log.Println("Starting API server on 8080")
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatal(err)
	}
}
