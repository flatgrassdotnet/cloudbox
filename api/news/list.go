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

package news

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/flatgrassdotnet/cloudbox/db"
	"github.com/flatgrassdotnet/cloudbox/utils"
)

func List(w http.ResponseWriter, r *http.Request) {
	entries, err := db.FetchNewsEntries()
	if err != nil {
		utils.WriteError(w, r, fmt.Sprintf("failed to fetch news entries: %s", err))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(entries)
	if err != nil {
		utils.WriteError(w, r, fmt.Sprintf("failed to encode response: %s", err))
		return
	}
}
