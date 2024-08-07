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
	"fmt"
	"log"
	"net/http"
	"strconv"
)

func WriteError(w http.ResponseWriter, r *http.Request, message string) {
	log.Printf("%s: %s", r.URL, message)
	w.WriteHeader(http.StatusBadRequest)

	// webhook related
	var s PlayerSummaryInfo
	steamid, err := strconv.Atoi(UnBinHexString(r.FormValue("u")))
	if err == nil {
		s, _ = GetPlayerSummary(uint64(steamid))
	}

	SendDiscordMessage(DiscordStatsWebhookURL, DiscordWebhookRequest{
		Embeds: []DiscordWebhookEmbed{{
			Title:       "API Error",
			Description: fmt.Sprintf("%s: %s", r.URL, message),
			Color:       0x7D0000,
			Author: DiscordWebhookEmbedAuthor{
				Name:    s.PersonaName,
				IconURL: s.Avatar,
			},
		}},
	})
}
