package Storage

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
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
	fullPath := strings.Trim(req.URL.Path, "/")
	pathParts := strings.Split(fullPath, "/")
	path := pathParts[0]

	switch path {
	case "users":
		{
			switch req.Method {
			case http.MethodGet:
				u.getUserHandler(w, req)
			case http.MethodPut:
				u.UpdateUserHandler(w, req)
			case http.MethodDelete:
				u.DeleteUserHandler(w, req)
			default:
				http.Error(w, fmt.Sprintf("expect method GET, DELETE or PUT at /task/, got %v",
					req.Method), http.StatusMethodNotAllowed)
				return
			}
		}
	case "tickets":
		{
			switch req.Method {
			case http.MethodPost:
				u.createTicketHandler(w, req)
			case http.MethodGet:
				u.getTicketHandler(w, req)
			case http.MethodPut:
				u.updateTicketHandler(w, req)
			case http.MethodDelete:
				u.deleteTicketHandler(w, req)
			default:
				http.Error(w, fmt.Sprintf("expect method GET, DELETE, PUT or POST at /task/, got %v",
					req.Method), http.StatusMethodNotAllowed)
				return
			}
		}
	case "flights":
		{
			switch req.Method {
			case http.MethodPost:
				u.createFlightHandler(w, req)
			case http.MethodGet:
				u.getFlightHandler(w, req)
			case http.MethodPut:
				u.updateFlightHandler(w, req)
			case http.MethodDelete:
				u.deleteFlightHandler(w, req)
			default:
				http.Error(w, fmt.Sprintf("expect method GET, DELETE, PUT or POST at /task/, got %v",
					req.Method), http.StatusMethodNotAllowed)
				return
			}
		}
	}
}

func RenderJSON(w http.ResponseWriter, v interface{}) {
	js, err1 := json.Marshal(v)
	if err1 != nil {
		http.Error(w, "problems with json", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}
