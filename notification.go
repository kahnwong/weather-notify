package main

import (
	"context"
	"os"

	"github.com/carlmjohnson/requests"
)

var (
	discordWebhookUrl = os.Getenv("DISCORD_WEBHOOK_URL")
)

type discordWebhook struct {
	Content string `json:"content"`
}

func notify(outputMessage string) error {
	body := discordWebhook{
		Content: outputMessage,
	}

	err := requests.
		URL(discordWebhookUrl).
		BodyJSON(&body).
		Fetch(context.Background())

	return err
}
