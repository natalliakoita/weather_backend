package service

import (
	"github.com/natalliakoita/weather_backend/datasource"
	"github.com/natalliakoita/weather_backend/models"
)

type DbService struct {
	DS datasource.DSWeatherInterface
}

func NewDbService(d datasource.DSWeatherInterface) DbService {
	return DbService{
		DS: d,
	}
}

func (dbs *DbService) AddWeather(w models.Weather) error {
	err := dbs.DS.AddWeather(w)
	if err != nil {
		return err
	}
	return nil
}

func (dbs *DbService) GetListWeatherRequest() ([]models.Weather, error) {
	resp, err := dbs.DS.GetListWeatherRequest()
	if err != nil {
		return []models.Weather{}, err
	}
	// TODO
	var answer []models.Weather
	for _, w := range resp {
		var ans models.Weather
		ans.ID = w.ID
		ans.City = w.City
		ans.Temperature = w.Temperature
		ans.TimeStamp = w.TimeStamp
		answer = append(answer, ans)
	}
	return answer, nil
}

type DbServiceInterface interface {
	AddWeather(w models.Weather) error
	GetListWeatherRequest() ([]models.Weather, error)
}
