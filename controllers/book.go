package controllers

import (
	"Biblioteca/database"
	"Biblioteca/models"
	"encoding/json"
	"net/http"
	"time"
)

func CreateBook(w http.ResponseWriter, r *http.Request) {
	var book struct {
		Title  string
		Author string
		Price  float64
		Format string
		Stock  int
	}

	err := json.NewDecoder(r.Body).Decode(&book)
	if err != nil || book.Format != "PDF" {
		http.Error(w, "Datos mal ingresados o formato no permitido", http.StatusBadRequest)
		return
	}

	_, err = database.DB.Exec(`INSERT INTO books (title, author, price, format, available, stock, created_at)
        VALUES (?, ?, ?, ?, ?, ?, ?)`,
		book.Title, book.Author, book.Price, book.Format, true, book.Stock, time.Now())

	if err != nil {
		http.Error(w, "Error al registrar libro", http.StatusInternalServerError)
		return
	}

	w.Write([]byte("Libro registrado correctamente"))
}

func ListBooks(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	title := q.Get("title")
	author := q.Get("author")
	available := q.Get("available")

	query := "SELECT id, title, author, price, format, available, stock, created_at FROM books WHERE 1=1"
	var args []interface{}

	if title != "" {
		query += " AND title LIKE ?"
		args = append(args, "%"+title+"%")
	}
	if author != "" {
		query += " AND author LIKE ?"
		args = append(args, "%"+author+"%")
	}
	if available == "true" {
		query += " AND available = true"
	}

	rows, err := database.DB.Query(query, args...)
	if err != nil {
		http.Error(w, "Error al consultar libros", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var books []models.Book
	for rows.Next() {
		var b models.Book
		if err := rows.Scan(&b.ID, &b.Title, &b.Author, &b.Price, &b.Format, &b.Available, &b.Stock, &b.CreatedAt); err == nil {
			books = append(books, b)
		}
	}

	json.NewEncoder(w).Encode(books)
}

func DecrementStock(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	_, err := database.DB.Exec("UPDATE books SET stock = stock - 1 WHERE id = ? AND stock > 0", id)
	if err != nil {
		http.Error(w, "Error al disminuir stock", http.StatusInternalServerError)
		return
	}

	w.Write([]byte("Stock decrementado"))
}

func IncrementStock(w http.ResponseWriter, r *http.Request) {
	var data struct {
		ID     int `json:"id"`
		Amount int `json:"amount"`
	}

	json.NewDecoder(r.Body).Decode(&data)

	_, err := database.DB.Exec("UPDATE books SET stock = stock + ? WHERE id = ?", data.Amount, data.ID)
	if err != nil {
		http.Error(w, "Error al incrementar stock", http.StatusInternalServerError)
		return
	}

	w.Write([]byte("Stock incrementado"))
}
