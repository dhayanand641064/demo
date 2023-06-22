package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/go-sql-driver/mysql"
)

type GitHubResponse struct {
	Login string `json:"login"`
}

func LoggedInHandler(w http.ResponseWriter, r *http.Request, githubData interface{}) {
	if githubData == nil {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintf(w, `{"error": "Unauthorized"}`)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	response := struct {
		Message string      `json:"message"`
		Data    interface{} `json:"data"`
	}{
		Message: "GitHub response processed successfully",
		Data:    githubData,
	}

	responseJSON, err := json.Marshal(response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `{"error": "Failed to marshal response JSON"}`)
		return
	}

	fmt.Fprintf(w, string(responseJSON))

	githubDataJSON, err := json.Marshal(githubData)
	if err != nil {
		fmt.Println("Failed to marshal GitHub data:", err)
		return
	}

	username := ExtractLogin(githubDataJSON)
	fmt.Printf("Username: %s\n", username)
	if username != "" {
		err := createMySQLEntry(username)
		if err != nil {
			fmt.Println("Failed to create MySQL entry:", err)
		}
	}
}

func ExtractLogin(githubData []byte) string {
	var data map[string]interface{}
	err := json.Unmarshal(githubData, &data)
	if err != nil {
		fmt.Println("Failed to unmarshal GitHub response:", err)
		return ""
	}

	login, ok := data["login"].(string)
	if !ok {
		fmt.Println("Failed to extract login field from GitHub response")
		return ""
	}

	return login
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
