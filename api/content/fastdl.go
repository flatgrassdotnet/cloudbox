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

package content

import (
	"database/sql"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/flatgrassdotnet/cloudbox/db"
	"github.com/flatgrassdotnet/cloudbox/utils"
)

func FastDL(w http.ResponseWriter, r *http.Request) {
	file := r.URL.Query().Get("file")
	if file == "" {
		utils.WriteError(w, r, "missing file value")
		return
	}

	id, err := db.FetchFileInfoFromPath(strings.TrimPrefix(file, "/"))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "file not found", http.StatusNotFound)
			return
		}

		utils.WriteError(w, r, fmt.Sprintf("failed to fetch file info: %s", err))
		return
	}

	f, err := db.GetContentFile(id)
	if err != nil {
		utils.WriteError(w, r, fmt.Sprintf("failed to open content file: %s", err))
		return
	}

	defer f.Body.Close()

	w.Header().Set("Content-Length", strconv.Itoa(int(*f.ContentLength)))
	io.Copy(w, f.Body)
}
