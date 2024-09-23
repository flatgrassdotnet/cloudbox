/*
	cloudbox - the toybox server emulator
	Copyright (C) 2024  patapancakes <patapancakes@pagefault.games>

	This program is free software: you can redistribute it and/or modify
	it under the terms of the GNU Affero General Public License as published by
	the Free Software Foundation, either version 3 of the License, or
	(at your option) any later version.

	This program is distributed in the hope that it will be useful,
	but WITHOUT ANY WARRANTY; without even the implied warranty of
	MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
	GNU Affero General Public License for more details.

	You should have received a copy of the GNU Affero General Public License
	along with this program.  If not, see <http://www.gnu.org/licenses/>.
*/

package utils

import (
	"bytes"
	"encoding/json"
	"net/http"
)

var (
	DiscordSaveWebhookURL  string
	DiscordStatsWebhookURL string
)

type DiscordWebhookRequest struct {
	Embeds          []DiscordWebhookEmbed `json:"embeds"`
	AllowedMentions struct {
		Parse []string `json:"parse"`
	} `json:"allowed_mentions"`
}

type DiscordWebhookEmbed struct {
	Title       string                    `json:"title"`
	Description string                    `json:"description"`
	Color       int                       `json:"color"`
	Author      DiscordWebhookEmbedAuthor `json:"author"`
	Image       DiscordWebhookEmbedImage  `json:"image"`
}

type DiscordWebhookEmbedAuthor struct {
	Name    string `json:"name"`
	IconURL string `json:"icon_url"`
}

type DiscordWebhookEmbedImage struct {
	URL string `json:"url"`
}

func SendDiscordMessage(url string, data DiscordWebhookRequest) error {
	if url == "" {
		return nil
	}

	body, err := json.Marshal(data)
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
