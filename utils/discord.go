/*
	reboxed - the toybox server emulator
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
		Username:  u.PersonaName,
		AvatarURL: u.Avatar,
		Content:   content,
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
