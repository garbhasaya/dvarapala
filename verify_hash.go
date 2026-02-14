package main

import (
	"fmt"
	"log"

	"golang.org/x/crypto/bcrypt"
)

func main() {
	password := "test1234"
	hash := "$2a$10$1JG1vP/d/QP4Qm/SVTDu/uMEXtzF856/Z07eanmNOtYKb7IDC7dLO"

	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		log.Fatalf("Verification failed: %v", err)
	}
	fmt.Println("Password match!")
}
