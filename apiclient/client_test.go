package apiclient

import (
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/natalliakoita/weather_backend/models"
	"github.com/stretchr/testify/assert"
)

// Example: https://dev.to/clavinjune/mocking-http-call-in-golang-15i5

type HTTPClientMock struct {
	DoFunc func(*http.Request) (*http.Response, error)
}

func (H HTTPClientMock) Do(r *http.Request) (*http.Response, error) {
	return H.DoFunc(r)
}

func TestApiWeather_GetWheater(t *testing.T) {
	type fields struct {
		Key string
	}
	type args struct {
		city     string
		httpResp *http.Response
		httpErr  error
	}
	tests := []struct {
		name       string
		fields     fields
		Body       string
		StatusCode int
		args       args
		want       *models.WheatherApiResponse
		wantErr    bool
	}{
		{
			name: "test 1",
			fields: fields{
				Key: "test",
			},
			Body: `{
				"coord": {
				  "lon": -0.1257,
				  "lat": 51.5085
				},
				"weather": [
				  {
					"id": 802,
					"main": "Clouds",
					"description": "scattered clouds",
					"icon": "03d"
				  }
				],
				"base": "stations",
				"main": {
				  "temp": 293.55,
				  "feels_like": 293.29,
				  "temp_min": 291.48,
				  "temp_max": 295.2,
				  "pressure": 1011,
				  "humidity": 63
				},
				"visibility": 10000,
				"wind": {
				  "speed": 6.69,
				  "deg": 240
				},
				"clouds": {
				  "all": 26
				},
				"dt": 1629047290,
				"sys": {
				  "type": 2,
				  "id": 2006068,
				  "country": "GB",
				  "sunrise": 1629002776,
				  "sunset": 1629055436
				},
				"timezone": 3600,
				"id": 2643743,
				"name": "London",
				"cod": 200
			  }`,
			StatusCode: 200,
			wantErr:    false,
		},
		{
			name: "some error",
			fields: fields{
				Key: "test",
			},
			StatusCode: 500,
			wantErr:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := &HTTPClientMock{}
			client.DoFunc = func(r *http.Request) (*http.Response, error) {
				return &http.Response{
					Body:       io.NopCloser(strings.NewReader(tt.Body)),
					StatusCode: tt.StatusCode,
				}, nil
			}

			a := ApiWeather{
				Client: client,
				Key:    tt.fields.Key,
			}
			got, err := a.GetWheater(tt.args.city)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, got.Timezone)
			}
		})
	}
}
