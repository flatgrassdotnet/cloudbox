/*
	cloudbox - the toybox server emulator
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
	"strconv"

	"github.com/flatgrassdotnet/cloudbox/db"
	"github.com/flatgrassdotnet/cloudbox/utils"
)

func List(w http.ResponseWriter, r *http.Request) {
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
	if offset < 0 {
		offset = 0
	}

	count, _ := strconv.Atoi(r.URL.Query().Get("count"))
	if count < 0 {
		count = 0
	}
	if count > 100 {
		count = 100
	}

	var sort string // must NOT be user input
	switch r.URL.Query().Get("sort") {
	case "mostfavs":
		sort = "favorites"
	case "mostlikes":
		sort = "goods"
	case "mostdls":
		sort = "downloads"
	case "random":
		sort = "RAND()"
	default: // newest
		sort = "id"
	}

	list, err := db.FetchPackageList(r.URL.Query().Get("type"), r.URL.Query().Get("author"), r.URL.Query().Get("search"), offset, count, sort)
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
