package helpers

import (
	"log"
	"os"
)

func GetGithubClientSecret() string {
	githubClientSecret, exists := os.LookupEnv("CLIENT_SECRET")
	if !exists {
		log.Fatal("Github Client Secret not defined in .env file")
	}
	return githubClientSecret
}
