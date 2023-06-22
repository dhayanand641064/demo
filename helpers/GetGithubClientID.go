package helpers

import (
	"log"
	"os"
)

func GetGithubClientID() string {
	githubClientID, exists := os.LookupEnv("CLIENT_ID")
	if !exists {
		log.Fatal("Github Client ID not defined in .env file")
	}
	return githubClientID
}
