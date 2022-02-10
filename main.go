package main

import (
	"BookingFlightRESTfullAPI/Storage"
	"log"
	"net/http"
)

//@title
//@version 1.0
//@description API Server for Booking flight

//@host localhost:8000
//@BasePath /

func main() {
	mux := http.NewServeMux()
	UsersServer := Storage.NewAdminStorage()
	mux.HandleFunc("/users/", UsersServer.AdminHandler)
	mux.HandleFunc("/tickets/", UsersServer.AdminHandler)
	mux.HandleFunc("/flights/", UsersServer.AdminHandler)
	log.Fatal(http.ListenAndServe(":8080", mux))
}
