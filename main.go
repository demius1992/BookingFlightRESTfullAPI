package main

import (
	"BookingFlightRESTfullAPI/Storage"
	"github.com/gorilla/mux"
	"log"
	"math/rand"
	"net/http"
	"time"
)

type AdminStorage struct {
	*Storage.UserStorage
	*Storage.TicketStorage
	*Storage.FlightStorage
}

func NewAdminStorage() *AdminStorage {
	return &AdminStorage{
		UserStorage:   Storage.NewUserStorage(),
		TicketStorage: Storage.NewTicketStorage(),
		FlightStorage: Storage.NewFlightStorage(),
	}
}

func main() {
	router := mux.NewRouter()
	router.StrictSlash(true)
	server := NewAdminStorage()

	router.HandleFunc("/users/{id:[0-9]+}/", server.DeleteUserHandler).Methods("DELETE")
	router.HandleFunc("/users/{id:[0-9]+}/", server.UpdateUserHandler).Methods("PUT")
	router.HandleFunc("/users/{id:[0-9]+}/", server.GetUserHandler).Methods("GET")

	router.HandleFunc("/tickets/", server.CreateTicketHandler).Methods("POST")
	router.HandleFunc("/tickets/{id:[0-9]+}/", server.GetTicketHandler).Methods("GET")
	router.HandleFunc("/tickets/{id:[0-9]+}/", server.UpdateTicketHandler).Methods("PUT")
	router.HandleFunc("/tickets/{id:[0-9]+}/", server.DeleteTicketHandler).Methods("DELETE")

	router.HandleFunc("/flights/", server.CreateFlightHandler).Methods("POST")
	router.HandleFunc("/flights/{id:[0-9]+}/", server.GetFlightHandler).Methods("GET")
	router.HandleFunc("/flights/{id:[0-9]+}/", server.UpdateFlightHandler).Methods("PUT")
	router.HandleFunc("/flights/{id:[0-9]+}/", server.DeleteFlightHandler).Methods("DELETE")

	router.Use(Logging)

	log.Fatal(http.ListenAndServe(":8080", router))
}

func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, req)
		s1 := rand.NewSource(time.Now().UnixNano())
		r1 := rand.New(s1)
		id := r1.Intn(100)
		log.Printf("unique id : %d, method: %s, incoming request: %s, request timing: %s",
			id, req.Method, req.RequestURI, time.Since(start))
	})
}
