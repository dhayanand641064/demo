package handlers

import (
	"encoding/json"
	"net/http"
	"fmt"

	"demo.com/helpers"
)

type UserCred struct {
	Login string `json:"login"`
	Id int `json:"id"`
}

func GithubCallbackHandler(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	githubAccessToken := helpers.GetGithubAccessToken(code)
	githubData := helpers.GetGithubData(githubAccessToken)
	githubOrgs := helpers.GetGithubOrganizations(githubAccessToken)

	// un-marshall the data before marshalling with githubOrgs
	var userCred UserCred
	err := json.Unmarshal([]byte(string(githubData)), &userCred)
	if err != nil {
		fmt.Println("Error parsing github data", err)
		return
	}

	response := struct {
		GithubData interface{}   `json:"githubData"`
		GithubOrgs []string 		 `json:"githubOrgs"`
	}{
		GithubData: userCred,
		GithubOrgs: githubOrgs,
	}

	responseJSON, _ := json.Marshal(response)
	http.Redirect(w, r, "/loggedin?githubData="+string(responseJSON), http.StatusSeeOther)
}
