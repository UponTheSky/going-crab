package main

import (
	"encoding/json"
	"flag"
	"log"
	"os"
)

type Secrets struct {
	ActorName string `json:"actor_name"`
	Ethnicity string `json:"ethnicity"`
}

func main() {
	flagSet := flag.NewFlagSet("gensecrets", flag.ExitOnError)

	actorName := flagSet.String("actorName", "", "-actorName \"Ken Jeong\"")
	ethnicity := flagSet.String("ethnicity", "", "-ethnicity \"Korean\"")

	if err := flagSet.Parse(os.Args[1:]); err != nil {
		log.Fatalf("parsing error: -actorName and -ethnicity flags must be provided: %v", err)
	}

	secrets := Secrets{ActorName: *actorName, Ethnicity: *ethnicity}

	encoded, err := json.Marshal(secrets)

	if err != nil {
		log.Fatalf("unexpected Json encoding error: %v", err)
	}

	if err := os.WriteFile("/tmp/gcrab_secrets.json", encoded, 0644); err != nil {
		log.Fatalf("unexpected Json file writing error: %v", err)
	}
}
