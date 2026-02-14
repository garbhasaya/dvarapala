package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	dbPath := "data/dvarapala.db"
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatalf("failed to open database: %v", err)
	}
	defer db.Close()

	email := "admin@dharma.com"
	password := "test1234"

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("failed to hash password: %v", err)
	}

	res, err := db.Exec("UPDATE user SET password = ? WHERE email = ?", string(hashedPassword), email)
	if err != nil {
		log.Fatalf("failed to update password: %v", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		log.Fatalf("failed to get rows affected: %v", err)
	}

	if rowsAffected == 0 {
		fmt.Printf("User with email %s not found\n", email)
	} else {
		fmt.Printf("Password updated successfully for user %s\n", email)
	}
}
