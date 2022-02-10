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

func (t *TicketStorage) createTicketHandler(w http.ResponseWriter, req *http.Request) {
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
	RenderJSON(w, createTicket)
}

func (t *TicketStorage) getTicketHandler(w http.ResponseWriter, req *http.Request) {
	log.Printf("handling ticket get at %s\n", req.URL.Path)
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
	RenderJSON(w, ticket)
}

func (t *TicketStorage) updateTicketHandler(w http.ResponseWriter, req *http.Request) {
	log.Printf("handling ticket update at %s\n", req.URL.Path)
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
		http.Error(w, "problem with ticket update", http.StatusBadRequest)
		return
	}
	updateTicket, err2 := t.UpdateTicket(id, &ticket)
	if err2 != nil {
		http.Error(w, "problem with ticket update", http.StatusBadRequest)
		return
	}
	RenderJSON(w, updateTicket)
}

func (t *TicketStorage) deleteTicketHandler(w http.ResponseWriter, req *http.Request) {
	log.Printf("handling ticket delete at %s\n", req.URL.Path)
	path := strings.Trim(req.URL.Path, "/")
	pathParts := strings.Split(path, "/")
	if len(pathParts) < 2 {
		http.Error(w, "expect /ticket/<id> in ticket handler", http.StatusBadRequest)
		return
	}
	id, _ := strconv.Atoi(pathParts[1])
	t.DeleteTicket(id)
	RenderJSON(w, id)
}
