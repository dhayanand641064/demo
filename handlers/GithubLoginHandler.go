package handlers

import (
	"fmt"
	"net/http"

	"demo.com/helpers"
)

func GithubLoginHandler(w http.ResponseWriter, r *http.Request) {
	githubClientID := helpers.GetGithubClientID()
	redirectURL := fmt.Sprintf("https://github.com/login/oauth/authorize?client_id=%s&redirect_uri=%s&scope=user,read:org", githubClientID, "http://localhost:3000/login/github/callback")
	http.Redirect(w, r, redirectURL, 301)
}
