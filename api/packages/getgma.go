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

package packages

import (
	"database/sql"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/flatgrassdotnet/cloudbox/common"
	"github.com/flatgrassdotnet/cloudbox/db"
	"github.com/flatgrassdotnet/cloudbox/utils"
)

type gmaDescription struct {
	Description string   `json:"description"`
	Type        string   `json:"type"`
	Tags        []string `json:"tags"`
}

func GetGMA(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		utils.WriteError(w, r, fmt.Sprintf("failed to parse id value: %s", err))
		return
	}

	rev, _ := strconv.Atoi(r.URL.Query().Get("rev"))
	if rev < 1 {
		rev, err = db.FetchPackageLatestRevision(id)
		if err != nil {
			utils.WriteError(w, r, fmt.Sprintf("failed to fetch package latest revision: %s", err))
			return
		}
	}

	pkg, err := db.FetchPackage(id, rev)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "package not found", http.StatusNotFound)
			return
		}

		utils.WriteError(w, r, fmt.Sprintf("failed to fetch package: %s", err))
		return
	}

	if len(pkg.Content) == 0 {
		w.WriteHeader(http.StatusOK)
		return
	}

	// magic
	w.Write([]byte("GMAD"))

	// version
	binary.Write(w, binary.LittleEndian, uint8(3))

	// steamid (unused)
	var author int
	if pkg.Author != "" {
		author, err = strconv.Atoi(pkg.Author)
		if err != nil {
			utils.WriteError(w, r, fmt.Sprintf("failed to convert author steamid: %s", err))
			return
		}
	}

	binary.Write(w, binary.LittleEndian, uint64(author))

	// timestamp
	binary.Write(w, binary.LittleEndian, uint64(pkg.Uploaded.Unix()))

	// required content (stubbed)
	binary.Write(w, binary.LittleEndian, uint8(0))

	// addon name
	w.Write([]byte(pkg.Name + "\000"))

	// addon description
	err = json.NewEncoder(w).Encode(gmaDescription{
		Description: pkg.Description,
		Type:        pkg.Type,
		Tags:        []string{"fun"},
	})
	if err != nil {
		utils.WriteError(w, r, fmt.Sprintf("failed to marshal package description: %s", err))
		return
	}

	binary.Write(w, binary.LittleEndian, uint8(0)) // null terminator

	// addon author (unused)
	w.Write([]byte(pkg.AuthorName + "\000"))

	// addon version (unused)
	binary.Write(w, binary.LittleEndian, uint32(pkg.Revision))

	// exclude non-whitelisted files
	var content []common.Content
	for _, item := range pkg.Content {
		whitelisted, err := isPathWhitelisted(item.Path)
		if err != nil {
			utils.WriteError(w, r, fmt.Sprintf("failed to check if path is whitelisted: %s", err))
			return
		}

		if whitelisted {
			content = append(content, item)
		}
	}

	// file list
	for i, item := range content {
		// file number
		binary.Write(w, binary.LittleEndian, uint32(i+1))

		// file name
		w.Write([]byte(strings.ToLower(item.Path) + "\000"))

		// file size
		binary.Write(w, binary.LittleEndian, uint64(item.Size))

		// file crc (skipped)
		binary.Write(w, binary.LittleEndian, uint32(0))
	}

	// end of file list marker
	binary.Write(w, binary.LittleEndian, uint32(0))

	// file content
	for _, item := range content {
		f, err := utils.GetContentFile(item.ID)
		if err != nil {
			utils.WriteError(w, r, fmt.Sprintf("failed to get content file data: %s", err))
			return
		}

		defer f.Close()

		io.Copy(w, f)
	}

	// content crc (skipped)
	binary.Write(w, binary.LittleEndian, uint32(0))
}

var gmaWhitelist = map[string]bool{
	"^lua/(.*).lua$":                               true,
	"^scenes/(.*).vcd$":                            true,
	"^particles/(.*).pcf$":                         true,
	"^resource/fonts/(.*).ttf$":                    true,
	"^scripts/vehicles/(.*).txt$":                  true,
	"^resource/localization/(.*)/(.*).properties$": true,
	"^maps/(.*).bsp$":                              true,
	"^maps/(.*).lmp$":                              true,
	"^maps/(.*).nav$":                              true,
	"^maps/(.*).ain$":                              true,
	"^maps/thumb/(.*).png$":                        true,
	"^sound/(.*).wav$":                             true,
	"^sound/(.*).mp3$":                             true,
	"^sound/(.*).ogg$":                             true,
	"^materials/(.*).vmt$":                         true,
	"^materials/(.*).vtf$":                         true,
	"^materials/(.*).png$":                         true,
	"^materials/(.*).jpg$":                         true,
	"^materials/(.*).jpeg$":                        true,
	"^materials/colorcorrection/(.*).raw$":         true,
	"^models/(.*).mdl$":                            true,
	"^models/(.*).phy$":                            true,
	"^models/(.*).ani$":                            true,
	"^models/(.*).vvd$":                            true,

	"^models/(.*).vtx$":       true,
	"^!models/(.*).sw.vtx$":   false, // These variations are unused by the game
	"^!models/(.*).360.vtx$":  false,
	"^!models/(.*).xbox.vtx$": false,

	"^gamemodes/(.*)/(.*).txt$":       true,
	"^!gamemodes/(.*)/(.*)/(.*).txt$": false, // Only in the root gamemode folder please!
	"^gamemodes/(.*)/(.*).fgd$":       true,
	"^!gamemodes/(.*)/(.*)/(.*).fgd$": false,

	"^gamemodes/(.*)/logo.png$":                   true,
	"^gamemodes/(.*)/icon24.png$":                 true,
	"^gamemodes/(.*)/gamemode/(.*).lua$":          true,
	"^gamemodes/(.*)/entities/effects/(.*).lua$":  true,
	"^gamemodes/(.*)/entities/weapons/(.*).lua$":  true,
	"^gamemodes/(.*)/entities/entities/(.*).lua$": true,
	"^gamemodes/(.*)/backgrounds/(.*).png$":       true,
	"^gamemodes/(.*)/backgrounds/(.*).jpg$":       true,
	"^gamemodes/(.*)/backgrounds/(.*).jpeg$":      true,
	"^gamemodes/(.*)/content/models/(.*).mdl$":    true,
	"^gamemodes/(.*)/content/models/(.*).phy$":    true,
	"^gamemodes/(.*)/content/models/(.*).ani$":    true,
	"^gamemodes/(.*)/content/models/(.*).vvd$":    true,

	"^gamemodes/(.*)/content/models/(.*).vtx$":       true,
	"^!gamemodes/(.*)/content/models/(.*).sw.vtx$":   false,
	"^!gamemodes/(.*)/content/models/(.*).360.vtx$":  false,
	"^!gamemodes/(.*)/content/models/(.*).xbox.vtx$": false,

	"^gamemodes/(.*)/content/materials/(.*).vmt$":                         true,
	"^gamemodes/(.*)/content/materials/(.*).vtf$":                         true,
	"^gamemodes/(.*)/content/materials/(.*).png$":                         true,
	"^gamemodes/(.*)/content/materials/(.*).jpg$":                         true,
	"^gamemodes/(.*)/content/materials/(.*).jpeg$":                        true,
	"^gamemodes/(.*)/content/materials/colorcorrection/(.*).raw$":         true,
	"^gamemodes/(.*)/content/scenes/(.*).vcd$":                            true,
	"^gamemodes/(.*)/content/particles/(.*).pcf$":                         true,
	"^gamemodes/(.*)/content/resource/fonts/(.*).ttf$":                    true,
	"^gamemodes/(.*)/content/scripts/vehicles/(.*).txt$":                  true,
	"^gamemodes/(.*)/content/resource/localization/(.*)/(.*).properties$": true,
	"^gamemodes/(.*)/content/maps/(.*).bsp$":                              true,
	"^gamemodes/(.*)/content/maps/(.*).nav$":                              true,
	"^gamemodes/(.*)/content/maps/(.*).ain$":                              true,
	"^gamemodes/(.*)/content/maps/thumb/(.*).png$":                        true,
	"^gamemodes/(.*)/content/sound/(.*).wav$":                             true,
	"^gamemodes/(.*)/content/sound/(.*).mp3$":                             true,
	"^gamemodes/(.*)/content/sound/(.*).ogg$":                             true,

	// static version of the data/ folder
	// (because you wouldn't be able to modify these)
	// We only allow filetypes here that are not already allowed above
	"^data_static/(.*).txt$":  true,
	"^data_static/(.*).dat$":  true,
	"^data_static/(.*).json$": true,
	"^data_static/(.*).xml$":  true,
	"^data_static/(.*).csv$":  true,
}

func isPathWhitelisted(path string) (bool, error) {
	for rule, allowed := range gmaWhitelist {
		matched, err := regexp.MatchString(rule, path)
		if err != nil {
			return false, err
		}

		if matched {
			return allowed, nil
		}
	}

	return false, nil
}
