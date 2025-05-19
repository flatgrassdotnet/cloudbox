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
	"slices"
	"strconv"

	"github.com/flatgrassdotnet/cloudbox/db"
	"github.com/flatgrassdotnet/cloudbox/utils"

	appticket "github.com/tmcarey/steam-appticket-go"
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

	token := utils.UnBinHex(utils.UnBinHexString(r.FormValue("token")))

	steamid := utils.UnBinHexString(r.FormValue("u"))

	vac := utils.UnBinHexString(r.FormValue("vac"))
	if !slices.Contains([]string{"good", "banned"}, vac) {
		utils.WriteError(w, r, "invalid vac value")
		return
	}

	appticket, err := appticket.ParseAppTicket(token, false)
	if err != nil || !appticket.IsValid || appticket.AppID != 4000 || strconv.Itoa(int(appticket.SteamID)) != steamid {
		fmt.Fprint(w, "chrome") // terminate game with anti-piracy error
		return
	}

	// store new profile or get its data
	s, err := utils.GetPlayerSummaries(steamid)
	if err != nil {
		utils.WriteError(w, r, fmt.Sprintf("failed to get player summary: %s", err))
		return
	}

	ticket := make([]byte, 24)
	_, err = rand.Read(ticket)
	if err != nil {
		utils.WriteError(w, r, fmt.Sprintf("failed to generate ticket: %s", err))
		return
	}

	err = db.InsertLogin(steamid, vac, ticket)
	if err != nil {
		utils.WriteError(w, r, fmt.Sprintf("failed to insert login: %s", err))
		return
	}

	base64.NewEncoder(base64.StdEncoding, w).Write(ticket)

	// webhook related
	err = utils.SendDiscordMessage(utils.DiscordStatsWebhookURL, utils.DiscordWebhookRequest{
		Embeds: []utils.DiscordWebhookEmbed{{
			Title: "Login",
			Color: 0x4096EE,
			Author: utils.DiscordWebhookEmbedAuthor{
				Name:    s[0].PersonaName,
				IconURL: s[0].Avatar,
			},
		}},
	})
	if err != nil {
		utils.WriteError(w, r, fmt.Sprintf("failed to send discord webhook message: %s", err))
		return
	}
}
