package controllers

import (
	"Biblioteca/database"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/gorilla/mux"
)

func UploadPDF(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	if id == 0 {
		http.Error(w, "id inválido", http.StatusBadRequest)
		return
	}

	if err := r.ParseMultipartForm(20 << 20); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	file, hdr, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "archivo requerido", http.StatusBadRequest)
		return
	}
	defer file.Close()

	if filepath.Ext(hdr.Filename) != ".pdf" {
		http.Error(w, "solo se aceptan PDFs", http.StatusUnsupportedMediaType)
		return
	}

	if err := os.MkdirAll("uploads", 0755); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Guardamos como uploads/book_<id>.pdf
	filename := "book_" + strconv.Itoa(id) + ".pdf"
	dstPath := filepath.Join("uploads", filename)
	dst, err := os.Create(dstPath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer dst.Close()
	if _, err = io.Copy(dst, file); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Ruta "/uploads/…"
	publicURL := "/uploads/" + filename
	if _, err = database.DB.Exec(
		`UPDATE books SET pdf_path = ? WHERE id = ?`, publicURL, id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("PDF cargado"))
}
