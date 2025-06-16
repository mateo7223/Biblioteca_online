package main

import (
	"Biblioteca/database"
	"Biblioteca/routes"
	"log"
	"net/http"
)

func main() {
	if err := database.InitDB(); err != nil {
		log.Fatal("Error al conectar a la BD:", err)
	}

	// Servir archivos estáticos
	fs := http.FileServer(http.Dir("./public"))
	http.Handle("/public/", http.StripPrefix("/public/", fs))

	// Agregar public de backend
	http.Handle("/", routes.Router())

	log.Println("Servidor en http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
