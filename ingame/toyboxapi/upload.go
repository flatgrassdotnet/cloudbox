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
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"reboxed/common"
	"reboxed/db"
	"reboxed/utils"
	"strconv"
)

func Upload(w http.ResponseWriter, r *http.Request) {
	// steamid64
	steamid, err := strconv.Atoi(r.URL.Query().Get("steamid"))
	if err != nil {
		utils.WriteError(w, r, fmt.Sprintf("failed to parse steamid value: %s", err))
		return
	}

	// "save" or "save_image"
	uploadType := r.URL.Query().Get("type")
	if uploadType != "save" && uploadType != "save_image" {
		utils.WriteError(w, r, "invalid upload type")
		return
	}

	// metadata
	meta := r.URL.Query().Get("meta")

	// unknown
	inc := r.URL.Query().Get("inc")

	body, err := io.ReadAll(r.Body)
	if err != nil {
		utils.WriteError(w, r, fmt.Sprintf("failed to read post body: %s", err))
		return
	}

	decoded := make([]byte, base64.StdEncoding.DecodedLen(len(body)))
	_, err = base64.StdEncoding.Decode(decoded, body)
	if err != nil {
		utils.WriteError(w, r, fmt.Sprintf("failed to parse post body: %s", err))
		return
	}

	id, err := db.InsertUpload(steamid, common.Upload{Type: uploadType, Metadata: meta, Include: inc, Data: decoded})
	if err != nil {
		utils.WriteError(w, r, fmt.Sprintf("failed to insert upload: %s", err))
		return
	}

	w.Write([]byte(strconv.Itoa(id)))
}
