package controllers

import (
	"Biblioteca/database"
	"database/sql"
	"encoding/json"
	"net/http"
	"time"
)

func PurchaseBook(w http.ResponseWriter, r *http.Request) {
	var req struct {
		UserID   int `json:"user_id"`
		BookID   int `json:"book_id"`
		Quantity int `json:"quantity"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "JSON inválido", http.StatusBadRequest)
		return
	}
	if req.Quantity <= 0 {
		req.Quantity = 1
	}

	tx, err := database.DB.Begin()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	var price float64
	var stock int
	if err := tx.
		QueryRow("SELECT price, stock FROM books WHERE id = ? FOR UPDATE", req.BookID).
		Scan(&price, &stock); err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Libro no encontrado", http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
	if stock < req.Quantity {
		http.Error(w, "Stock insuficiente", http.StatusConflict)
		return
	}

	if _, err = tx.Exec(`UPDATE books SET stock = stock - ? WHERE id = ?`,
		req.Quantity, req.BookID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	total := price * float64(req.Quantity)
	if _, err = tx.Exec(`
		INSERT INTO purchases(user_id, book_id, quantity, total, purchased_at)
		VALUES(?,?,?,?,?)`,
		req.UserID, req.BookID, req.Quantity, total, time.Now()); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err = tx.Commit(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Compra realizada"))
}

func MyPurchases(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id")

	rows, err := database.DB.Query(`
        SELECT p.id,
               p.book_id,
               b.title,
               b.author,
               b.pdf_path,      -- ← aquí
               p.quantity,
               p.total,
               p.purchased_at
        FROM purchases p
        JOIN books b ON b.id = p.book_id
        WHERE p.user_id = ?
        ORDER BY p.purchased_at DESC`, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	type view struct {
		ID          int       `json:"id"`
		BookID      int       `json:"book_id"`
		Title       string    `json:"title"`
		Author      string    `json:"author"`
		FileURL     string    `json:"file_url"`
		Quantity    int       `json:"quantity"`
		Total       float64   `json:"total"`
		PurchasedAt time.Time `json:"purchased_at"`
	}

	var list []view
	for rows.Next() {
		var v view
		if err := rows.Scan(&v.ID, &v.BookID, &v.Title, &v.Author,
			&v.FileURL, &v.Quantity, &v.Total, &v.PurchasedAt); err == nil {
			list = append(list, v)
		}
	}
	json.NewEncoder(w).Encode(list)
}
