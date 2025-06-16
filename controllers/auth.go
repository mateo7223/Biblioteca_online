package controllers

import (
	"Biblioteca/database"
	"database/sql"
	"encoding/json"
	"net/http"
)

func Login(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, "Error en los datos enviados", http.StatusBadRequest)
		return
	}

	var storedPassword string
	err = database.DB.QueryRow("SELECT password FROM users WHERE username = ?", data.Username).Scan(&storedPassword)
	if err == sql.ErrNoRows || storedPassword != data.Password {
		http.Error(w, "Credenciales inválidas", http.StatusUnauthorized)
		return
	}

	w.Write([]byte("inicio exitoso"))
}
