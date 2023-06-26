package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/go-sql-driver/mysql"
)

type GithubResponse struct {
	GithubData struct {
		Login string `json:"login"`
		Id    int    `json:"id"`
	}
}

func LoggedInHandler(w http.ResponseWriter, r *http.Request, githubData string) {
	if githubData == "" {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintf(w, `{"error": "Unauthorized"}`)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	// un-marshall the data before marshalling with githubOrgs
	var data GithubResponse
	err := json.Unmarshal([]byte(string(githubData)), &data)
	if err != nil {
		fmt.Println("Error parsing github data", err)
		return
	}

	response := struct {
		Message string      `json:"message"`
		Data    interface{} `json:"data"`
	}{
		Message: "GitHub response processed successfully",
		Data:    data,
	}

	responseJSON, err := json.Marshal(response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `{"error": "Failed to marshal response JSON"}`)
		return
	}

	fmt.Fprintf(w, string(responseJSON))

	username := data.GithubData.Login
	fmt.Printf("Username: %s\n", username)
	if username != "" {
		err := createMySQLEntry(username)
		if err != nil {
			fmt.Println("Failed to create MySQL entry:", err)
		}
	}
}

func createMySQLEntry(username string) error {
	dbUsername := os.Getenv("MYSQL_USERNAME")
	password := os.Getenv("MYSQL_PASSWORD")
	host := os.Getenv("MYSQL_HOST")
	port := os.Getenv("MYSQL_PORT")
	database := os.Getenv("MYSQL_DATABASE")

	cfg := mysql.Config{
		User:                 dbUsername,
		Passwd:               password,
		Net:                  "tcp",
		Addr:                 fmt.Sprintf("%s:%s", host, port),
		DBName:               database,
		AllowNativePasswords: true,
	}

	db, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		return err
	}
	defer db.Close()

	// Create the database if it doesn't exist
	_, err = db.Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s", database))
	if err != nil {
		return err
	}

	// Select the created database
	_, err = db.Exec(fmt.Sprintf("USE %s", database))
	if err != nil {
		return err
	}

	// Create the table if it doesn't exist
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id INT AUTO_INCREMENT PRIMARY KEY,
			username VARCHAR(255)
		)`)
	if err != nil {
		return err
	}

	stmt, err := db.Prepare("INSERT INTO users (username) VALUES (?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(username)
	if err != nil {
		return err
	}

	return nil
}
