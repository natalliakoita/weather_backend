package handlers

import (
	"encoding/json"
	"net/http"
	"time"

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
	city := "Barcelona"

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
		ID:          resp.ID,
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
	ID          int       `json:"id"`
	City        string    `json:"city"`
	TimeStamp   time.Time `json:"dt"`
	Temperature float32   `json:"temperature"`
}

// write web response to client
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