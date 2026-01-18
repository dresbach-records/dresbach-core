package main

import (
	"fmt"
	"log"
	"os"

	"golang.org/x/crypto/bcrypt"
)

// Este programa gera um hash bcrypt para uma senha fornecida como argumento.
func main() {
	if len(os.Args) != 2 {
		log.Fatalf("Uso: go run cmd/tools/hash_password.go <sua-senha>")
	}

	password := os.Args[1]
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("Erro ao gerar o hash da senha: %v", err)
	}

	fmt.Println(string(hashedPassword))
}
