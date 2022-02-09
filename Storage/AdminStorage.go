package Storage

import (
	"fmt"
	"net/http"
)

type AdminStorage struct {
	*UserStorage
	*TicketStorage
	*FlightStorage
}

func NewAdminStorage() *AdminStorage {
	return &AdminStorage{
		UserStorage:   NewUserStorage(),
		TicketStorage: NewTicketStorage(),
		FlightStorage: NewFlightStorage(),
	}
}

func (u *AdminStorage) AdminHandler(w http.ResponseWriter, req *http.Request) {
	if req.URL.Path == "/users/" {
		// Запрос направлен к "/users/", без идущего в конце ID.
		if req.Method == http.MethodGet {
			u.getUserHandler(w, req)
		} else if req.Method == http.MethodPut {
			u.UpdateUserHandler(w, req)
		} else if req.Method == http.MethodDelete {
			u.DeleteUserHandler(w, req)
		} else {
			http.Error(w, fmt.Sprintf("expect method GET, DELETE or POST at /task/, got %v", req.Method), http.StatusMethodNotAllowed)
			return
		}
	} else if req.URL.Path == "/tickets/" {
		// Запрос направлен к "/tickets/", без идущего в конце ID.
		if req.Method == http.MethodGet {
			u.getTicketHandler(w, req)
		} else if req.Method == http.MethodPut {
			u.updateTicketHandler(w, req)
		} else if req.Method == http.MethodDelete {
			u.deleteTicketHandler(w, req)
		} else if req.Method == http.MethodPost {
			u.createTicketHandler(w, req)
		} else {
			http.Error(w, fmt.Sprintf("expect method GET, DELETE or POST at /task/, got %v", req.Method), http.StatusMethodNotAllowed)
			return
		}
	} else if req.URL.Path == "/flights/" {
		// Запрос направлен к "/flights/", без идущего в конце ID.
		if req.Method == http.MethodGet {
			u.getFlightHandler(w, req)
		} else if req.Method == http.MethodPut {
			u.updateFlightHandler(w, req)
		} else if req.Method == http.MethodDelete {
			u.deleteFlightHandler(w, req)
		} else if req.Method == http.MethodPost {
			u.createFlightHandler(w, req)
		} else {
			http.Error(w, fmt.Sprintf("expect method GET, DELETE or POST at /task/, got %v", req.Method), http.StatusMethodNotAllowed)
			return
		}
	}
}
