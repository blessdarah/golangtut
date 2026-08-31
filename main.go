package main

import (
	"blessdarah/tuts/user"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type api struct {
	store user.Store
}

func (a *api) ListUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	users := a.store.GetAll()

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(users)
}

// Adduser adds a new user to the store.
// example: POST http://localhost:8080/users
// Content-Type: application/json
//
//	{
//	    "name": "John Doe",
//	    "age": 30,
//	    "email": "john@doe.com"
//	    "zipCode": "12345",
//	    "city": "New York",
//	    "street": "Broadway"
//	}
func (a *api) AddUser(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var user user.User

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = a.store.Add(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
}

func main() {
	api := &api{
		store: user.NewUserStore(),
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/users", api.ListUsers)
	mux.HandleFunc("/users/create", api.AddUser)

	fmt.Println("Server started on port 8080")
	log.Fatal(http.ListenAndServe(":8080", mux))

}
