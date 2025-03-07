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

package utils

import (
	"archive/zip"
	"io/fs"
	"os"
	"path/filepath"
	"strconv"
)

// it's the caller's responsibility to close both the zip and content file
func GetContentFile(id int, rev int) (fs.File, *os.File, error) {
	f, err := os.Open(filepath.Join("data", "cdn", strconv.Itoa(id), strconv.Itoa(rev)))
	if err != nil {
		return nil, nil, err
	}

	stat, err := f.Stat()
	if err != nil {
		return nil, nil, err
	}

	zr, err := zip.NewReader(f, stat.Size())
	if err != nil {
		return nil, nil, err
	}

	zf, err := zr.Open("file")
	if err != nil {
		return nil, nil, err
	}

	return zf, f, nil
}
