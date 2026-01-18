package main

import (
	"log"
	"net/http"
	"os"

	"hosting-backend/internal/database"
	"hosting-backend/internal/handlers/auth"

	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

	db, err := database.ConnectMySQL()
	if err != nil {
		log.Fatalf("Erro ao conectar no MySQL: %v", err)
	}
	defer db.Close()

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	// Rota de Autenticação
	http.HandleFunc("/auth/login", auth.LoginHandler(db))

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}

	log.Println("API rodando na porta", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
