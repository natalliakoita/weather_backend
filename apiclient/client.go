package apiclient

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"os"

	"github.com/natalliakoita/weather_backend/models"
)

type HTTPClient interface {
	Do(*http.Request) (*http.Response, error)
}

type ApiWeather struct {
	Client HTTPClient
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

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}

	resp, err := a.Client.Do(req)
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
