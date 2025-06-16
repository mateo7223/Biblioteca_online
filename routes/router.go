package routes

import (
	"Biblioteca/controllers"
	"net/http"

	"github.com/gorilla/mux"
)

func Router() http.Handler {
	r := mux.NewRouter()
	r.HandleFunc("/users", controllers.CreateUser).Methods("POST")
	r.HandleFunc("/login", controllers.Login).Methods("POST")
	r.HandleFunc("/books", controllers.CreateBook).Methods("POST")
	r.HandleFunc("/books", controllers.ListBooks).Methods("GET")
	r.HandleFunc("/books/decrement", controllers.DecrementStock).Methods("POST")
	r.HandleFunc("/books/increment", controllers.IncrementStock).Methods("POST")
	r.HandleFunc("/users", controllers.ListUsers).Methods("GET")
	return r
}
