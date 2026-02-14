//go:build tools

package main

import (
	"fmt"
	"os"

	"golang.org/x/crypto/bcrypt"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: go run verify_hash.go <password> <hash>")
		return
	}
	password := os.Args[1]
	hash := os.Args[2]
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		fmt.Println("Password does not match hash.")
		return
	}
	fmt.Println("Password matches hash.")
}
