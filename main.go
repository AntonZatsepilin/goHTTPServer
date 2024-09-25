package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

type User struct {
	ID   int    `json:"id"`
	Name string `json:"name,omitempty"`
}

var users = []User{{1, "Anton"}, {2, "Max"}, {3, "Tim"}}

func main() {
	http.HandleFunc("/users", loggerMiddleware(handleUsers))

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}

func loggerMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("[%s] %s\n", r.Method, r.URL)
		next(w, r)
	}
}

func handleUsers(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getUsers(w, r)
	case http.MethodPost:
		addUser(w, r)
	default:
		w.WriteHeader(http.StatusNotImplemented)
	}
}

func getUsers(w http.ResponseWriter, r *http.Request) {
	resp, err := json.Marshal(users)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(resp)
}

func addUser(w http.ResponseWriter, r *http.Request) {
	reqBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	var user User
	err = json.Unmarshal(reqBytes, &user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	users = append(users, user)
}
