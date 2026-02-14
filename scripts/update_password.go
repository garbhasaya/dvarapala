//go:build tools

package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"golang.org/x/crypto/bcrypt"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Usage: go run update_password.go <email> <new_password>")
		return
	}
	email := os.Args[1]
	newPassword := os.Args[2]

	db, err := sql.Open("sqlite3", "./data/dvarapala.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	hash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal(err)
	}

	stmt, err := db.Prepare("UPDATE users SET password = ? WHERE email = ?")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	res, err := stmt.Exec(string(hash), email)
	if err != nil {
		log.Fatal(err)
	}

	affected, err := res.RowsAffected()
	if err != nil {
		log.Fatal(err)
	}

	if affected == 0 {
		fmt.Println("No user found with that email.")
		return
	}

	fmt.Println("Password updated successfully.")
}
