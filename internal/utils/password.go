package utils

import "golang.org/x/crypto/bcrypt"

// HashPassword gera um hash bcrypt de uma senha.
// O custo (cost) define o quão computacionalmente caro é o hash.
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14) // O custo 14 é um bom padrão
	return string(bytes), err
}

// CheckPasswordHash compara uma senha em texto puro com um hash bcrypt.
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
