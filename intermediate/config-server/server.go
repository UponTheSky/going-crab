package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
)

type Secrets struct {
	ActorName string `json:"actor_name"`
	Ethnicity string `json:"ethnicity"`
}

func readSecretsFromJson(path string) (Secrets, error) {
	secretBytes, err := os.ReadFile(path)

	if err != nil && !errors.Is(err, io.EOF) {
		return Secrets{}, err
	}

	secrets := Secrets{}
	json.Unmarshal(secretBytes, &secrets)

	return secrets, nil
}

func main() {
	// environment variable
	// if secret, ok := os.LookupEnv("secret"); ok {
	// 	fmt.Printf("secret: %v\n", secret)
	// }

	// from file
	secretsPath := "/tmp/gcrab_secrets.json"
	secrets, err := readSecretsFromJson(secretsPath)

	if err != nil {
		log.Fatalf("reading json data from %v", secrets)
	}

	fmt.Printf("The secrets here: actor name - %v, ethnicity - %v", secrets.ActorName, secrets.Ethnicity)
}
