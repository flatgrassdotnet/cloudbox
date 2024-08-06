package utils

import (
	"bytes"
	"encoding/json"
	"net/http"
)

var DiscordWebhookURL string

type DiscordWebhookRequest struct {
	Username        string `json:"username"`
	AvatarURL       string `json:"avatar_url"`
	Content         string `json:"content"`
	AllowedMentions struct {
		Parse []string `json:"parse"`
	} `json:"allowed_mentions"`
}

func SendDiscordMessage(url string, steamid int64, content string) error {
	if url == "" {
		return nil
	}

	u, err := GetPlayerSummary(steamid)
	if err != nil {
		return err
	}

	body, err := json.Marshal(DiscordWebhookRequest{
		Username: u.PersonaName,
		AvatarURL: u.Avatar,
		Content: content,
	})
	if err != nil {
		return err
	}

	resp, err := http.Post(url, "application/json", bytes.NewReader(body))
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	return nil
}
