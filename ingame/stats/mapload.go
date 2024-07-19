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

package stats

import (
	"fmt"
	"log"
	"net/http"
	"reboxed/db"
	"reboxed/utils"
	"strconv"
)

// mapload records statistics about map usage
func MapLoad(w http.ResponseWriter, r *http.Request) {
	if !utils.ValidateKey(r.URL.String()) {
		utils.WriteError(w, r, "invalid key")
		return
	}

	// game version
	version, err := strconv.Atoi(utils.UnBinHexString(r.FormValue("v")))
	if err != nil {
		utils.WriteError(w, r, fmt.Sprintf("failed to parse v value: %s", err))
		return
	}

	// steamid64
	steamid, err := strconv.Atoi(utils.UnBinHexString(r.FormValue("u")))
	if err != nil {
		utils.WriteError(w, r, fmt.Sprintf("failed to parse u value: %s", err))
		return
	}

	// duration taken to load (in seconds)
	duration, err := strconv.ParseFloat(utils.UnBinHexString(r.URL.Query().Get("time")), 32)
	if err != nil {
		utils.WriteError(w, r, fmt.Sprintf("failed to parse time value: %s", err))
		return
	}

	// map
	mapName := utils.UnBinHexString(r.URL.Query().Get("map"))

	// platform ("win32", "linux", or "osx")
	platform := utils.UnBinHexString(r.URL.Query().Get("platform"))
	if !(platform == "win32" || platform == "linux" || platform == "osx") {
		utils.WriteError(w, r, "invalid platform value")
		return
	}

	log.Printf("New map load: v=%d, u=%d, time=%f, map=%s, platform=%s", version, steamid, duration, mapName, platform)

	err = db.InsertMapLoad(version, steamid, duration, mapName, platform)
	if err != nil {
		utils.WriteError(w, r, fmt.Sprintf("failed to insert map load: %s", err))
		return
	}

	w.WriteHeader(http.StatusOK)
}
