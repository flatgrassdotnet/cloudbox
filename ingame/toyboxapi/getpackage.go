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

// getpackage returns package metadata
func GetPackage(w http.ResponseWriter, r *http.Request) {
	if !utils.ValidateKey(r.URL.String()) {
		utils.WriteError(w, r, "invalid key")
		return
	}

	steamid, err := strconv.Atoi(utils.UnBinHexString(r.FormValue("u")))
	if err != nil {
		utils.WriteError(w, r, fmt.Sprintf("failed to parse u value: %s", err))
		return
	}

	scriptid, err := strconv.Atoi(utils.UnBinHexString(r.FormValue("scriptid")))
	if err != nil {
		utils.WriteError(w, r, fmt.Sprintf("failed to parse scriptid value: %s", err))
		return
	}

	rev, err := strconv.Atoi(utils.UnBinHexString(r.FormValue("rev")))
	if err != nil {
		utils.WriteError(w, r, fmt.Sprintf("failed to parse rev value: %s", err))
		return
	}

	// getscript also specifies "type" but scriptids are unique between types
	// we don't need its value because of this

	pkg, err := db.FetchPackage(scriptid, rev)
	if err != nil {
		utils.WriteError(w, r, fmt.Sprintf("failed to fetch package: %s", err))
		return
	}

	w.Write(pkg.Marshal())

	// webhook related
	s, err := utils.GetPlayerSummary(uint64(steamid))
	if err != nil {
		utils.WriteError(w, r, fmt.Sprintf("failed to get player summary: %s", err))
		return
	}

	err = utils.SendDiscordMessage(utils.DiscordStatsWebhookURL, utils.DiscordWebhookRequest{
		Embeds: []utils.DiscordWebhookEmbed{{
			Title:       "Package Download",
			Description: fmt.Sprintf("%s (%dr%d/%s)", pkg.Name, pkg.ID, pkg.Revision, pkg.Type),
			Color:       0x4096EE,
			Author: utils.DiscordWebhookEmbedAuthor{
				Name:    s.PersonaName,
				IconURL: s.Avatar,
			},
			Image: utils.DiscordWebhookEmbedImage{
				URL: fmt.Sprintf("https://img.reboxed.fun/%d_thumb_128.png", pkg.ID),
			},
		}},
	})
	if err != nil {
		utils.WriteError(w, r, fmt.Sprintf("failed to send discord webhook message: %s", err))
		return
	}
}
