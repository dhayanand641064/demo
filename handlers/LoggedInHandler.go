package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/couchbase/gocb/v2"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
)

type GithubResponse struct {
	GithubData struct {
		Login string `json:"login"`
		Id    int    `json:"id"`
	}
}

func LoggedInHandler(w http.ResponseWriter, r *http.Request, githubData string) (string, error) {
	if githubData == "" {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintf(w, `{"error": "Unauthorized"}`)
		return "", nil
	}

	w.Header().Set("Content-Type", "application/json")

	var data GithubResponse
	err := json.Unmarshal([]byte(string(githubData)), &data)
	if err != nil {
		fmt.Println("Error parsing github data", err)
		return "", err
	}

	response := struct {
		Message string      `json:"message"`
		Data    interface{} `json:"data"`
	}{
		Message: "GitHub response processed successfully",
		Data:    data,
	}

	responseJSON, err := json.Marshal(response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `{"error": "Failed to marshal response JSON"}`)
		return "", err
	}

	fmt.Fprintf(w, string(responseJSON))

	username := data.GithubData.Login
	fmt.Printf("Username: %s\n", username)
	if username != "" {
		insertedID, err := createCouchbaseEntry(username)
		if err != nil {
			fmt.Println("Failed to create Couchbase entry:", err)
			return "", err

		} else {
			fmt.Printf("Inserted Document ID: %s\n", insertedID)
			return insertedID, nil
		}
	}

	return "", nil
}

type User struct {
	ID       string `json:"id"`
	Username string `json:"username"`
}

func createCouchbaseEntry(username string) (string, error) {
	err := godotenv.Load()
	if err != nil {
		return "", fmt.Errorf("error loading .env file: %w", err)
	}

	connectionString := os.Getenv("COUCHBASE_CONNECTION_STRING")
	bucketName := os.Getenv("COUCHBASE_BUCKET_NAME")
	dbUsername := os.Getenv("COUCHBASE_USERNAME")
	password := os.Getenv("COUCHBASE_PASSWORD")

	cluster, err := gocb.Connect(connectionString, gocb.ClusterOptions{
		Username: dbUsername,
		Password: password,
	})
	if err != nil {
		return "", fmt.Errorf("error connecting to Couchbase: %w", err)
	}

	bucket := cluster.Bucket(bucketName)

	collection := bucket.DefaultCollection()

	id := uuid.New().String()

	user := User{
		ID:       id,
		Username: username,
	}

	_, err = collection.Upsert(id, user, nil)
	if err != nil {
		return "", fmt.Errorf("error inserting document into Couchbase: %w", err)
	}

	fmt.Printf("Inserted Document ID: %s\n", id)

	return id, nil
}
