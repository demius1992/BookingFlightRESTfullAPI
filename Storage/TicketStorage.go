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

func (t *TicketStorage) GetTicket(id int) (Tickets, error) {
	t.Lock()
	defer t.Unlock()
	if ticket, ok := t.TicketMap[id]; !ok {
		return ticket, errors.New("no such ticket")
	} else {
		return ticket, nil
	}
}

func (t *TicketStorage) CreateTicket(tickets *Tickets) {
	t.Lock()
	defer t.Unlock()
	tickets.Id = t.counter
	t.TicketMap[tickets.Id] = *tickets
	t.counter++
}

func (t *TicketStorage) UpdateTicket(id int, tickets *Tickets) error {
	t.Lock()
	defer t.Unlock()
	if _, ok := t.TicketMap[id]; !ok {
		return errors.New("no such ticket")
	} else {
		t.TicketMap[id] = *tickets
		return nil
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

func (t *TicketStorage) createTicketHandler(w http.ResponseWriter, req *http.Request) {
	var ticket Tickets
	err := json.NewDecoder(req.Body).Decode(&ticket)
	if err != nil {
		http.Error(w, "problem with ticket creation", http.StatusBadRequest)
		return
	}
	t.CreateTicket(&ticket)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
}

func (t *TicketStorage) getTicketHandler(w http.ResponseWriter, req *http.Request) {
	path := strings.Trim(req.URL.Path, "/")
	pathParts := strings.Split(path, "/")
	if len(pathParts) < 2 {
		http.Error(w, "expect /ticket/<id> in ticket handler", http.StatusBadRequest)
		return
	}
	id, _ := strconv.Atoi(pathParts[1])
	ticket, err := t.GetTicket(id)
	if err != nil {
		http.Error(w, "wrong id or no such ticket", http.StatusBadRequest)
		return
	}
	js, err1 := json.Marshal(ticket)
	if err1 != nil {
		http.Error(w, "problems with json", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func (t *TicketStorage) updateTicketHandler(w http.ResponseWriter, req *http.Request) {
	var ticket Tickets
	path := strings.Trim(req.URL.Path, "/")
	pathParts := strings.Split(path, "/")
	if len(pathParts) < 2 {
		http.Error(w, "expect /user/<id> in user handler", http.StatusBadRequest)
		return
	}
	id, _ := strconv.Atoi(pathParts[1])
	err := json.NewDecoder(req.Body).Decode(&ticket)
	if err != nil {
		http.Error(w, "problem with ticket creation", http.StatusBadRequest)
		return
	}
	t.UpdateTicket(id, &ticket)
}

func (t *TicketStorage) deleteTicketHandler(w http.ResponseWriter, req *http.Request) {
	path := strings.Trim(req.URL.Path, "/")
	pathParts := strings.Split(path, "/")
	if len(pathParts) < 2 {
		http.Error(w, "expect /ticket/<id> in ticket handler", http.StatusBadRequest)
		return
	}
	id, _ := strconv.Atoi(pathParts[1])
	t.DeleteTicket(id)
}
