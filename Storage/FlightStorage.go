package Storage

import (
	"BookingFlightRESTfullAPI/Service"
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"log"
	"mime"
	"net/http"
	"strconv"
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

func (f *FlightStorage) CreateFlightHandler(w http.ResponseWriter, req *http.Request) {
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
	Service.RenderJSON(w, createFlight)
}

func (f *FlightStorage) GetFlightHandler(w http.ResponseWriter, req *http.Request) {
	log.Printf("handling flight get at %s\n", req.URL.Path)
	id, _ := strconv.Atoi(mux.Vars(req)["id"])
	flight, err := f.GetFlight(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	Service.RenderJSON(w, flight)
}

func (f *FlightStorage) UpdateFlightHandler(w http.ResponseWriter, req *http.Request) {
	log.Printf("handling flight update at %s\n", req.URL.Path)
	var flight Flights
	id, _ := strconv.Atoi(mux.Vars(req)["id"])
	err := json.NewDecoder(req.Body).Decode(&flight)
	if err != nil {
		http.Error(w, "problem with flight update", http.StatusBadRequest)
		return
	}
	updateFlight, err2 := f.UpdateFlight(id, &flight)
	if err2 != nil {
		http.Error(w, err2.Error(), http.StatusNotFound)
		return
	}
	Service.RenderJSON(w, updateFlight)
}

func (f *FlightStorage) DeleteFlightHandler(w http.ResponseWriter, req *http.Request) {
	log.Printf("handling flight delete at %s\n", req.URL.Path)
	id, _ := strconv.Atoi(mux.Vars(req)["id"])
	err := f.DeleteFlight(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
	}
	Service.RenderJSON(w, id)
}
