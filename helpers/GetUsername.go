package helpers

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

type User struct {
	ID       int
	Username string
	// Add other fields as per your table structure
}

func GetUsername() (string, error) {
	//dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUsername, dbPassword, dbHost, dbPort, dbDatabase)
	db, err := sql.Open("mysql", "root:root@tcp(localhost:3306)/db01")
	if err != nil {
		return "", err
	}
	defer db.Close()

	query := "SELECT username FROM users WHERE id = (SELECT MAX(id) FROM users)"
	row := db.QueryRow(query)

	var username string
	if err := row.Scan(&username); err != nil {
		fmt.Println("Error occurred:", err)
		return "", err
	}

	return username, nil
}
