package apiclient

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"os"

	"github.com/natalliakoita/weather_backend/models"
)

type ApiWeather struct {
	Client *http.Client
	Key    string
}

func NewApiWeather(conn *http.Client) ApiWeather {
	key := os.Getenv("KEY")
	return ApiWeather{
		Client: conn,
		Key:    key,
	}
}

func (a ApiWeather) GetWheater(city string) (*models.WheatherApiResponse, error) {
	baseUrl := "http://api.openweathermap.org/data/2.5/weather"
	u, err := url.Parse(baseUrl)
	if err != nil {
		return nil, err
	}
	q := u.Query()
	q.Set("q", city)
	q.Set("appid", a.Key)
	u.RawQuery = q.Encode()

	resp, err := a.Client.Get(u.String())
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	bodyData, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	w := new(models.WheatherApiResponse)
	if err = json.Unmarshal(bodyData, w); err != nil {
		return nil, err
	}
	return w, nil
}

type ApiWeatherInterface interface {
	GetWheater(city string) (*models.WheatherApiResponse, error)
}
