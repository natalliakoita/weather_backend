package service

import (
	"time"

	"github.com/natalliakoita/weather_backend/apiclient"
	"github.com/natalliakoita/weather_backend/models"
)

type ApiService struct {
	DS apiclient.ApiWeatherInterface
}

func NewApiService(d apiclient.ApiWeatherInterface) ApiService {
	return ApiService{
		DS: d,
	}
}

func (dbs *ApiService) GetWheater(city string) (models.Weather, error) {
	resp, err := dbs.DS.GetWheater(city)
	if err != nil {
		return models.Weather{}, err
	}
	var m models.Weather
	m.City = resp.Name
	tm := time.Unix(int64(resp.Dt), 0)
	m.TimeStamp = tm
	m.Temperature = float32(resp.Main.Temp)

	return m, nil
}

type ApiServiceInterface interface {
	GetWheater(city string) (models.Weather, error)
}
