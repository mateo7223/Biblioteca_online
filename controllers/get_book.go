package controllers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"Biblioteca/database"
	"Biblioteca/models"

	"github.com/gorilla/mux"
)

func GetBook(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	const q = `SELECT id, title, author, price, format,
	                  available, stock, created_at
	           FROM books
	           WHERE id = ? LIMIT 1`

	var b models.Book
	err := database.DB.QueryRow(q, id).Scan(
		&b.ID, &b.Title, &b.Author, &b.Price, &b.Format,
		&b.Available, &b.Stock, &b.CreatedAt,
	)
	if err == sql.ErrNoRows {
		http.Error(w, "libro no encontrado", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(b)
}
