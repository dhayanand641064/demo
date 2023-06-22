package helpers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func GetGithubOrganizations(accessToken string) []string {
	req, reqerr := http.NewRequest("GET", "https://api.github.com/user/orgs", nil)
	if reqerr != nil {
		log.Panic("API Request creation failed")
	}

	authorizationHeaderValue := fmt.Sprintf("token %s", accessToken)
	req.Header.Set("Authorization", authorizationHeaderValue)

	resp, resperr := http.DefaultClient.Do(req)
	if resperr != nil {
		log.Panic("Request failed")
	}

	defer resp.Body.Close()
	respbody, _ := ioutil.ReadAll(resp.Body)

	type githubOrg struct {
		Login string `json:"login"`
	}

	var orgs []githubOrg
	json.Unmarshal(respbody, &orgs)

	orgNames := make([]string, len(orgs))
	for i, org := range orgs {
		orgNames[i] = org.Login
	}

	return orgNames
}
