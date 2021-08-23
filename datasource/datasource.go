package datasource

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
	"github.com/natalliakoita/weather_backend/models"
)

// create new connection in database
func NewDB() (*sql.DB, error) {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	conn, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, err
	}

	err = conn.Ping()
	if err != nil {
		return nil, err
	}

	return conn, nil
}

type DS struct {
	DB *sql.DB
}

func NewDS(conn *sql.DB) DS {
	return DS{
		DB: conn,
	}
}

func (ds *DS) AddWeather(w models.Weather) error {
	insertStatement := "INSERT INTO weather (city, dt, temperature) VALUES ($1, $2, $3)"

	_, err := ds.DB.Exec(insertStatement, w.City, w.TimeStamp, w.Temperature)

	if err != nil {
		return err
	}
	return nil
}

func (ds *DS) GetListWeatherRequest() ([]models.Weather, error) {
	getList := `SELECT id, city, dt, temperature FROM  weather`
	rows, err := ds.DB.Query(getList)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	ws := []models.Weather{}
	for rows.Next() {
		w := models.Weather{}
		err := rows.Scan(&w.ID, &w.City, &w.TimeStamp, &w.Temperature)
		if err != nil {
			fmt.Println(err)
			continue
		}
		ws = append(ws, w)

	}
	return ws, nil
}

type DSWeatherInterface interface {
	AddWeather(w models.Weather) error
	GetListWeatherRequest() ([]models.Weather, error)
}
