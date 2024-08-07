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

package content

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"reboxed/utils"
	"strconv"
)

func Get(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		utils.WriteError(w, r, fmt.Sprintf("failed to parse id value: %s", err))
		return
	}

	rev, err := strconv.Atoi(r.URL.Query().Get("rev"))
	if err != nil {
		rev = 1
	}

	b, err := os.ReadFile(fmt.Sprintf("data/cdn/%d/%d", id, rev))
	if err != nil {
		utils.WriteError(w, r, fmt.Sprintf("failed to read zip: %s", err))
		return
	}

	zr, err := zip.NewReader(bytes.NewReader(b), int64(len(b)))
	if err != nil {
		utils.WriteError(w, r, fmt.Sprintf("failed to create zip reader: %s", err))
		return
	}

	f, err := zr.Open("file")
	if err != nil {
		utils.WriteError(w, r, fmt.Sprintf("failed to open file from zip: %s", err))
		return
	}

	defer f.Close()

	fb, err := io.ReadAll(f)
	if err != nil {
		utils.WriteError(w, r, fmt.Sprintf("failed to read file from zip: %s", err))
		return
	}

	w.Write(fb)
}
