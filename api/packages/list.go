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

package packages

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reboxed/db"
	"reboxed/utils"
)

func List(w http.ResponseWriter, r *http.Request) {
	list, err := db.FetchPackageList(r.URL.Query().Get("type"))
	if err != nil {
		utils.WriteError(w, r, fmt.Sprintf("failed to fetch package list: %s", err))
		return
	}

	resp, err := json.Marshal(list)
	if err != nil {
		utils.WriteError(w, r, fmt.Sprintf("failed to marshal response: %s", err))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(resp)
}
