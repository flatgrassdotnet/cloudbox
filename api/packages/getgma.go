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
	"encoding/binary"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strconv"

	"github.com/flatgrassdotnet/cloudbox/common"
	"github.com/flatgrassdotnet/cloudbox/db"
	"github.com/flatgrassdotnet/cloudbox/utils"
)

type GMADescription struct {
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
		utils.WriteError(w, r, fmt.Sprintf("failed to fetch package: %s", err))
		return
	}

	if len(pkg.Content) == 0 {
		w.WriteHeader(http.StatusOK)
		return
	}

	// magic
	w.Write([]byte("GMAD"))

	// gma version
	w.Write([]byte{3})

	// steamid (unused)
	var author int
	if pkg.Author != "" {
		author, err = strconv.Atoi(pkg.Author)
		if err != nil {
			utils.WriteError(w, r, fmt.Sprintf("failed to convert author steamid: %s", err))
			return
		}
	}

	steamid := make([]byte, 8)
	binary.LittleEndian.PutUint64(steamid, uint64(author))
	w.Write(steamid)

	// timestamp (unused)
	timestamp := make([]byte, 8)
	binary.LittleEndian.PutUint64(timestamp, uint64(pkg.Uploaded.Unix()))
	w.Write(timestamp)

	// required content (unused)
	w.Write([]byte{0x00})

	// addon name
	w.Write(append([]byte(pkg.Name), 0))

	// addon description
	err = json.NewEncoder(w).Encode(GMADescription{
		Description: pkg.Description,
		Type:        pkg.Type,
		Tags:        []string{"fun"},
	})
	if err != nil {
		utils.WriteError(w, r, fmt.Sprintf("failed to marshal package description: %s", err))
		return
	}

	w.Write([]byte{0x00})

	// addon author
	w.Write(append([]byte(pkg.AuthorName), 0x00))

	// addon version
	version := make([]byte, 4)
	binary.LittleEndian.PutUint32(version, uint32(pkg.Revision))

	w.Write(version)

	var content []common.Content
	for _, item := range pkg.Content {
		whitelisted, err := IsGMAPathWhitelisted(item.Path)
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
		fileNum := make([]byte, 4)
		binary.LittleEndian.PutUint32(fileNum, uint32(i+1))
		w.Write(fileNum)

		// file name
		w.Write(append([]byte(item.Path), 0x00))

		// file size
		fileSize := make([]byte, 8)
		binary.LittleEndian.PutUint64(fileSize, uint64(item.Size))
		w.Write(fileSize)

		// file crc (skipped)
		w.Write(make([]byte, 4))
	}

	// end of file list marker
	w.Write(make([]byte, 4))

	// file content
	for _, item := range content {
		data, err := utils.GetContentFile(item.ID, item.Revision)
		if err != nil {
			utils.WriteError(w, r, fmt.Sprintf("failed to get content file data: %s", err))
			return
		}

		w.Write(data)
	}

	// content crc (skipped)
	w.Write(make([]byte, 4))
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

func IsGMAPathWhitelisted(path string) (bool, error) {
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
