package Storage

import (
	"encoding/json"
	"errors"
	"log"
	"mime"
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

func (f *FlightStorage) CreateFlight(flights *Flights) Flights {
	f.Lock()
	defer f.Unlock()
	flights.Id = f.counter
	f.flightMap[flights.Id] = *flights
	f.counter++
	return f.flightMap[flights.Id]
}

func (f *FlightStorage) UpdateFlight(id int, flights *Flights) (Flights, error) {
	f.Lock()
	defer f.Unlock()
	if _, ok := f.flightMap[id]; !ok {
		return Flights{}, errors.New("no such flight")
	} else {
		f.flightMap[id] = *flights
		return f.flightMap[id], nil
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
	log.Printf("handling flight create at %s\n", req.URL.Path)
	var flights Flights
	contentType := req.Header.Get("Content-Type")
	mediatype, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if mediatype != "application/json" {
		http.Error(w, "expect application/json Content-Type", http.StatusUnsupportedMediaType)
		return
	}
	err2 := json.NewDecoder(req.Body).Decode(&flights)
	if err2 != nil {
		http.Error(w, "problem with ticket creation", http.StatusBadRequest)
		return
	}
	createFlight := f.CreateFlight(&flights)
	RenderJSON(w, createFlight)
}

func (f *FlightStorage) getFlightHandler(w http.ResponseWriter, req *http.Request) {
	log.Printf("handling flight get at %s\n", req.URL.Path)
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
	RenderJSON(w, flight)
}

func (f *FlightStorage) updateFlightHandler(w http.ResponseWriter, req *http.Request) {
	log.Printf("handling flight update at %s\n", req.URL.Path)
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
		http.Error(w, "problem with flight update", http.StatusBadRequest)
		return
	}
	updateFlight, err2 := f.UpdateFlight(id, &flight)
	if err2 != nil {
		http.Error(w, "problem with flight update", http.StatusBadRequest)
		return
	}
	RenderJSON(w, updateFlight)
}

func (f *FlightStorage) deleteFlightHandler(w http.ResponseWriter, req *http.Request) {
	log.Printf("handling flight delete at %s\n", req.URL.Path)
	path := strings.Trim(req.URL.Path, "/")
	pathParts := strings.Split(path, "/")
	if len(pathParts) < 2 {
		http.Error(w, "expect /ticket/<id> in ticket handler", http.StatusBadRequest)
		return
	}
	id, _ := strconv.Atoi(pathParts[1])
	f.DeleteFlight(id)
	RenderJSON(w, id)
}
