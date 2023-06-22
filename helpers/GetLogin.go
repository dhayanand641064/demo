package helpers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func GetLogin(accessToken string) string {
	req, err := http.NewRequest("GET", "https://api.github.com/user/login", nil)
	if err != nil {
		log.Panic("API request creation failed")
	}

	authorizationHeaderValue := fmt.Sprintf("token %s", accessToken)
	req.Header.Set("Authorization", authorizationHeaderValue)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Panic("Request failed")
	}

	defer resp.Body.Close()
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Panic("Failed to read response body")
	}

	var githubResponse struct {
		Login string `json:"login"`
	}

	err = json.Unmarshal(respBody, &githubResponse)
	if err != nil {
		log.Panic("Failed to extract login data from GitHub response")
	}

	return githubResponse.Login
}
