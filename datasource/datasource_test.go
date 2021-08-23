package datasource

import (
	"database/sql"
	"errors"
	"log"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	_ "github.com/lib/pq"
	"github.com/natalliakoita/weather_backend/models"
	"github.com/stretchr/testify/assert"
)

func NewMock() (*sql.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	return db, mock
}

var w = Weather{
	ID:          1,
	City:        "Minsk",
	TimeStamp:   time.Now(),
	Temperature: 298,
}

var wm = models.Weather{
	ID:          1,
	City:        "Minsk",
	TimeStamp:   time.Now(),
	Temperature: 298,
}

func TestAddWeather(t *testing.T) {
	db, mock := NewMock()
	defer db.Close()

	ds := NewDS(db)

	query := "INSERT INTO weather (city, dt, temperature) VALUES ($1, $2, $3)"

	t.Run("Success call", func(t *testing.T) {
		mock.ExpectExec(regexp.QuoteMeta(query)).
			WithArgs(wm.City, wm.TimeStamp, wm.Temperature).
			WillReturnResult(sqlmock.NewResult(0, 1))

		err := ds.AddWeather(wm)
		assert.NoError(t, err)

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})

	t.Run("Error call", func(t *testing.T) {
		mock.ExpectExec(regexp.QuoteMeta(query)).
			WithArgs(wm.City, wm.TimeStamp, wm.Temperature).
			WillReturnError(errors.New("error"))

		err := ds.AddWeather(wm)
		assert.Error(t, err)

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})
}

func TestGetListWeatherRequest(t *testing.T) {
	db, mock := NewMock()
	defer db.Close()

	ds := NewDS(db)

	query := "SELECT id, city, dt, temperature FROM  weather"

	t.Run("Success call", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "city", "dt", "temperature"}).
			AddRow(w.ID, w.City, w.TimeStamp, w.Temperature)

		mock.ExpectQuery(query).WillReturnRows(rows)

		resp, err := ds.GetListWeatherRequest()
		assert.NotEmpty(t, resp)
		assert.NoError(t, err)
		assert.Len(t, resp, 1)

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})

	t.Run("Error call", func(t *testing.T) {
		mock.ExpectQuery(query).WillReturnError(errors.New("error"))

		resp, err := ds.GetListWeatherRequest()
		assert.Empty(t, resp)
		assert.Error(t, err)

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})
}

func TestNewDS(t *testing.T) {
	db, _ := NewMock()
	defer db.Close()
	ds := NewDS(db)
	assert.IsType(t, DS{}, ds)
}
