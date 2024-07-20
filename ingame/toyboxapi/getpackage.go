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
	"fmt"
	"log"
	"net/http"
	"reboxed/db"
	"reboxed/utils"
	"strconv"
)

// getpackage returns package metadata
func GetPackage(w http.ResponseWriter, r *http.Request) {
	if !utils.ValidateKey(r.URL.String()) {
		utils.WriteError(w, r, "invalid key")
		return
	}

	scriptid, err := strconv.Atoi(utils.UnBinHexString(r.FormValue("scriptid")))
	if err != nil {
		utils.WriteError(w, r, fmt.Sprintf("failed to parse scriptid value: %s", err))
		return
	}

	rev, err := strconv.Atoi(utils.UnBinHexString(r.FormValue("rev")))
	if err != nil {
		utils.WriteError(w, r, fmt.Sprintf("failed to parse rev value: %s", err))
		return
	}

	// getscript also specifies "type" but scriptids are unique between types
	// we don't need its value because of this

	pkg, err := db.FetchPackage(scriptid, rev)
	if err != nil {
		utils.WriteError(w, r, fmt.Sprintf("failed to fetch package: %s", err))
		return
	}

	log.Printf("New package download: %s/%d (%s)", pkg.Type, pkg.ID, pkg.Name)

	w.Write(pkg.Marshal())
}
