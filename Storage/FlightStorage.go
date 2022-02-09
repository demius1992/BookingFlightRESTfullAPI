package Storage

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

/*
POST   /flights/                : создает рейс и возвращает его ID
GET    /flights/<flightsid>      :  возвращает рейс по его ID
DELETE /flights/<flightsid>      :  удаляет рейс по его ID
PUT    /flights/<flightsid>      :  изменяет рейс по его ID
*/

type Flights struct {
	Id                                 int       `json:"id"`
	StartDateAndTime                   time.Time `json:"startDateAndTime"`
	EndDateAndTime                     time.Time `json:"endDateAndTime"`
	DeparturePoint                     string    `json:"departurePoint"`
	DestinationPoint                   string    `json:"destinationPoint"`
	PriceOfTicketsOfEachRank           int       `json:"priceOfTicketsOfEachRank"`
	TotalNumberOfTicketsOfEachRank     int       `json:"totalNumberOfTicketsOfEachRank"`
	NumberOfAvailableTicketsOfEachRank int       `json:"numberOfAvailableTicketsOfEachRank"`
	DateOfCreation                     time.Time `json:"dateOfCreation"`
}

type FlightStorage struct {
	counter   int
	flightMap map[int]Flights
	sync.Mutex
}

func NewFlightStorage() *FlightStorage {
	return &FlightStorage{
		flightMap: make(map[int]Flights),
		counter:   1,
	}
}

func (f *FlightStorage) GetFlight(id int) (Flights, error) {
	f.Lock()
	defer f.Unlock()
	if flight, ok := f.flightMap[id]; !ok {
		return flight, errors.New("no such flight")
	} else {
		return flight, nil
	}
}

func (f *FlightStorage) CreateFlight(flights *Flights) {
	f.Lock()
	defer f.Unlock()
	flights.Id = f.counter
	f.flightMap[flights.Id] = *flights
	f.counter++
}

func (f *FlightStorage) UpdateFlight(id int, flights *Flights) error {
	f.Lock()
	defer f.Unlock()
	if _, ok := f.flightMap[id]; !ok {
		return errors.New("no such flight")
	} else {
		f.flightMap[id] = *flights
		return nil
	}
}

func (f *FlightStorage) DeleteFlight(id int) error {
	f.Lock()
	defer f.Unlock()
	if _, ok := f.flightMap[id]; !ok {
		return errors.New("no such flight")
	} else {
		delete(f.flightMap, id)
		return nil
	}
}

func (f *FlightStorage) createFlightHandler(w http.ResponseWriter, req *http.Request) {
	var flights Flights
	err := json.NewDecoder(req.Body).Decode(&flights)
	if err != nil {
		http.Error(w, "problem with ticket creation", http.StatusBadRequest)
		return
	}
	f.CreateFlight(&flights)
}

func (f *FlightStorage) getFlightHandler(w http.ResponseWriter, req *http.Request) {
	path := strings.Trim(req.URL.Path, "/")
	pathParts := strings.Split(path, "/")
	if len(pathParts) < 2 {
		http.Error(w, "expect /flight/<id> in flight handler", http.StatusBadRequest)
		return
	}
	id, _ := strconv.Atoi(pathParts[1])
	flight, err := f.GetFlight(id)
	if err != nil {
		http.Error(w, "wrong id or no such flight", http.StatusBadRequest)
		return
	}
	js, err1 := json.Marshal(flight)
	if err1 != nil {
		http.Error(w, "problems with json", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func (f *FlightStorage) updateFlightHandler(w http.ResponseWriter, req *http.Request) {
	var flight Flights
	path := strings.Trim(req.URL.Path, "/")
	pathParts := strings.Split(path, "/")
	if len(pathParts) < 2 {
		http.Error(w, "expect /flight/<id> in flight handler", http.StatusBadRequest)
		return
	}
	id, _ := strconv.Atoi(pathParts[1])
	err := json.NewDecoder(req.Body).Decode(&flight)
	if err != nil {
		http.Error(w, "problem with ticket creation", http.StatusBadRequest)
		return
	}
	f.UpdateFlight(id, &flight)
}

func (f *FlightStorage) deleteFlightHandler(w http.ResponseWriter, req *http.Request) {
	path := strings.Trim(req.URL.Path, "/")
	pathParts := strings.Split(path, "/")
	if len(pathParts) < 2 {
		http.Error(w, "expect /ticket/<id> in ticket handler", http.StatusBadRequest)
		return
	}
	id, _ := strconv.Atoi(pathParts[1])
	f.DeleteFlight(id)
}
