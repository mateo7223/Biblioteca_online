# Biblioteca_online

Plataforma en Go para vender libros PDF.
REQUISITOS
- Go 1.22+
- MySQL 8
INSTALACIÓN
git clone <repo>
cd Biblioteca_online
go mod download
go run main.go
ESTRUCTURA
Biblioteca_online/
■■■ controllers/ # Lógica HTTP
■■■ models/ # Modelos
■■■ routes/ # Rutas
■■■ database/ # Conexión BD
■■■ public/ # HTML/CSS
API BÁSICA
POST /users -> crear usuario
POST /login -> login
GET /books -> listar libros
POST /books -> crear libro
POST /purchases -> comprar libro