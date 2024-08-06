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

package toyboxapi

import (
	"fmt"
	"net/http"
	"reboxed/db"
	"reboxed/utils"
	"strconv"
)

// error records errors from users
func Error(w http.ResponseWriter, r *http.Request) {
	if !utils.ValidateKey(r.URL.String()) {
		utils.WriteError(w, r, "invalid key")
		return
	}

	// game version
	version, err := strconv.Atoi(utils.UnBinHexString(r.FormValue("v")))
	if err != nil {
		utils.WriteError(w, r, fmt.Sprintf("failed to parse v value: %s", err))
		return
	}

	// steamid64
	steamid, err := strconv.Atoi(utils.UnBinHexString(r.FormValue("u")))
	if err != nil {
		utils.WriteError(w, r, fmt.Sprintf("failed to parse u value: %s", err))
		return
	}

	// error
	error := utils.UnBinHexString(r.URL.Query().Get("error"))

	// content
	content := utils.UnBinHexString(r.URL.Query().Get("content"))

	// realm ("client", or "server")
	realm := utils.UnBinHexString(r.URL.Query().Get("realm"))
	if !(realm == "client" || realm == "server") {
		utils.WriteError(w, r, "invalid realm value")
		return
	}

	// platform ("win32", "linux", or "osx")
	platform := utils.UnBinHexString(r.URL.Query().Get("platform"))
	if !(platform == "win32" || platform == "linux" || platform == "osx") {
		utils.WriteError(w, r, "invalid platform value")
		return
	}

	err = db.InsertError(version, steamid, error, content, realm, platform)
	if err != nil {
		utils.WriteError(w, r, fmt.Sprintf("failed to insert error: %s", err))
		return
	}

	w.WriteHeader(http.StatusOK)

	// webhook related
	s, err := utils.GetPlayerSummary(int64(steamid))
	if err != nil {
		utils.WriteError(w, r, fmt.Sprintf("failed to get player summary: %s", err))
		return
	}

	err = utils.SendDiscordMessage(utils.DiscordStatsWebhookURL, utils.DiscordWebhookRequest{
		Embeds: []utils.DiscordWebhookEmbed{{
			Title:       "Error",
			Description: error,
			Color:       4232942, // #4096EE
			Author: utils.DiscordWebhookEmbedAuthor{
				Name:    s.PersonaName,
				IconURL: s.Avatar,
			},
		}},
	})
	if err != nil {
		utils.WriteError(w, r, fmt.Sprintf("failed to send discord webhook message: %s", err))
		return
	}
}
