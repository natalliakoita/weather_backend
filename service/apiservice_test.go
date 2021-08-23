package service

import (
	"errors"
	"testing"

	"github.com/natalliakoita/weather_backend/models"
	"github.com/stretchr/testify/assert"
)

type FakeApiclient struct {
	MockGetWheaterFn func(city string) (*models.WheatherApiResponse, error)
}

func (fake *FakeApiclient) GetWheater(city string) (*models.WheatherApiResponse, error) {
	return fake.MockGetWheaterFn(city)
}

func TestApiService_GetWheater(t *testing.T) {
	type args struct {
		w         *models.WheatherApiResponse
		city      string
		respError error
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		lenResp int
	}{
		{
			name: "call when len(resp) > 0",
			args: args{
				city: "Minsk",
				w: &models.WheatherApiResponse{
					Dt: 1,
					Main: models.Main{
						Temp: 298,
					},
					Sys: models.Sys{
						Country: "Minsk",
					},
				},
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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fake := &FakeApiclient{
				MockGetWheaterFn: func(string) (*models.WheatherApiResponse, error) { return tt.args.w, tt.args.respError },
			}
			dbs := &ApiService{
				DS: fake,
			}
			got, err := dbs.GetWheater(tt.args.city)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.IsType(t, models.Weather{}, got)
				assert.NoError(t, err)
			}
		})
	}
}
