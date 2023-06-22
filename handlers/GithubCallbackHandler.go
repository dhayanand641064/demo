package handlers

import (
	"encoding/json"
	"net/http"

	"demo.com/helpers"
)

func GithubCallbackHandler(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	githubAccessToken := helpers.GetGithubAccessToken(code)
	githubData := helpers.GetGithubData(githubAccessToken)
	githubOrgs := helpers.GetGithubOrganizations(githubAccessToken)

	response := struct {
		GithubData string   `json:"githubData"`
		GithubOrgs []string `json:"githubOrgs"`
	}{
		GithubData: githubData,
		GithubOrgs: githubOrgs,
	}

	responseJSON, _ := json.Marshal(response)

	http.Redirect(w, r, "/loggedin?githubData="+string(responseJSON), http.StatusSeeOther)
}
