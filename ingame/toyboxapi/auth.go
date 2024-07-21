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
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"reboxed/db"
	"reboxed/utils"
	"strconv"
	"strings"
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

	// game version
	version, err := strconv.Atoi(utils.UnBinHexString(r.FormValue("v")))
	if err != nil {
		utils.WriteError(w, r, fmt.Sprintf("failed to parse v value: %s", err))
		return
	}

	// "u" value (steamid64) is ignored - we get it from steam
	// "vac" value is ignored - we get it from steam

	user, err := utils.GetSteamUserInfo(utils.UnBinHexString(r.FormValue("token")))
	if err != nil {
		// net/http errors shouldn't cause the game to exit
		if !strings.Contains(err.Error(), "net/http:") {
			w.Write([]byte("chrome")) // terminate game with anti-piracy error
		}

		utils.WriteError(w, r, fmt.Sprintf("failed to validate steam ticket: %s", err))
		return
	}

	steamid, _ := strconv.Atoi(user.SteamID)

	vac := "good"
	if user.VACBanned {
		vac = "banned"
	}

	log.Printf("New login: v=%d, u=%d, vac=%s", version, steamid, vac)

	ticket := make([]byte, 24)
	_, err = rand.Read(ticket)
	if err != nil {
		utils.WriteError(w, r, fmt.Sprintf("failed to generate ticket: %s", err))
		return
	}

	err = db.InsertLogin(version, steamid, vac, ticket)
	if err != nil {
		utils.WriteError(w, r, fmt.Sprintf("failed to insert login: %s", err))
		return
	}

	w.Write([]byte(base64.StdEncoding.EncodeToString(ticket)))
}