package controllers

import (
	"Biblioteca/database"
	"database/sql"
	"encoding/json"
	"net/http"
	"strings"
)

func Login(w http.ResponseWriter, r *http.Request) {
	// --- 1. Leer credenciales ---
	var (
		username string
		password string
		isJSON   bool
	)

	ct := r.Header.Get("Content-Type")
	if strings.HasPrefix(ct, "application/json") {
		isJSON = true
		var payload struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			http.Error(w, "JSON inválido", http.StatusBadRequest)
			return
		}
		username, password = payload.Username, payload.Password
	} else {
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Datos inválidos", http.StatusBadRequest)
			return
		}
		username = r.FormValue("user")
		password = r.FormValue("pass")
	}

	const q = `SELECT id, password, role FROM users WHERE username = ? LIMIT 1`
	var (
		id     int
		dbPass string
		role   string
	)
	err := database.DB.QueryRow(q, username).Scan(&id, &dbPass, &role)
	switch {
	case err == sql.ErrNoRows || dbPass != password:
		http.Error(w, "Credenciales inválidas", http.StatusUnauthorized)
		return
	case err != nil:
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//  Responder
	if isJSON || r.Header.Get("X-Requested-With") == "XMLHttpRequest" {
		resp := struct {
			ID       int    `json:"id"`
			Username string `json:"username"`
			Role     string `json:"role"`
		}{ID: id, Username: username, Role: role}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
		return
	}

	target := "/public/menu.html"
	if role == "Lector" {
		target = "/public/reader_menu.html"
	}
	http.Redirect(w, r, target, http.StatusSeeOther)
}
