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
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/flatgrassdotnet/cloudbox/utils"
)

func Get(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		utils.WriteError(w, r, fmt.Sprintf("failed to parse id value: %s", err))
		return
	}

	rev, _ := strconv.Atoi(r.URL.Query().Get("rev"))
	if rev < 1 {
		rev = 1
	}

	f, zf, err := utils.GetContentFile(id, rev)
	if err != nil {
		utils.WriteError(w, r, fmt.Sprintf("failed to get content file data: %s", err))
		return
	}

	defer f.Close()
	defer zf.Close()

	io.Copy(w, f)
}
