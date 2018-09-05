package main

import (
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/nlopes/slack"
)

type envConfig struct {
	// Port is server port to be listened.
	Port string `envconfig:"PORT" default:"3000"`

	// BotToken is bot user token to access to slack API.
	BotToken string `envconfig:"BOT_TOKEN" required:"true"`

	// VerificationToken is used to validate interactive messages from slack.
	VerificationToken string `envconfig:"VERIFICATION_TOKEN" required:"true"`

	// BotID is bot user ID.
	BotID string `envconfig:"BOT_ID" required:"true"`

	// ChannelID is slack channel ID where bot is working.
	// Bot responses to the mention in this channel.
	ChannelID string `envconfig:"CHANNEL_ID" required:"true"`
}

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
		// verificationToken: os.Getenv("VARIFICATION_TOKEN"),
		verificationToken: "JxCpnDjGI9QjRVNdHZnlRe9V",
	})

	// log.Printf("[INFO] Server listening on :%s", env.Port)
	const port = "3000"
	log.Printf("[INFO] Server listening on :%s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Printf("[ERROR] %s", err)
		return 1
	}

	return 0
}
