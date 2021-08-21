package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/natalliakoita/weather_backend/service"
)

type ApiHandler struct {
	dbSvc  service.DbServiceInterface
	apiSvc service.ApiServiceInterface
}

func NewApiHandler(d service.DbServiceInterface, a service.ApiServiceInterface) ApiHandler {
	cont := ApiHandler{
		dbSvc:  d,
		apiSvc: a,
	}
	return cont
}

func (u *ApiHandler) GetWeatherByCity(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	city, ok := vars["city"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	resp, err := u.apiSvc.GetWheater(city)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = u.dbSvc.AddWeather(resp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	result := GetWeatherResponse{
		City:        resp.City,
		TimeStamp:   resp.TimeStamp,
		Temperature: resp.Temperature,
	}
	err = result.writeToWeb(w)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

type GetWeatherResponse struct {
	City        string    `json:"city"`
	TimeStamp   time.Time `json:"dt"`
	Temperature float32   `json:"temperature"`
}

func (c GetWeatherResponse) writeToWeb(w http.ResponseWriter) error {
	b, err := json.Marshal(c)
	if err != nil {
		return err
	}
	w.Header().Add("Content-Type", "application/json")
	if _, err := w.Write(b); err != nil {
		return err
	}
	return nil
}

type WeatherListResponse struct {
	Cities []GetWeatherResponse `json:"cities"`
}

func (c WeatherListResponse) writeToWeb(w http.ResponseWriter) {
	b, err := json.Marshal(c)
	if err != nil {
		log.Fatal(err)
	}
	w.Header().Add("Content-Type", "application/json")
	if _, err := w.Write(b); err != nil {
		log.Fatal(err)
	}
}

func (u *ApiHandler) WeatherListRequest(w http.ResponseWriter, req *http.Request) {
	cities, err := u.dbSvc.GetListWeatherRequest()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	response := []GetWeatherResponse{}
	for _, city := range cities {
		q := GetWeatherResponse{}
		q.City = city.City
		q.TimeStamp = city.TimeStamp
		q.Temperature = city.Temperature

		response = append(response, q)
	}

	resp := WeatherListResponse{}
	resp.Cities = response

	resp.writeToWeb(w)

	w.WriteHeader(http.StatusOK)
}
