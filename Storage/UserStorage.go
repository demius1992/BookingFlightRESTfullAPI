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

func (u *UserStorage) UpdateUser(id int, user Users) error {
	u.Lock()
	defer u.Unlock()
	if _, ok := u.userMap[id]; !ok {
		return errors.New("no such user")
	} else {
		u.userMap[id] = user
		return nil
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

func (u *UserStorage) getUserHandler(w http.ResponseWriter, req *http.Request) {
	path := strings.Trim(req.URL.Path, "/")
	pathParts := strings.Split(path, "/")
	if len(pathParts) < 2 {
		http.Error(w, "expect /user/<id> in user handler", http.StatusBadRequest)
		return
	}
	id, _ := strconv.Atoi(pathParts[1])
	users, _ := u.GetUser(id)
	js, err := json.Marshal(users)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	var err2 error
	_, err2 = w.Write(js)
	if err2 != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (u *UserStorage) UpdateUserHandler(w http.ResponseWriter, req *http.Request) {
	var user Users
	path := strings.Trim(req.URL.Path, "/")
	pathParts := strings.Split(path, "/")
	if len(pathParts) < 2 {
		http.Error(w, "expect /user/<id> in user handler", http.StatusBadRequest)
		return
	}
	id, _ := strconv.Atoi(pathParts[1])
	err := json.NewDecoder(req.Body).Decode(&user)
	if err != nil {
		http.Error(w, "json problem", http.StatusBadRequest)
		return
	}
	u.UpdateUser(id, user)
}

func (u *UserStorage) DeleteUserHandler(w http.ResponseWriter, req *http.Request) {
	path := strings.Trim(req.URL.Path, "/")
	pathParts := strings.Split(path, "/")
	if len(pathParts) < 2 {
		http.Error(w, "expect /user/<id> in user handler", http.StatusBadRequest)
		return
	}
	id, _ := strconv.Atoi(pathParts[1])
	err := u.DeleteUser(id)
	if err != nil {
		http.Error(w, "expect /user/<id> in user handler", http.StatusBadRequest)
		return
	}
}
