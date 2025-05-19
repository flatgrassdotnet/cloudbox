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
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/flatgrassdotnet/cloudbox/db"
	"github.com/flatgrassdotnet/cloudbox/utils"
)

func GetZIP(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		utils.WriteError(w, r, fmt.Sprintf("failed to parse id value: %s", err))
		return
	}

	o, err := db.GetContentFile(id)
	if err != nil {
		utils.WriteError(w, r, fmt.Sprintf("failed to open content file for reading: %s", err))
		return
	}

	defer o.Body.Close()

	// GM12 won't show download progress without Content-Length
	buf := new(bytes.Buffer)

	zw := zip.NewWriter(buf)

	file, err := zw.Create("file")
	if err != nil {
		utils.WriteError(w, r, fmt.Sprintf("failed to open zip file for writing: %s", err))
		return
	}

	io.Copy(file, o.Body)

	zw.Close()

	w.Header().Set("Content-Type", "application/zip")
	w.Header().Set("Content-Length", strconv.Itoa(buf.Len()))
	io.Copy(w, buf)
}
