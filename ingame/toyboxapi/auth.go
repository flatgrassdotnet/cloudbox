/*
	cloudbox - the toybox server emulator
	Copyright (C) 2024-2025  patapancakes <patapancakes@pagefault.games>

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
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"

	"github.com/flatgrassdotnet/cloudbox/db"
	"github.com/flatgrassdotnet/cloudbox/utils"
)

// auth logs someone into the toybox api
func Auth(w http.ResponseWriter, r *http.Request) {
	// the server returns a 32 character login token on success
	// any value under 32 characters is treated as an error response
	// the server responding "chrome" will make an anti piracy message appear

	err := r.ParseForm()
	if err != nil {
		utils.WriteError(w, r, fmt.Sprintf("failed to parse form data: %s", err))
		return
	}

	if !utils.ValidateKey(r.Form.Encode()) {
		utils.WriteError(w, r, "invalid key")
		return
	}

	// "u" value (steamid64) is ignored - we get it from steam
	// "vac" value is ignored - we get it from steam

	user, err := utils.AuthenticateUserTicket(utils.UnBinHexString(r.FormValue("token")))
	if err != nil {
		// net/http errors shouldn't cause the game to exit
		if !strings.Contains(err.Error(), "net/http:") {
			w.Write([]byte("chrome")) // terminate game with anti-piracy error
		}

		utils.WriteError(w, r, fmt.Sprintf("failed to validate steam ticket: %s", err))
		return
	}

	vac := "good"
	if user.VACBanned {
		vac = "banned"
	}

	ticket := make([]byte, 24)
	_, err = rand.Read(ticket)
	if err != nil {
		utils.WriteError(w, r, fmt.Sprintf("failed to generate ticket: %s", err))
		return
	}

	err = db.InsertLogin(user.SteamID, vac, ticket)
	if err != nil {
		utils.WriteError(w, r, fmt.Sprintf("failed to insert login: %s", err))
		return
	}

	w.Write([]byte(base64.StdEncoding.EncodeToString(ticket)))

	// webhook related
	s, err := utils.GetPlayerSummary(user.SteamID)
	if err != nil {
		utils.WriteError(w, r, fmt.Sprintf("failed to get player summary: %s", err))
		return
	}

	err = utils.SendDiscordMessage(utils.DiscordStatsWebhookURL, utils.DiscordWebhookRequest{
		Embeds: []utils.DiscordWebhookEmbed{{
			Title: "Login",
			Color: 0x4096EE,
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
