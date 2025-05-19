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

package auth

import (
	"encoding/base64"
	"fmt"
	"net/http"

	"github.com/flatgrassdotnet/cloudbox/db"
	"github.com/flatgrassdotnet/cloudbox/utils"
)

func GetID(w http.ResponseWriter, r *http.Request) {
	ticket, err := base64.StdEncoding.DecodeString(r.URL.Query().Get("ticket"))
	if err != nil {
		utils.WriteError(w, r, fmt.Sprintf("failed to decode ticket value: %s", err))
		return
	}

	steamid, err := db.FetchSteamIDFromTicket(ticket)
	if err != nil {
		utils.WriteError(w, r, fmt.Sprintf("failed to fetch steamid from ticket: %s", err))
		return
	}

	fmt.Fprint(w, steamid)
}
