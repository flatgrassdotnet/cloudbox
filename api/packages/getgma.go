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
	"archive/zip"
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
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

	var buf bytes.Buffer

	// magic
	buf.Write([]byte("GMAD"))

	// gma version
	buf.Write([]byte{3})

	// steamid (unused)
	author, err := strconv.Atoi(pkg.Author)
	if err != nil {
		utils.WriteError(w, r, fmt.Sprintf("failed to convert author steamid: %s", err))
		return
	}

	steamid := make([]byte, 8)
	binary.LittleEndian.PutUint64(steamid, uint64(author))
	buf.Write(steamid)

	// timestamp (unused)
	timestamp := make([]byte, 8)
	binary.LittleEndian.PutUint64(timestamp, uint64(pkg.Uploaded.Unix()))
	buf.Write(timestamp)

	// required content (unused)
	buf.Write([]byte{0x00})

	// addon name
	buf.Write(append([]byte(pkg.Name), 0))

	// addon description
	description, err := json.Marshal(GMADescription{
		Description: pkg.Description,
		Type:        pkg.Type,
		Tags:        []string{"fun"},
	})
	if err != nil {
		utils.WriteError(w, r, fmt.Sprintf("failed to marshal package description: %s", err))
		return
	}

	buf.Write(append([]byte(description), 0x00))

	// addon author
	buf.Write(append([]byte(pkg.AuthorName), 0x00))

	// addon version
	version := make([]byte, 4)
	binary.LittleEndian.PutUint32(version, uint32(pkg.Revision))

	buf.Write(version)

	// file list
	for i, content := range pkg.Content {
		// file number
		fileNum := make([]byte, 4)
		binary.LittleEndian.PutUint32(fileNum, uint32(i+1))
		buf.Write(fileNum)

		// file name
		buf.Write(append([]byte(content.Path), 0x00))

		// file size
		fileSize := make([]byte, 8)
		binary.LittleEndian.PutUint64(fileSize, uint64(content.Size))
		buf.Write(fileSize)

		// file crc (skipped)
		buf.Write(make([]byte, 4))
	}

	// end of file list marker
	buf.Write(make([]byte, 4))

	// file content
	for _, content := range pkg.Content {
		b, err := os.ReadFile(fmt.Sprintf("data/cdn/%d/%d", content.ID, content.Revision))
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

		buf.Write(fb)
	}

	// content crc (skipped)
	buf.Write(make([]byte, 4))

	w.Write(buf.Bytes())
}
