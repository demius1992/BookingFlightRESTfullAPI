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
POST   /tickets/                : создает билет и возвращает его ID
GET    /tickets/<ticketid>      :  возвращает билет по его ID
DELETE /tickets/<ticketsid>      :  удаляет билет по его ID
PUT    /tickets/<ticketsid>      :  изменяет билет по его ID
*/

type Tickets struct {
	Id             int       `json:"id"`
	Flight         Flights   `json:"flight"`
	User           Users     `json:"user"`
	Rank           string    `json:"rank"` //(Economy, Business, Deluxe)
	Price          int       `json:"price"`
	DateOfCreation time.Time `json:"dateOfCreation"`
}

type TicketStorage struct {
	counter   int
	TicketMap map[int]Tickets
	sync.Mutex
}

func NewTicketStorage() *TicketStorage {
	return &TicketStorage{
		TicketMap: make(map[int]Tickets),
		counter:   1,
	}
}

func (t *TicketStorage) CreateTicket(tickets *Tickets) Tickets {
	t.Lock()
	defer t.Unlock()
	tickets.Id = t.counter
	t.TicketMap[tickets.Id] = *tickets
	t.counter++
	return t.TicketMap[tickets.Id]
}

func (t *TicketStorage) GetTicket(id int) (Tickets, error) {
	t.Lock()
	defer t.Unlock()
	if ticket, ok := t.TicketMap[id]; !ok {
		return ticket, errors.New("no such ticket")
	} else {
		return ticket, nil
	}
	//return Tickets{
	//	Id:             12,
	//	Flight:         Flights{},
	//	User:           Users{},
	//	Rank:           "economy",
	//	Price:          1,
	//	DateOfCreation: time.Time{},
	//}, nil
}

func (t *TicketStorage) UpdateTicket(id int, tickets *Tickets) (Tickets, error) {
	t.Lock()
	defer t.Unlock()
	if _, ok := t.TicketMap[id]; !ok {
		return Tickets{}, errors.New("no such ticket")
	} else {
		t.TicketMap[id] = *tickets
		return t.TicketMap[id], nil
	}
}

func (t *TicketStorage) DeleteTicket(id int) error {
	t.Lock()
	defer t.Unlock()
	if _, ok := t.TicketMap[id]; !ok {
		return errors.New("no such ticket")
	} else {
		delete(t.TicketMap, id)
		return nil
	}
}

func (t *TicketStorage) CreateTicketHandler(w http.ResponseWriter, req *http.Request) {
	log.Printf("handling ticket create at %s\n", req.URL.Path)
	var ticket Tickets
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
	err2 := json.NewDecoder(req.Body).Decode(&ticket)
	if err2 != nil {
		http.Error(w, "problem with ticket creation", http.StatusBadRequest)
		return
	}
	createTicket := t.CreateTicket(&ticket)
	Service.RenderJSON(w, createTicket)
}

func (t *TicketStorage) GetTicketHandler(w http.ResponseWriter, req *http.Request) {
	log.Printf("handling ticket get at %s\n", req.URL.Path)
	id, _ := strconv.Atoi(mux.Vars(req)["id"])
	ticket, err := t.GetTicket(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	Service.RenderJSON(w, ticket)
}

func (t *TicketStorage) UpdateTicketHandler(w http.ResponseWriter, req *http.Request) {
	log.Printf("handling ticket update at %s\n", req.URL.Path)
	var ticket Tickets
	id, _ := strconv.Atoi(mux.Vars(req)["id"])
	err := json.NewDecoder(req.Body).Decode(&ticket)
	if err != nil {
		http.Error(w, "problem with ticket update", http.StatusBadRequest)
		return
	}
	updateTicket, err2 := t.UpdateTicket(id, &ticket)
	if err2 != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	Service.RenderJSON(w, updateTicket)
}

func (t *TicketStorage) DeleteTicketHandler(w http.ResponseWriter, req *http.Request) {
	log.Printf("handling ticket delete at %s\n", req.URL.Path)
	id, _ := strconv.Atoi(mux.Vars(req)["id"])
	err := t.DeleteTicket(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
	}
	Service.RenderJSON(w, id)
}
