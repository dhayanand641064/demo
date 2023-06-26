package main

import (
	"fmt"
	"log"
	"net/http"

	"demo.com/handlers"
	"demo.com/helpers"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("No .env file found")
	}
}

func main() {
	http.HandleFunc("/", handlers.RootHandler)
	http.HandleFunc("/login/github/", handlers.GithubLoginHandler)
	http.HandleFunc("/login/github/callback", handlers.GithubCallbackHandler)
	http.HandleFunc("/loggedin", func(w http.ResponseWriter, r *http.Request) {
		githubData := r.URL.Query().Get("githubData")
		handlers.LoggedInHandler(w, r, githubData)
		lastUsername, err := helpers.GetUsername()
		if err != nil {
			fmt.Println("Error retrieving last username:", err)
			return
		}
		payload := helpers.Input{
			Username: lastUsername,
		}

		token, err := helpers.GenerateToken(payload)
		if err != nil {
			fmt.Println("Error generating token:", err)
			return
		}

		fmt.Println("Generated token:", token)
	})

	fmt.Println("[ UP ON PORT 3000 ]")
	log.Panic(http.ListenAndServe(":3000", nil))
}
