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

	http.Handle("/", routes.Router())

	log.Println("Servidor en http://localhost:8080")

	http.Handle("/uploads/",
		http.StripPrefix("/uploads/",
			http.FileServer(http.Dir("./uploads"))))

	http.ListenAndServe(":8080", nil)
}
