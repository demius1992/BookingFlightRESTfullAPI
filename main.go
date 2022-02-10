package main

import (
	"BookingFlightRESTfullAPI/Storage"
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	UsersServer := Storage.NewAdminStorage()
	mux.HandleFunc("/users/", UsersServer.AdminHandler)
	mux.HandleFunc("/tickets/", UsersServer.AdminHandler)
	mux.HandleFunc("/flights/", UsersServer.AdminHandler)
	log.Fatal(http.ListenAndServe(":8080", mux))
}
