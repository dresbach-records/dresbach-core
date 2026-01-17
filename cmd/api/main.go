package main

import (
	"log"
	"net/http"
	"os"

	_ "hosting-backend/internal/database"

	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

	// db, err := database.ConnectMySQL()
	// if err != nil {
	// 	log.Fatalf("Erro ao conectar no MySQL: %v", err)
	// }
	// defer db.Close()

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}

	log.Println("API rodando na porta", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
