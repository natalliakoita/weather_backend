package service

import (
	"errors"
	"testing"
	"time"

	"github.com/natalliakoita/weather_backend/models"
	"github.com/stretchr/testify/assert"
)

type FakeDatasource struct {
	MockAddWeatherFn            func(models.Weather) error
	MockGetListWeatherRequestFn func() ([]models.Weather, error)
}

func (fake *FakeDatasource) AddWeather(w models.Weather) error {
	return fake.MockAddWeatherFn(w)
}

func (fake *FakeDatasource) GetListWeatherRequest() ([]models.Weather, error) {
	return fake.MockGetListWeatherRequestFn()
}

func NewFakeDatasource() *FakeDatasource {
	return &FakeDatasource{
		MockAddWeatherFn: func(w models.Weather) error { return nil },
		MockGetListWeatherRequestFn: func() ([]models.Weather, error) {
			return []models.Weather{}, nil
		},
	}
}

func TestDbService_AddWeather(t *testing.T) {
	type args struct {
		w         models.Weather
		respError error
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "test1",
			args: args{
				respError: nil,
			},
			wantErr: false,
		},
		{
			name: "test2",
			args: args{
				respError: errors.New("error"),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fake := &FakeDatasource{
				MockAddWeatherFn: func(w models.Weather) error { return tt.args.respError },
			}
			dbs := &DbService{
				DS: fake,
			}
			err := dbs.AddWeather(tt.args.w)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestDbService_GetListWeatherRequest(t *testing.T) {
	type args struct {
		respError error
		respAray  []models.Weather
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		lenResp int
	}{
		{
			name: "succses call",
			args: args{
				respError: nil,
			},
			wantErr: false,
		},
		{
			name: "call with a some error",
			args: args{
				respError: errors.New("error"),
			},
			wantErr: true,
		},
		{
			name: "call when len(resp) > 0",
			args: args{
				respAray: []models.Weather{
					{
						ID:          1,
						City:        "minsk",
						TimeStamp:   time.Now(),
						Temperature: 298,
					},
				},
			},
			wantErr: false,
			lenResp: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fake := &FakeDatasource{
				MockGetListWeatherRequestFn: func() ([]models.Weather, error) {
					return tt.args.respAray, tt.args.respError
				},
			}
			dbs := &DbService{
				DS: fake,
			}
			resp, err := dbs.GetListWeatherRequest()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.Len(t, resp, tt.lenResp)
				assert.NoError(t, err)
			}
		})
	}
}
