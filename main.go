package main

import (
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/nlopes/slack"
)

func main() {
	os.Exit(_main(os.Args[1:]))
}

func _main(args []string) int {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Listening slack event and response
	log.Printf("[INFO] Start slack event listening")
	client := slack.New(os.Getenv("BOT_TOKEN"))
	slackListener := &SlackListener{
		client:    client,
		botID:     os.Getenv("BOT_ID"),
		channelID: os.Getenv("CHANNEL_ID"),
	}
	go slackListener.ListenAndResponse()

	// Register handler to receive interactive message
	// responses from slack (kicked by user action)
	http.Handle("/interaction", interactionHandler{
		verificationToken: os.Getenv("VERIFICATION_TOKEN"),
		slackClient:       client,
	})

	const port = "3000"
	log.Printf("[INFO] Server listening on :%s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Printf("[ERROR] %s", err)
		return 1
	}

	return 0
}
