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

package db

import (
	"encoding/json"

	"github.com/flatgrassdotnet/cloudbox/common"
)

func InsertUpload(steamid int, upload common.Upload) (int, error) {
	includes, _ := json.Marshal(upload.Includes)

	r, err := handle.Exec("INSERT INTO uploads (steamid, type, meta, includes, data) VALUES (?, ?, ?, ?, ?)", steamid, upload.Type, upload.Metadata, includes, upload.Data)
	if err != nil {
		return 0, err
	}

	i, err := r.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(i), nil
}

func FetchUpload(id int) (common.Upload, error) {
	var upload common.Upload
	var includes string
	err := handle.QueryRow("SELECT type, meta, includes, data FROM uploads WHERE id = ?", id).Scan(&upload.Type, &upload.Metadata, &includes, &upload.Data)
	if err != nil {
		return upload, err
	}

	json.Unmarshal([]byte(includes), &upload.Includes)

	return upload, nil
}

func DeleteUpload(id int) error {
	_, err := handle.Exec("DELETE FROM uploads WHERE id = ?", id)
	if err != nil {
		return err
	}

	return nil
}
