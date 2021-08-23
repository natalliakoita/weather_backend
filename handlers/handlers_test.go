package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gorilla/mux"
	"github.com/natalliakoita/weather_backend/models"
	"github.com/stretchr/testify/assert"
)

// Example: https://gist.github.com/p4tin/64cf99e30b034ee100b9

type FakeApiSVC struct {
	MockGetWheaterFn func(string) (models.Weather, error)
}

func (fake *FakeApiSVC) GetWheater(city string) (models.Weather, error) {
	return fake.MockGetWheaterFn(city)
}

type FakeDbSVC struct {
	MockAddWeatherFn            func(models.Weather) error
	MockGetListWeatherRequestFn func() ([]models.Weather, error)
}

func (fake *FakeDbSVC) AddWeather(w models.Weather) error {
	return fake.MockAddWeatherFn(w)
}

func (fake *FakeDbSVC) GetListWeatherRequest() ([]models.Weather, error) {
	return fake.MockGetListWeatherRequestFn()
}

func MockWeatherModel() models.Weather {
	var w models.Weather
	w.ID = 1
	w.City = "Minsk"
	w.TimeStamp = time.Now()
	w.Temperature = 100
	return w
}

func TestApiHandler_GetWeatherByCity(t *testing.T) {
	type args struct {
		city        string
		testWeather models.Weather
		dbErr       error
		apiErr      error
	}
	tests := []struct {
		name       string
		args       args
		wantErr    bool
		statusCode int
	}{
		{
			name: "test 1",
			args: args{
				city:        "Minsk",
				testWeather: MockWeatherModel(),
				dbErr:       nil,
				apiErr:      nil,
			},
			wantErr:    false,
			statusCode: http.StatusOK,
		},
		{
			name: "test 2",
			args: args{
				city:   "Minsk",
				apiErr: errors.New("a some err"),
			},
			wantErr:    true,
			statusCode: http.StatusInternalServerError,
		},
		{
			name: "test 3",
			args: args{
				city:   "Minsk",
				testWeather: MockWeatherModel(),
				dbErr: errors.New("a some err"),
			},
			wantErr:    true,
			statusCode: http.StatusInternalServerError,
		},
		{
			name: "test 4",
			args: args{
				city:   "",
			},
			wantErr:    true,
			statusCode: http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeDbSVC := &FakeDbSVC{
				MockAddWeatherFn: func(models.Weather) error { return tt.args.dbErr },
				// MockGetListWeatherRequestFn: func() ([]models.Weather, error) { return tt.args.weatherList, tt.args.dbErr },
			}

			fakeApiSVC := &FakeApiSVC{
				MockGetWheaterFn: func(string) (models.Weather, error) { return tt.args.testWeather, tt.args.apiErr },
			}

			u := &ApiHandler{
				dbSvc:  fakeDbSVC,
				apiSvc: fakeApiSVC,
			}

			path := fmt.Sprintf("/api/v0/%s/weather", tt.args.city)
			req, err := http.NewRequest("GET", path, nil)
			if err != nil {
				t.Fatal(err)
			}
			rr := httptest.NewRecorder()

			// Example: https://stackoverflow.com/questions/34435185/unit-testing-for-functions-that-use-gorilla-mux-url-parameters
			vars := map[string]string{}
			if tt.args.city != "" {
				vars["city"] = tt.args.city
			}

			req = mux.SetURLVars(req, vars)

			handler := http.HandlerFunc(u.GetWeatherByCity)
			handler.ServeHTTP(rr, req)

			assert.Equal(t, rr.Code, tt.statusCode)
			b, err := getBody(tt.args.testWeather)
			assert.NoError(t, err)
			if !tt.wantErr {
				assert.Equal(t, b, rr.Body.Bytes())
			}
		})
	}
}

func getBody(w models.Weather) ([]byte, error) {
	result := GetWeatherResponse{
		City:        w.City,
		TimeStamp:   w.TimeStamp,
		Temperature: w.Temperature,
	}
	b, err := json.Marshal(result)
	if err != nil {
		return []byte{}, err
	}
	return b, nil
}

func TestApiHandler_WeatherListRequest(t *testing.T) {
	type args struct {
		weatherList []models.Weather
		dbErr       error
	}
	tests := []struct {
		name       string
		args       args
		wantErr    bool
		statusCode int
	}{
		{
			name: "test 1",
			args: args{
				dbErr:       nil,
				weatherList: []models.Weather{
					MockWeatherModel(),
				},
			},
			wantErr:    false,
			statusCode: http.StatusOK,
		},
		{
			name: "test 2",
			args: args{
				dbErr: errors.New("a some err"),
			},
			wantErr:    true,
			statusCode: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeDbSVC := &FakeDbSVC{
				MockGetListWeatherRequestFn: func() ([]models.Weather, error) { return tt.args.weatherList, tt.args.dbErr },
			}

			fakeApiSVC := &FakeApiSVC{}

			u := &ApiHandler{
				dbSvc:  fakeDbSVC,
				apiSvc: fakeApiSVC,
			}

			path := "/api/v0/weather"
			req, err := http.NewRequest("GET", path, nil)
			if err != nil {
				t.Fatal(err)
			}
			rr := httptest.NewRecorder()

			handler := http.HandlerFunc(u.WeatherListRequest)
			handler.ServeHTTP(rr, req)

			assert.Equal(t, rr.Code, tt.statusCode)
			b, err := getArrayBody(tt.args.weatherList)
			assert.NoError(t, err)
			if !tt.wantErr {
				assert.Equal(t, b, rr.Body.Bytes())
			}
		})
	}
}

func getArrayBody(w []models.Weather) ([]byte, error) {
	response := []GetWeatherResponse{}
	for _, city := range w {
		q := GetWeatherResponse{}
		q.City = city.City
		q.TimeStamp = city.TimeStamp
		q.Temperature = city.Temperature

		response = append(response, q)
	}

	result := WeatherListResponse{}
	result.Cities = response

	b, err := json.Marshal(result)
	if err != nil {
		return []byte{}, err
	}
	return b, nil
}