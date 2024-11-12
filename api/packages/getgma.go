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
	"strconv"

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
		utils.WriteError(w, r, fmt.Sprintf("requested package with no content: %dr%d", id, rev))
		return
	}

	w.Header().Set("X-Package-ID", strconv.Itoa(pkg.ID))
	w.Header().Set("X-Package-Revision", strconv.Itoa(pkg.Revision))
	w.Header().Set("X-Package-Type", pkg.Type)
	w.Header().Set("X-Package-Name", pkg.Name)

	// magic
	w.Write([]byte("GMAD"))

	// gma version
	w.Write([]byte{3})

	// steamid (unused)
	author, err := strconv.Atoi(pkg.Author)
	if err != nil {
		utils.WriteError(w, r, fmt.Sprintf("failed to convert author steamid: %s", err))
		return
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

	// file list
	for i, content := range pkg.Content {
		// file number
		fileNum := make([]byte, 4)
		binary.LittleEndian.PutUint32(fileNum, uint32(i+1))
		w.Write(fileNum)

		// file name
		w.Write(append([]byte(content.Path), 0x00))

		// file size
		fileSize := make([]byte, 8)
		binary.LittleEndian.PutUint64(fileSize, uint64(content.Size))
		w.Write(fileSize)

		// file crc (skipped)
		w.Write(make([]byte, 4))
	}

	// end of file list marker
	w.Write(make([]byte, 4))

	// file content
	for _, content := range pkg.Content {
		data, err := utils.GetContentFile(content.ID, content.Revision)
		if err != nil {
			utils.WriteError(w, r, fmt.Sprintf("failed to get content file data: %s", err))
			return
		}

		w.Write(data)
	}

	// content crc (skipped)
	w.Write(make([]byte, 4))
}
