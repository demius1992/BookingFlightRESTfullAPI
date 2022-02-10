package Storage

import (
	"BookingFlightRESTfullAPI/Service"
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"
)

/*
GET    /users/<userId>      :  возвращает одного пользователя по его ID
DELETE /users/<userId>      :  удаляет пользователя по его ID
PUT    /users/<userId>      :  изменяет одного пользователя по его ID
*/

type Users struct {
	Id             int       `json:"id"`
	Phone          uint64    `json:"phone"`
	Email          string    `json:"email"`
	BookedTickets  bool      `json:"bookedTickets"`
	Password       string    `json:"password"`
	DateOfCreation time.Time `json:"dateOfCreation"`
	DatOfUpdate    time.Time `json:"datOfUpdate"`
}

type UserStorage struct {
	counter int
	userMap map[int]Users
	sync.Mutex
}

func NewUserStorage() *UserStorage {
	return &UserStorage{
		userMap: make(map[int]Users),
		counter: 1,
	}
}

func (u *UserStorage) GetUser(id int) (Users, error) {
	u.Lock()
	defer u.Unlock()
	if user, ok := u.userMap[id]; !ok {
		return user, errors.New("no such user")
	} else {
		return user, nil
	}
}

func (u *UserStorage) UpdateUser(id int, user Users) (Users, error) {
	u.Lock()
	defer u.Unlock()
	if _, ok := u.userMap[id]; !ok {
		return Users{}, errors.New("no such user")
	} else {
		u.userMap[id] = user
		return u.userMap[id], nil
	}
}

func (u *UserStorage) DeleteUser(id int) error {
	u.Lock()
	defer u.Unlock()
	if _, ok := u.userMap[id]; !ok {
		return errors.New("no such user")
	} else {
		delete(u.userMap, id)
		return nil
	}
}

func (u *UserStorage) GetUserHandler(w http.ResponseWriter, req *http.Request) {
	log.Printf("handling user get at %s\n", req.URL.Path)
	id, _ := strconv.Atoi(mux.Vars(req)["id"])
	users, err := u.GetUser(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	Service.RenderJSON(w, users)
}

func (u *UserStorage) UpdateUserHandler(w http.ResponseWriter, req *http.Request) {
	log.Printf("handling user update at %s\n", req.URL.Path)
	var user Users
	id, _ := strconv.Atoi(mux.Vars(req)["id"])
	err := json.NewDecoder(req.Body).Decode(&user)
	if err != nil {
		http.Error(w, "problem with user update", http.StatusBadRequest)
		return
	}
	updateUser, err2 := u.UpdateUser(id, user)
	if err2 != nil {
		http.Error(w, err2.Error(), http.StatusNotFound)
		return
	}
	Service.RenderJSON(w, updateUser)
}

func (u *UserStorage) DeleteUserHandler(w http.ResponseWriter, req *http.Request) {
	log.Printf("handling user delete at %s\n", req.URL.Path)
	id, _ := strconv.Atoi(mux.Vars(req)["id"])
	err := u.DeleteUser(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	Service.RenderJSON(w, id)
}
