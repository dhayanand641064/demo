package main

import (
	"fmt"
	"log"
	"net/http"

	"demo.com/handlers"
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
	})

	fmt.Println("[ UP ON PORT 3000 ]")
	log.Panic(http.ListenAndServe(":3000", nil))
}
