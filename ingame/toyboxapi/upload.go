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
	"bytes"
	"encoding/base64"
	"encoding/csv"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/flatgrassdotnet/cloudbox/common"
	"github.com/flatgrassdotnet/cloudbox/db"
	"github.com/flatgrassdotnet/cloudbox/utils"
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

	// includes
	var includes []int
	if r.URL.Query().Get("inc") != "" {
		incs, err := csv.NewReader(bytes.NewReader([]byte(r.URL.Query().Get("inc")))).Read()
		if err != nil {
			utils.WriteError(w, r, fmt.Sprintf("failed to decode inc: %s", err))
			return
		}

		for _, inc := range incs {
			// ignore blank entries
			if inc == "" {
				continue
			}

			id, err := strconv.Atoi(inc)
			if err != nil {
				utils.WriteError(w, r, fmt.Sprintf("failed to decode inc value: %s", err))
				return
			}

			includes = append(includes, id)
		}
	}

	body, err := io.ReadAll(base64.NewDecoder(base64.StdEncoding, r.Body))
	if err != nil {
		utils.WriteError(w, r, fmt.Sprintf("failed to decode request body: %s", err))
		return
	}

	id, err := db.InsertUpload(steamid, common.Upload{Type: uploadType, Metadata: meta, Includes: includes, Data: body})
	if err != nil {
		utils.WriteError(w, r, fmt.Sprintf("failed to insert upload: %s", err))
		return
	}

	fmt.Fprint(w, id)
}
